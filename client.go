package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"
)

//func main() {
//	conn, _ := net.Dial("tcp", ":8081")
//	for {
//		reader := bufio.NewReader(os.Stdin)
//		fmt.Print("Text to send: ")
//		text, _ := reader.ReadString('\n')
//		conn.Write([]byte(text))
//		msg,_ := bufio.NewReader(conn).ReadString('\n')
//		fmt.Print("Message from server: ", msg)
//	}
//}

func main() {
	wg := sync.WaitGroup{}
	numberOfUsers := 5
	for i := 0; i < numberOfUsers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			open()
		}()
	}
	wg.Wait()
}

func open() {
	conn, _ := net.Dial("tcp", ":8081")
	for i := 0; i < 10; i++ {
		text := strconv.Itoa(rand.Intn(1000))
		fmt.Fprintln(conn, text)
		msg, _ := bufio.NewReader(conn).ReadString('\n')
		log.Println(msg)
		time.Sleep(time.Millisecond*500)
	}
	conn.Close()
}