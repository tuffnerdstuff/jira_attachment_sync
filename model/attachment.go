package model

import (
	"fmt"
	"strings"
	"time"
)

const ENDPOINT_ATTACHMENT = "attachment"

type Attachment struct {
	ID          string `json:"id"`
	URL         string `json:"content"`
	Filename    string `json:"filename"`
	MimeType    string `json:"mimeType"`
	CreatedTime string `json:"created"`
	Size        int    `json:"size"`
}

func (a *Attachment) GetFilenameWithDatePrefix() string {
	prefix := a.ID
	// golang is whack! These are "magic numbers" in the pattern string ...
	createdTime, err := time.Parse("2006-01-02T15:04:05.000-0700", a.CreatedTime)
	if err == nil {
		prefix = createdTime.Format("2006-01-02")
	}
	return fmt.Sprintf("%s_%s", prefix, a.Filename)
}

func (a *Attachment) IsCompressed() bool {

	extensions := []string{".zip", ".rar", ".7z", ".001"}
	for _, extension := range extensions {
		if strings.HasSuffix(strings.ToLower(a.Filename), extension) {
			return true
		}

	}
	return false
}
