package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

func main() {
	listen, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln(err)
	}
	defer listen.Close()
	log.Println("Server is listening...")

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println(err)
			conn.Close()
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	input := make([]byte, 1024*4)

	n, err := conn.Read(input)
	if err != nil || n == 0 {
		log.Println("Read error:", err)
		return
	}
	source := string(input[0:n])
	res, err := multiplication(source)
	if err != nil {
		log.Println("Multiplication error:", err)
		return
	}

	n, err = conn.Write([]byte(res))
	if err != nil || n == 0 {
		log.Println("Write error:", err)
		return
	}
}

func multiplication(line string) (string, error) {
	fields := strings.Fields(line)
	if len(fields) < 1 {
		return "", fmt.Errorf("empty line")
	}

	var res string
	for _, f := range fields {
		split := strings.Split(f, ",")
		num1, err := strconv.Atoi(split[0])
		if err != nil {
			return "", err
		}
		num2, err := strconv.Atoi(split[1])
		if err != nil {
			return "", err
		}
		res += fmt.Sprintf("%d\r\n", num1*num2)
	}
	res += fmt.Sprintf("\r\n")
	return res, nil
}
