package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"strconv"
)

var (
	// Lista para armazenar os nós conectados
	connectedPeers []net.Conn
	mu             sync.Mutex
)

// Função para atualizar a lista de nós e enviar para todos os clientes
func updateNodeList() {
	// Cria uma lista de nós conectados, mas sempre inclui o nó local (localhost) como primeiro
	var nodeList []string
	mu.Lock()
	for i, peer := range connectedPeers {
		// Se o cliente for o próprio nó, substitui o IP por "localhost"
		if peer.LocalAddr().String() == peer.RemoteAddr().String() {
			nodeList = append([]string{fmt.Sprintf("1. localhost:%s", peer.LocalAddr().String())}, nodeList...)
		} else {
			nodeList = append(nodeList, fmt.Sprintf("%d. %s", i+1, peer.RemoteAddr().String()))
		}
	}
	mu.Unlock()

	// Envia a lista de nós para todos os clientes
	if len(nodeList) > 0 {
		message := "Lista de nós conectados:\n" + strings.Join(nodeList, "\n") + "\nEscolha o número do nó para se comunicar (ou 'X' para sair): "
		mu.Lock()
		for _, client := range connectedPeers {
			client.Write([]byte(message))
		}
		mu.Unlock()
	}
}

// Função para gerenciar a comunicação com cada cliente
func handleClient(client net.Conn) {
	defer client.Close()

	// Adiciona o nó à lista de nós conectados
	mu.Lock()
	connectedPeers = append(connectedPeers, client)
	mu.Unlock()
	fmt.Printf("Novo nó conectado: %s. Total de nós: %d\n", client.RemoteAddr().String(), len(connectedPeers))

	// Envia a lista de nós para o cliente
	updateNodeList()

	buf := make([]byte, 1024)
	for {
		// Lê a mensagem do cliente
		n, err := client.Read(buf)
		if err != nil {
			fmt.Println("Erro ao ler mensagem:", err)
			break
		}
		clientMessage := string(buf[:n])

		// Se o cliente escolheu sair
		if strings.ToUpper(clientMessage) == "X" {
			fmt.Printf("Nó %s escolheu sair.\n", client.RemoteAddr().String())
			break
		}

		// Se o cliente solicitou atualização da lista de nós
		if clientMessage == "u" || clientMessage == "update" {
			fmt.Printf("Nó %s solicitou atualização da lista de nós.\n", client.RemoteAddr().String())
			updateNodeList()
		} else {
			// Processa a escolha de nó
			if chosenNodeNumber, err := strconv.Atoi(clientMessage); err == nil && chosenNodeNumber >= 1 && chosenNodeNumber <= len(connectedPeers) {
				selectedNode := connectedPeers[chosenNodeNumber-1]
				client.Write([]byte(fmt.Sprintf("Você escolheu o nó %s. Conectando...\n", selectedNode.RemoteAddr().String())))
				fmt.Printf("Nó %s escolheu o nó %s para se comunicar.\n", client.RemoteAddr().String(), selectedNode.RemoteAddr().String())
			} else {
				client.Write([]byte("Escolha inválida. Conexão encerrada.\n"))
			}
		}
	}

	// Remove o nó da lista e atualiza a lista
	mu.Lock()
	for i, peer := range connectedPeers {
		if peer == client {
			connectedPeers = append(connectedPeers[:i], connectedPeers[i+1:]...)
			break
		}
	}
	mu.Unlock()

	updateNodeList() // Atualiza a lista para todos os clientes
}

// Função para iniciar o servidor de descoberta
func startDiscoveryServer() {
	// Inicia o servidor de descoberta na porta 9999
	ln, err := net.Listen("tcp", "localhost:9999")
	if err != nil {
		fmt.Println("Erro ao iniciar servidor:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Servidor de descoberta iniciado...")

	for {
		// Aceita novas conexões
		client, err := ln.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}

		// Cria uma goroutine para gerenciar cada cliente
		go handleClient(client)
	}
}

func main() {
	startDiscoveryServer()
}
