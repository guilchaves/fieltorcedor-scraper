# FielTorcedorBot

Um bot em Go para monitorar o site Fiel Torcedor e enviar notificações automáticas sobre novos jogos do Corinthians, incluindo lembretes para abertura de categorias de ingressos, via Telegram.

## Funcionalidades

- Busca automática de novos jogos do Corinthians no site Fiel Torcedor.
- Envio de notificações detalhadas para o Telegram, com emojis e formatação.
- Lembretes automáticos para abertura de categorias de ingressos.
- Comandos interativos no Telegram:
  - `/start` — Mensagem de boas-vindas.
  - `/jogos` — Lista todos os próximos jogos encontrados.
  - `/help` — Mostra todos os comandos disponíveis.
- Evita envio de mensagens duplicadas, mesmo após reiniciar o bot.
- Arquitetura limpa (hexagonal), facilitando manutenção e extensibilidade.

## Como usar

### 1. Pré-requisitos

- Go 1.24 ou superior
- Uma conta de bot no Telegram ([crie com o BotFather](https://t.me/botfather))
- Um chat ID do Telegram (pode ser seu próprio usuário ou grupo)
- (Opcional) Docker, se quiser rodar em container

### 2. Configuração

Crie um arquivo `.env` na raiz do projeto com o seguinte conteúdo:

```env
TELEGRAM_BOT_TOKEN=seu_token_aqui
TELEGRAM_CHAT_ID=seu_chat_id_aqui
DATABASE_URL=seu_host_postgresql_aqui
```

### 3. Instale as dependências

```bash
go mod tidy
```
### 4. Execute o bot

```bash
go run ./cmd/bot/main.go
```

O bot irá:</br>
Fazer uma busca inicial por jogos.</br>
Agendar buscas automáticas diárias às 11h, 15h, 17h e 18h.</br>
Responder aos comandos do Telegram.</br>

### 5. Comandos do Telegram

    `/start` — Mensagem de boas-vindas.</br>
    `/jogos` — Lista todos os próximos jogos do Corinthians, com detalhes.</br>
    `/help` — Mostra todos os comandos disponíveis.</br>

### 6. Estrutura do Projeto
```
internal/
  adapters/
    in/                # Schedulers e entrada de comandos
    out/
      fieltorcedor/    # Scraper do site Fiel Torcedor
      notification/    # Envio de mensagens para o Telegram
      notifiedgames/   # Persistência dos jogos notificados
  core/
    domain/            # Modelos de domínio (Game, Category, etc)
    ports/             # Interfaces (ports) da arquitetura hexagonal
    service/           # Lógica de negócio (casos de uso)
  handlers/            # Handlers dos comandos do Telegram
cmd/
  bot/                 # main.go (ponto de entrada)
.env                   # Variáveis de ambiente
notified_games.txt     # Persistência dos jogos já notificados
```
### 7. Personalização

    Para adicionar novos comandos, crie um novo handler em internal/handlers/ e registre no main.go.
    Para alterar os horários das buscas automáticas, edite os crons no main.go.

### 8. Licença

[MIT](https://github.com/guilchaves/fieltorcedorbot/blob/main/README.md)
