# Sistema Cliente-Servidor TCP em Go

Este projeto foi desenvolvido como atividade avaliativa da disciplina Sistemas Distribuídos, ministrada pelo professor doutor Rodrigo Cambiolo.

O objetivo é implementar uma aplicação cliente-servidor com sockets TCP, autenticação simples e comandos de navegação em diretórios por sessão de usuário.

## Visão geral

O sistema é composto por:

- um servidor TCP que escuta na porta 8080;
- um cliente de terminal que envia comandos em texto;
- gerenciamento de sessão por conexão;
- acesso restrito a diretórios de cada usuário dentro da pasta users_files.

Cada conexão de cliente é tratada em uma goroutine independente no servidor.

## Estrutura do projeto

- client/client.go: cliente TCP interativo via terminal.
- tcp_server/main.go: ponto de entrada do servidor.
- tcp_server/src/server: lógica de rede, leitura e resposta de comandos.
- tcp_server/src/commands: implementação dos comandos suportados.
- tcp_server/src/session: controle de sessões por cliente conectado.
- tcp_server/src/user: controle de usuários e autenticação.
- tcp_server/src/utils: validações e listagem de diretórios/arquivos.
- users_files: área de arquivos acessível pelo servidor.

## Pré-requisitos

- Go instalado (versão 1.25.1 ou superior recomendada).
- Terminal (Git Bash, PowerShell, CMD ou terminal integrado do VS Code).

## Como executar

### 1. Iniciar o servidor

No terminal 1:

```bash
cd tcp_server
go run ./
```

O servidor ficará ouvindo em :8080.

### 2. Iniciar o cliente

No terminal 2, a partir da raiz do projeto:

```bash
go run client/client.go
```

Você verá o prompt:

```text
>>>
```

## Comandos disponíveis no cliente

### CONNECT

Formato esperado:

```text
CONNECT usuario, senha
```

Exemplo:

```text
CONNECT admin, password
```

Observação: o cliente transforma automaticamente a senha em hash antes de enviar para o servidor.

### PWD

Retorna o diretório atual da sessão.

```text
PWD
```

### CHDIR

Muda o diretório atual da sessão (somente dentro da raiz autorizada do usuário).

```text
CHDIR /admin
CHDIR /admin/fotos
```

### GETDIRS

Lista subdiretórios do diretório atual.

```text
GETDIRS
```

### GETFILES

Lista arquivos do diretório atual.

```text
GETFILES
```

### EXIT

Encerra a sessão no servidor e fecha a conexão do cliente.

```text
EXIT
```

## Usuário padrão para testes

O servidor cria automaticamente o usuário abaixo na inicialização:

- usuário: admin
- senha: password

## Fluxo sugerido de teste

Após iniciar servidor e cliente, execute:

```text
CONNECT admin, password
PWD
GETDIRS
CHDIR /admin/fotos
PWD
GETFILES
EXIT
```

## Possíveis erros comuns

- ERROR: NOT_AUTHENTICATED: enviado comando sem autenticação prévia.
- ERROR: INVALID_COMMAND: sintaxe inválida do comando.
- ERROR: INVALID_DIRECTORY: tentativa de acessar diretório fora da raiz permitida.
- erro ao conectar: servidor não está em execução na porta 8080.

## Objetivo pedagógico

Este projeto reforça os seguintes conceitos:

- comunicação cliente-servidor com TCP;
- concorrência com goroutines;
- controle de sessão por conexão;
- validação de acesso a diretórios;
- organização de código em camadas (comandos, sessão, servidor e utilitários).
