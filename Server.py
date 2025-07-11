import socket
import threading

# Lista para armazenar os nós conectados
connected_peers = []

def update_node_list():
    """
    Envia a lista de nós para todos os clientes conectados,
    atualizando-os sempre que há uma mudança.
    """
    node_list = "\n".join([f"{i + 1}. {peer.getpeername()[0]}:{peer.getpeername()[1]}" for i, peer in enumerate(connected_peers)])
    
    # Atualiza a lista para todos os clientes
    if node_list:
        for client_socket in connected_peers:
            try:
                client_socket.send(f"Lista de nós conectados:\n{node_list}\nDigite o número do nó para se comunicar (ou 'X' para sair): ".encode('utf-8'))
            except:
                pass  # Se o cliente desconectar, ignoramos o erro

def handle_client(client_socket, client_address):
    """
    Gerencia a comunicação com um cliente.
    """
    # Adiciona o nó à lista de nós conectados
    connected_peers.append(client_socket)
    print(f"Novo nó conectado: {client_address}. Total de nós: {len(connected_peers)}")

    # Envia a lista de nós para o cliente
    update_node_list()

    try:
        while True:
            # Recebe a mensagem do cliente
            client_message = client_socket.recv(1024).decode('utf-8').strip()

            if client_message.upper() == 'X':
                print(f"Nó {client_address} escolheu sair.")
                break

            if client_message == "u" or client_message == "update":
                print(f"Nó {client_address} solicitou atualização da lista de nós.")
                update_node_list()
            else:
                # Caso o nó escolha algum número, processa a escolha
                try:
                    chosen_node_number = int(client_message)
                    if 1 <= chosen_node_number <= len(connected_peers):
                        selected_node = connected_peers[chosen_node_number - 1]
                        client_socket.send(f"Você escolheu o nó {selected_node.getpeername()}. Conectando...".encode('utf-8'))
                        print(f"Nó {client_address} escolheu o nó {selected_node.getpeername()} para se comunicar.")
                    else:
                        client_socket.send("Escolha inválida. Conexão encerrada.".encode('utf-8'))
                except ValueError:
                    client_socket.send("Entrada inválida. Conexão encerrada.".encode('utf-8'))

    finally:
        # Remove o nó da lista e fecha a conexão
        connected_peers.remove(client_socket)
        client_socket.close()
        update_node_list()  # Atualiza a lista para todos os clientes

def start_discovery_server():
    """
    Inicia o servidor de descoberta.
    """
    discovery_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    discovery_socket.bind(('localhost', 9999))  # Configura o servidor na porta 9999
    discovery_socket.listen(5)

    print("Servidor de descoberta iniciado...")

    while True:
        # Aceita novas conexões de clientes
        client_socket, client_address = discovery_socket.accept()

        # Cria uma thread para gerenciar cada novo cliente
        client_thread = threading.Thread(target=handle_client, args=(client_socket, client_address))
        client_thread.start()

if __name__ == "__main__":
    start_discovery_server()
