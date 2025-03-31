package main

import (
	"os"

	"sendemail/internal/consumer"
	"sendemail/internal/service"
	"sendemail/pkg/rabbitmq"
	"sendemail/utils"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		utils.LogIfDevelopment("⚠️ Arquivo .env não encontrado, usando valores padrão")
	}

	conn, err := rabbitmq.Connect()
	if err != nil {
		utils.FatalIf(err, "❌ Erro ao conectar ao RabbitMQ: %v", err)
	}
	defer conn.Close()

	emailService := service.NewEmailService()

	queueName := os.Getenv("RABBITMQ_QUEUE")
	consumer := consumer.NewConsumer(conn, emailService)

	if err := consumer.Start(queueName); err != nil {
		utils.FatalIf(err, "❌ Variável de ambiente RABBITMQ_QUEUE não definida")
	}
}
