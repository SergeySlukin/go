package main

import (
	"net"
	"fmt"
	"bufio"
)

func main()  {
	listener, _ := net.Listen("tcp", ":5000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("can not connect!")
			conn.Close()
			continue
		}

		fmt.Println("Connected")

		//создаем реадер для чтения информации из сокета
		bufReader := bufio.NewReader(conn)
		fmt.Println("Start reading")
		go func(conn net.Conn) {
			defer conn.Close()
			for {
				rbyte, err :=  bufReader.ReadByte()
				if err != nil {
					fmt.Println("Can not read!", err)
					break
				}

				fmt.Print(string(rbyte))
			}
		}(conn)

	}
}
