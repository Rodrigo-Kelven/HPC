import socket
import os  # Para limpar a tela

def choose_node():
    """
    Conecta-se ao servidor central, exibe a lista de nós e permite que o usuário escolha
    com qual nó ele deseja se comunicar.
    """
    # Conectando ao servidor de descoberta
    discovery_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    discovery_socket.connect(('localhost', 9999))

    try:
        while True:
            # Recebe a lista de nós disponíveis
            node_list = discovery_socket.recv(1024).decode('utf-8')
            print(node_list)

            # Se não houver outros nós, informar o cliente e permitir a escolha de saída
            if "Não há outros nós para se comunicar" in node_list:
                print("Nenhum nó disponível para se comunicar.")

            # Se o cliente for o próprio nó, substitui o IP pelo 'localhost'
            node_list_with_localhost = node_list.replace(f"{discovery_socket.getsockname()[0]}:{discovery_socket.getsockname()[1]}", "localhost")

            print("\nLista atualizada de nós:")
            print(node_list_with_localhost)

            # Solicita uma escolha do usuário
            choice = input("\nEscolha o número do nó para se comunicar (ou 'X' para sair, 'u' para atualizar): ").strip()

            # Se o cliente digitar 'X', ele sai
            if choice.upper() == 'X':
                print("Saindo...")
                discovery_socket.close()
                return

            # Se o cliente digitar 'u', ele solicita uma atualização manual
            elif choice.upper() == 'U':
                # Limpa a tela antes de atualizar
                os.system('clear')
                print("Atualizando lista de nós...\n")
                discovery_socket.send(b"u")
                continue

            # Se o cliente digitar algo diferente de 'u' ou 'X', trata-se de um comando de atualização automática
            else:
                # Limpa a tela antes de atualizar
                os.system('clear')
                print("Atualizando lista de nós automaticamente...\n")
                discovery_socket.send(b"update")

            # Envia a escolha do nó para o servidor
            chosen_node = input("Escolha o número do nó para se comunicar (ou 'X' para sair): ")

            if chosen_node.upper() == 'X':
                print("Saindo...")
                discovery_socket.close()
                return

            # Envia o número do nó escolhido para o servidor
            discovery_socket.send(chosen_node.encode('utf-8'))

            # Recebe a resposta do servidor
            response = discovery_socket.recv(1024).decode('utf-8')
            print(response)

    except KeyboardInterrupt:
        print("\nProcesso interrompido pelo usuário. Saindo...")
        discovery_socket.close()

if __name__ == "__main__":
    choose_node()
