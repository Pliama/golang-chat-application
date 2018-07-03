package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Println(err.Error())
	}
	users := make(map[int]string)
	//TODO: add check for existing usernames
	aconns := make(map[net.Conn]int)
	conns := make(chan net.Conn)
	dconns := make(chan net.Conn)
	msgs := make(chan string)
	i := 0

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err.Error())
			}
			fmt.Println("Client connected ", conn.LocalAddr())
			conns <- conn
			rd := bufio.NewReader(conn)
			nick, err := rd.ReadString('\n')
			nick = nick[0 : len(nick)-2]
			if err != nil {
				log.Println(err.Error())
			}
			users[i] = nick
		}
	}()

	for {
		select {
		//read incomming connections
		case conn := <-conns:
			aconns[conn] = i
			i++
			go func(conn net.Conn, i int) {
				rd := bufio.NewReader(conn)
				for {
					m, err := rd.ReadString('\n')
					if err != nil {
						break
					}
					name := users[i]
					msgs <- fmt.Sprintf("%v:%v", name, m)
				}
				//Done reading from it
				dconns <- conn
			}(conn, i)
		case msg := <-msgs:
			//broadcast
			for conn := range aconns {
				conn.Write([]byte(msg))
			}
		case dconn := <-dconns:
			log.Printf("Client %v disconnected\n", aconns[dconn])
			delete(users, i)
			delete(aconns, dconn)
		}
	}
}
