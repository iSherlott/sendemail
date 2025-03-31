# 📧 Serviço de Envio de E-mails

Este serviço permite o envio de e-mails utilizando **RabbitMQ**, **Azure Blob Storage** para templates e suporte a **anexos** via **Base64** ou **URL**.

## 📌 Funcionalidades
- 📬 **Envio de e-mails assíncrono** utilizando RabbitMQ
- 🖼 **Templates dinâmicos** carregados do Azure Blob Storage
- 📎 **Suporte a anexos** em Base64 ou via URL
- 📜 **Logs controlados por ambiente** (Development vs. Produção)
- 🔧 **Configuração via `.env`**

---

## 🏗 Estrutura do Projeto
```
📂 sendemail/
├── 📂 internal/
│   ├── 📂 models/            # Modelos de dados (EmailPayload, Attachment)
│   ├── 📂 consumer/          # Lógica para consumir mensagens do RabbitMQ
├── 📂 pkg/
│   ├── 📂 azure/             # Recuperação de templates do Azure Blob Storage
│   ├── 📂 email/             # Configuração e envio de e-mails
│   ├── 📂 rabbitmq/          # Conexão com o RabbitMQ
│   ├── 📂 utils/             # Funções auxiliares (logs condicionais, etc.)
├── 📂 service/               # Lógica principal do envio de e-mails
├── .env                      # Configurações do ambiente
├── go.mod                     # Dependências do projeto
└── README.md                  # Documentação do projeto
```

---

## 🚀 Como Rodar o Serviço

### 📌 1. Configurar Variáveis de Ambiente
Crie um arquivo `.env` na raiz do projeto e defina as seguintes variáveis:

```ini
# Configuração do RabbitMQ
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
RABBITMQ_QUEUE=email_queue

# Configuração do SMTP
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_EMAIL=seuemail@gmail.com
SMTP_PASSWORD=suasenha

# Configuração do Azure Blob Storage
AZURE_STORAGE_ACCOUNT=seuarmazem
AZURE_STORAGE_KEY=sua_chave_base64
AZURE_CONTAINER_NAME=templates

# Ambiente de execução
ENV=Development  # Ou "Production"
```

### 📌 2. Instalar Dependências
```sh
go mod tidy
```

### 📌 3. Rodar o Serviço
```sh
go run main.go
```

O serviço ficará escutando mensagens do RabbitMQ e processará o envio de e-mails automaticamente.

---

## 📩 Payload de Envio de E-mail
O serviço recebe mensagens no RabbitMQ com o seguinte formato:

```json
{
  "to": ["destinatario@email.com"],
  "subject": "Bem-vindo ao nosso serviço!",
  "body": "Olá, seja bem-vindo!",
  "unsubscribeLink": "https://meusite.com/unsubscribe",
  "template": "welcome",
  "attachments": [
    {
      "filename": "documento.pdf",
      "content": "Base64EncodedStringAqui"
    },
    {
      "filename": "imagem.png",
      "url": "https://meusite.com/imagem.png"
    }
  ]
}
```

### 📌 Explicação dos Campos
- **to**: Lista de destinatários do e-mail.
- **subject**: Assunto do e-mail.
- **body**: Conteúdo HTML do e-mail (substituído no template).
- **unsubscribeLink**: Link para descadastramento.
- **template**: Nome do template armazenado no Azure Blob Storage (exemplo: `welcome.html`).
- **attachments**: Lista de anexos. Cada anexo pode ter:
  - `filename`: Nome do arquivo.
  - `content`: String em **Base64** do arquivo (opcional).
  - `url`: Link do arquivo (opcional).

---

## 🛠 Como Funciona o Serviço

1. O serviço **escuta mensagens** no RabbitMQ.
2. Ao receber uma mensagem, ele **busca o template no Azure Blob Storage**.
3. O template é **renderizado** substituindo os placeholders com os dados do payload.
4. O e-mail é **enviado via SMTP**, anexando arquivos conforme necessário.

---

## 📜 Logs Controlados por Ambiente
A função `utils.LogIfDevelopment` permite exibir logs apenas em ambiente de desenvolvimento.

```go
utils.LogIfDevelopment("📩 Mensagem recebida: %s", msg.Body)
```

Em `Production`, os logs serão omitidos automaticamente.

---

## 🔗 Dependências Utilizadas
- [Gomail](https://github.com/go-gomail/gomail) - Envio de e-mails via SMTP.
- [RabbitMQ Go Client](https://github.com/streadway/amqp) - Consumo de mensagens do RabbitMQ.
- [Azure Blob Storage SDK](https://github.com/Azure/azure-storage-blob-go) - Acesso aos templates armazenados.
- [godotenv](https://github.com/joho/godotenv) - Carregamento de variáveis de ambiente.

---

## 📜 Licença
Este projeto é distribuído sob a licença **MIT**.

