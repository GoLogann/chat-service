# Chat Service

Este é um serviço de chat que utiliza WebSocket, AWS SQS, Redis e Go para gerenciar sessões de usuários e mensagens em tempo real.

## Funcionalidades

- **WebSocket**: Conexões em tempo real para comunicação bidirecional com os usuários.
- **AWS SQS**: Envia e recebe mensagens para processamento assíncrono.
- **Redis**: Cache de mensagens para otimizar o tempo de resposta.
- **Worker Pool**: Gerenciamento de tarefas assíncronas para processamento de mensagens.

## Arquitetura

A arquitetura do sistema é composta por:

- **WebSocket Manager**: Gerencia as conexões WebSocket ativas e trata o envio de mensagens.
- **SQS Producer**: Envia mensagens para uma fila SQS para processamento assíncrono.
- **SQS Consumer**: Consome mensagens da fila SQS, processa-as e atualiza o cache Redis.
- **Redis Cache**: Armazena respostas de mensagens para garantir que mensagens já enviadas não sejam repetidas.

## Configuração

### Dependências

- Go 1.23.5
- Redis
- AWS SDK (para SQS)
- Gorilla WebSocket
