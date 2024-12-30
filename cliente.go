package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Conecta ao servidor na porta 8080
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Erro ao conectar no servidor:", err)
		os.Exit(1)
	}
	defer conn.Close()

	go func() {
		// Lê mensagens do servidor e exibe no terminal
		for {
			message, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Println("Erro de leitura do servidor:", err)
				break
			}
			fmt.Print(message)
		}
	}()

	// Envia mensagens para o servidor
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Você: ")
		scanner.Scan()
		message := scanner.Text()

		if message == "sair" {
			conn.Write([]byte("sair\n"))
			break
		}

		conn.Write([]byte(message + "\n"))
	}
}
