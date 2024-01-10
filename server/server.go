package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Data struct {
	Username string
	Message  string
}

const (
	connHost = "localhost"
	connPort = "8080"
	connType = "tcp"
)

// ANSI escape codes for text colors
var colors = []string{
	"\033[31m", // Red
	"\033[32m", // Green
	"\033[33m", // Yellow
	"\033[34m", // Blue
	"\033[35m", // Magenta
	"\033[36m", // Cyan
}

var clients = make(map[string]net.Conn) // Map to store connections using UUID
var sub = make(map[string](map[string]bool))
var mutex = &sync.Mutex{}

func handleConnection(conn net.Conn) {
	clientUUID := uuid.New().String()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	chosenColor := colors[r.Intn(len(colors))]
	group, _ := bufio.NewReader(conn).ReadString('\n')
	group = group[:len(group)-1]
	conn.Write([]byte(chosenColor + "\n"))
	fmt.Println("New client connected:", clientUUID)

	mutex.Lock()
	clients[clientUUID] = conn

	mutex.Unlock()

	defer func() {
		mutex.Lock()
		delete(clients, clientUUID)
		mutex.Unlock()
		conn.Close()
	}()
	if sub[group] == nil {
		sub[group] = make(map[string]bool)
	}
	sub[group][clientUUID] = true
	var data Data
	for {
		decoder := gob.NewDecoder(conn)
		err := decoder.Decode(&data)
		if err != nil {
			fmt.Println("Client", clientUUID, "left.")
			return
		}

		for g, v := range sub {
			if g == group {
				for i, f := range v {
					if i != clientUUID && f {
						clients[i].Write([]byte("[" + data.Username + "]: " + data.Message))
					}
				}
			}
		}

	}
}

func main() {
	fmt.Println("Starting", connType, "server on", connHost+":"+connPort)
	l, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		log.Fatalf("Error listening: %v", err)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatalf("Error connecting: %v", err)
		}
		go handleConnection(c)
	}
}
