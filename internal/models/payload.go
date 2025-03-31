package models

type EmailPayload struct {
	To              []string     `json:"to"`
	Subject         string       `json:"subject"`
	Body            string       `json:"body"`
	UnsubscribeLink string       `json:"unsubscribeLink"`
	Template        string       `json:"template"`
	Attachments     []Attachment `json:"attachments,omitempty"`
}
