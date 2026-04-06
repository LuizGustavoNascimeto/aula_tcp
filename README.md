# Atividade 1 - Programação com Sockets TCP

Projeto da disciplina **Sistemas Distribuídos** Prof. Rodrigo Campiolo.

Este repositório contém uma implementação simples de comunicação **cliente-servidor usando TCP** em Go.
O foco da atividade é praticar sockets TCP, troca de mensagens em formato texto (String UTF) e suporte a múltiplos clientes.

## Estrutura do projeto

- `server/server.go`: inicia o servidor TCP na porta `8080`, aceita múltiplas conexões e responde cada mensagem com `ACK: <mensagem>`.
- `client/client.go`: conecta ao servidor TCP, lê texto do terminal, envia ao servidor e valida a resposta.

## Como funciona

1. O servidor abre a porta `:8080`.
2. Cada cliente que conecta é tratado em uma goroutine separada (concorrência).
3. O cliente envia uma mensagem (texto digitado no terminal).
4. O servidor recebe e responde com um ACK.

## Pré-requisitos

- Go instalado (versão 1.20+ recomendada).
- Terminal (PowerShell, CMD, Git Bash ou terminal integrado do VS Code).

## Como rodar o projeto

No diretório raiz do projeto (`aula_tcp`), siga os passos:

1. Inicie o servidor:

```bash
go run server/server.go
```

2. Em outro terminal, inicie um cliente:

```bash
go run client/client.go
```

3. Digite uma mensagem no terminal do cliente e pressione Enter.

4. Observe:

- No servidor: mensagem recebida do cliente.
- No cliente: confirmação da resposta `ACK`.

## Testar múltiplos clientes

Para validar que o servidor atende múltiplos clientes, abra mais terminais e execute novamente:

```bash
go run client/client.go
```

Cada cliente pode enviar mensagens de forma independente.

## Observações

- A leitura de mensagens está baseada em `\n` (Enter no terminal), ou seja, cada mensagem deve terminar com quebra de linha.
- Se o servidor não estiver ativo, o cliente não consegue conectar.
- A porta padrão do projeto é `8080`.

## Contexto acadêmico

Este projeto foi feito como atividade prática de faculdade para reforçar conceitos de:

- sockets TCP;
- concorrência com goroutines;
- comunicação cliente-servidor;
- noções iniciais de sistemas distribuídos.
