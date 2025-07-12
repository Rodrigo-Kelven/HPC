package main

import (
	"fmt"
	"net"
)

func chooseNode() {
	// Conectando ao servidor central
	discoverySocket, err := net.Dial("tcp", "localhost:9999")
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		return
	}
	defer discoverySocket.Close()

	for {
		// Recebendo a lista de nós disponíveis
		nodeList := make([]byte, 1024)
		n, err := discoverySocket.Read(nodeList)
		if err != nil {
			fmt.Println("Erro ao receber dados:", err)
			return
		}

		// Limpa a tela (comando para Linux/Unix)
		fmt.Print("\033[H\033[2J") // ANSI escape code para limpar a tela
		// Exibe a lista dos outros nós conectados
		fmt.Println(string(nodeList[:n]))

		// Solicita uma escolha do usuário
		var choice string
		fmt.Print("\nEscolha o número do nó para se comunicar (ou 'X' para sair, 'u' para atualizar): ")
		fmt.Scanln(&choice)

		if choice == "X" || choice == "x" {
			fmt.Println("Saindo...")
			break
		} else if choice == "u" || choice == "U" {
			// Atualiza a lista de nós
			discoverySocket.Write([]byte("u"))
			continue
		} else {
			// Envia a escolha do nó para o servidor
			discoverySocket.Write([]byte(choice))
		}
	}
}

func main() {
	chooseNode()
}
