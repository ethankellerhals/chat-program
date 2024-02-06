package main

import (
	"bufio"
	"log"
	"net"
	"strings"
	"sync"
)

type Client struct {
    conn net.Conn
    name string
}

var (
    clients   []*Client
    clientsMu sync.Mutex
)

func main() {
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        log.Fatal("Error starting server:", err)
    }
    defer listener.Close()

    log.Println("Server started, listening on port 8080")

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Println("Error accepting connection:", err)
            continue
        }

        go handleClient(conn)
    }
}
func handleClient(conn net.Conn) {
    defer conn.Close()

    log.Println("New client connected:", conn.RemoteAddr())

    client := &Client{conn: conn}

    // Ask the client for their name
    client.write("Enter your name: ")
    reader := bufio.NewReader(conn)
    name, err := reader.ReadString('\n')
    if err != nil {
        log.Println("Error reading client name:", err)
        return
    }
    client.name = strings.TrimSpace(name)

    clientsMu.Lock()
    clients = append(clients, client)
    clientsMu.Unlock()

    client.write("Welcome, " + client.name + "!\n")

    // Read messages from the client
    for {
        message, err := reader.ReadString('\n')
        if err != nil {
            break
        }
        broadcast(client.name+": "+strings.TrimSpace(message), client)
    }

    // Remove the client from the list
    clientsMu.Lock()
    defer clientsMu.Unlock()
    for i, c := range clients {
        if c == client {
            clients = append(clients[:i], clients[i+1:]...)
            break
        }
    }
    log.Println("Client disconnected:", client.name)
}

func (c *Client) write(message string) {
    _, err := c.conn.Write([]byte(message))
    if err != nil {
        log.Println("Error writing to client:", err)
    }
}

func broadcast(message string, sender *Client) {
	message = "You: " + message + "\n"
    clientsMu.Lock()
    defer clientsMu.Unlock()
    for _, client := range clients {
        if client != sender {
            client.write(message)
			log.Printf("Broadcasting message to %s: %s", client.name, message)
        }
    }
}
