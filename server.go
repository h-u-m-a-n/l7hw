package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"
)

type Semaphore struct {
	ch chan struct{}
}

func NewSemaphore(size int) *Semaphore {
	return &Semaphore{ch: make(chan struct{}, size)}
}

func (s *Semaphore) Acquire(n int) {
	for i := 0; i < n; i++ {
		s.ch <- struct{}{}
	}
}

func (s *Semaphore) Release(n int) {
	for i := 0; i < n; i++ {
		<- s.ch
	}
}

func square(n int, sec int) int {
	time.Sleep(time.Second*time.Duration(sec))
	return n*n
}

func handleConnection(conn net.Conn, ctx context.Context) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {

		text := scanner.Text()
		n, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			log.Fatal("error while parsing: ", err)
		}
		result := square(int(n), rand.Intn(3))
		_, err = fmt.Fprintf(conn, "%v: client - %v squared is %v\n", conn.RemoteAddr() ,n, result)
		if err != nil {
			log.Fatal("error while responding to client")
		}

		select {
		case <-ctx.Done():
			fmt.Fprint(conn, conn.RemoteAddr(), " - Server is down\n")
			return
		default:

		}
	}
}

func main() {
	args := os.Args

	if len(args) == 1 {
		log.Println("Enter maximum number of connections.")
		return
	}
	n, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatal("error occurred: ", err)
	}

	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Accept connections on port")

	var (
		wg sync.WaitGroup
		ctx, _ = signal.NotifyContext(context.Background(), os.Interrupt)
		semaphore = NewSemaphore(n)
	)
	//defer stopFunc()

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Server is closing")
				return
			default:
				conn, err := ln.Accept()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("Calling handleConnection")
				wg.Add(1)
				semaphore.Acquire(1)
				go func() {
					defer wg.Done()
					defer semaphore.Release(1)
					handleConnection(conn, ctx)
				}()
			}
		}
	}()
	<-ctx.Done()
	wg.Wait()
	log.Println("Server closed")
}
