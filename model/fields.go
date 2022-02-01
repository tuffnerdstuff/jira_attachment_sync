package model

type Fields struct {
	Attachments []Attachment `json:"attachment"`
	Summary     string       `json:"summary"`
	Description string       `json:"description"`
}
