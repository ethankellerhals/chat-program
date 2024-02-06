package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}
	defer conn.Close()

	go receiveMessages(conn) // Add goroutine to handle receiving messages

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter your name: ")
	scanner.Scan()
	name := scanner.Text()

	_, err = fmt.Fprintf(conn, "%s\n", name)
	if err != nil {
		log.Println("Error sending name:", err)
		return
	}

	for scanner.Scan() {
		message := scanner.Text()
		_, err := conn.Write([]byte(strings.TrimSpace(message) + "\n"))
		if err != nil {
			log.Println("Error sending message:", err)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error reading standard input:", err)
	}
}

func receiveMessages(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		if strings.HasPrefix(message, "You:") {
			fmt.Println("\033[1m" + message + "\033[0m")
		} else {
			fmt.Println(scanner.Text())
		}
	}
	if err := scanner.Err(); err != nil {
		log.Println("Error reading from server:", err)
	}
}
