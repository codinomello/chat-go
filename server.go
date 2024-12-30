package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

var clients []net.Conn
var mu sync.Mutex

// Função para broadcast de mensagens para todos os clientes
func broadcast(message string, sender net.Conn) {
	mu.Lock()
	defer mu.Unlock()

	for _, client := range clients {
		if client != sender {
			_, err := client.Write([]byte(message))
			if err != nil {
				fmt.Println("Erro ao enviar mensagem para o cliente:", err)
			}
		}
	}
}

// Função para gerenciar a comunicação de um cliente
func handleClient(conn net.Conn) {
	defer conn.Close()
	clients = append(clients, conn)
	defer func() {
		mu.Lock()
		// Remove o cliente desconectado
		for i, client := range clients {
			if client == conn {
				clients = append(clients[:i], clients[i+1:]...)
				break
			}
		}
		mu.Unlock()
	}()

	conn.Write([]byte("Bem-vindo ao chat! Digite suas mensagens.\n"))

	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Erro de leitura do cliente:", err)
			break
		}
		message := string(buffer[:n])
		message = strings.TrimSpace(message)
		if message == "sair" {
			conn.Write([]byte("Você saiu do chat.\n"))
			break
		}
		// Broadcast da mensagem para todos os clientes conectados
		broadcast(message+"\n", conn)
	}
}

func main() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
		os.Exit(1)
	}
	defer listen.Close()

	fmt.Println("Servidor de chat iniciado na porta 8080")

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}
		go handleClient(conn)
	}
}
