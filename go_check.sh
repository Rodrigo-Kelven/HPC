#!/bin/bash

# Função para exibir mensagens
function info() {
    echo -e "\033[1;32m$1\033[0m"
}

# Verificar se as ferramentas estão instaladas
function check_tools() {
    for tool in gofmt golint goimports staticcheck; do
        if ! command -v $tool &> /dev/null; then
            echo "Erro: A ferramenta '$tool' não está instalada."
            exit 1
        fi
    done
}

# Formatar o código com gofmt
function run_gofmt() {
    info "Rodando gofmt para formatação de código..."
    gofmt -w .
    if [ $? -eq 0 ]; then
        echo "gofmt: Código formatado com sucesso."
    else
        echo "gofmt: Erro ao formatar código."
    fi
}

# Rodar o golint para análise de boas práticas
function run_golint() {
    info "Rodando golint para análise de boas práticas..."
    golint ./...
    if [ $? -eq 0 ]; then
        echo "golint: Código aprovado."
    else
        echo "golint: Encontrados problemas de estilo."
    fi
}

# Rodar o goimports para formatar e organizar as importações
function run_goimports() {
    info "Rodando goimports para organizar importações..."
    goimports -w .
    if [ $? -eq 0 ]; then
        echo "goimports: Importações organizadas com sucesso."
    else
        echo "goimports: Erro ao organizar importações."
    fi
}

# Rodar o staticcheck para análise estática avançada
function run_staticcheck() {
    info "Rodando staticcheck para análise estática avançada..."
    staticcheck ./...
    if [ $? -eq 0 ]; then
        echo "staticcheck: Código aprovado."
    else
        echo "staticcheck: Encontrados problemas no código."
    fi
}

# Rodar todas as verificações
function run_all_checks() {
    check_tools
    run_gofmt
    run_golint
    run_goimports
    #run_staticcheck
}

# Executar as verificações
run_all_checks
