package email

import (
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"strconv"

	"sendemail/internal/models"
	"sendemail/utils"

	"gopkg.in/gomail.v2"
)

func SendEmail(to []string, subject, body string, attachments []models.Attachment) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_EMAIL"))
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	err := AttachFiles(m, attachments)
	if err != nil {
		utils.LogIfDevelopment("❌ Erro ao anexar arquivos: %v", err)
		return err
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	smtpUser := os.Getenv("SMTP_EMAIL")
	smtpPass := os.Getenv("SMTP_PASSWORD")

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	if err := d.DialAndSend(m); err != nil {
		utils.LogIfDevelopment("❌ Erro ao enviar e-mail: %v", err)
		return err
	}

	utils.LogIfDevelopment("✅ E-mail enviado para: %v", to)
	return nil
}

func AttachFiles(m *gomail.Message, attachments []models.Attachment) error {
	for _, att := range attachments {
		if att.Content != "" {
			decodedContent, err := base64.StdEncoding.DecodeString(att.Content)
			if err != nil {
				utils.LogIfDevelopment("❌ Erro ao decodificar anexo %s: %v", att.Filename, err)
				continue
			}
			m.Attach(att.Filename, gomail.SetCopyFunc(func(w io.Writer) error {
				_, err := w.Write(decodedContent)
				return err
			}))
		} else if att.URL != "" {
			resp, err := http.Get(att.URL)
			if err != nil {
				utils.LogIfDevelopment("❌ Erro ao baixar anexo %s: %v", att.Filename, err)
				continue
			}
			defer resp.Body.Close()

			m.Attach(att.Filename, gomail.SetCopyFunc(func(w io.Writer) error {
				_, err := io.Copy(w, resp.Body)
				return err
			}))
		}
	}
	return nil
}
