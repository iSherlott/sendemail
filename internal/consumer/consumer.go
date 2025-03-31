package consumer

import (
	"sendemail/internal/service"
	"sendemail/utils"

	"github.com/streadway/amqp"
)

type Consumer struct {
	conn         *amqp.Connection
	emailService *service.EmailService
}

func NewConsumer(conn *amqp.Connection, emailService *service.EmailService) *Consumer {
	return &Consumer{conn: conn, emailService: emailService}
}

func (c *Consumer) Start(queueName string) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			utils.LogIfDevelopment("📩 Mensagem recebida: %s", msg.Body)

			if err := c.emailService.SendEmail(msg.Body); err != nil {
				utils.LogIfDevelopment("❌ Erro ao processar e-mail: %v", err)
			}
		}
	}()

	utils.LogIfDevelopment("🎧 Aguardando mensagens na fila: %s", queueName)
	<-forever

	return nil
}
