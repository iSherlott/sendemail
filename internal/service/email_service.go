package service

import (
	"bytes"
	"encoding/json"
	"html/template"

	"sendemail/internal/models"
	"sendemail/pkg/azure"
	"sendemail/pkg/email"
	"sendemail/utils"
)

type EmailService struct{}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (s *EmailService) SendEmail(message []byte) error {
	var emailPayload models.EmailPayload

	if err := json.Unmarshal(message, &emailPayload); err != nil {
		utils.LogIfDevelopment("❌ Erro ao decodificar a mensagem: %v", err)
		return err
	}

	templateContent, err := azure.GetTemplateFromBlob(emailPayload.Template + ".html")
	if err != nil {
		utils.LogIfDevelopment("❌ Erro ao recuperar o template: %v", err)
		return err
	}

	emailBody, err := s.renderTemplate(templateContent, emailPayload)
	if err != nil {
		utils.LogIfDevelopment("❌ Erro ao renderizar o template: %v", err)
		return err
	}

	return email.SendEmail(emailPayload.To, emailPayload.Subject, emailBody, emailPayload.Attachments)
}

func (s *EmailService) renderTemplate(templateContent string, payload models.EmailPayload) (string, error) {
	tmpl, err := template.New("email").Parse(templateContent)
	if err != nil {
		return "", err
	}

	var renderedContent bytes.Buffer
	err = tmpl.Execute(&renderedContent, payload)
	if err != nil {
		return "", err
	}

	return renderedContent.String(), nil
}
