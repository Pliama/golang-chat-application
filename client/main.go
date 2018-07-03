package main

import (
	"net"
	"bufio"
	"fmt"
	"os"
	"time"
)

func main(){
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		if _, t := err.(*net.OpError); t {
			fmt.Println("Some problem connecting.")
		} else {
			fmt.Println("Unknown error: " + err.Error())
		}
		os.Exit(1)
	}
	go readConnection(conn)
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("text:>")
		text, _ := reader.ReadString('\n')
		conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
		_, err := conn.Write([]byte(text))
		if err != nil {
			fmt.Println("Error writing to stream.")
			break
		}
	}
}
func readConnection(conn net.Conn) {
	for {
		scanner := bufio.NewScanner(conn)

		for {
			ok := scanner.Scan()
			text := scanner.Text()
				fmt.Println(text)
			if !ok {
				fmt.Println("Reached EOF on server connection.")
				break
			}
		}
	}
}
