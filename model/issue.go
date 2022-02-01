package model

import "fmt"

const ENDPOINT_ISSUE = "issue"

type Issue struct {
	ID     string `json:"id"`
	Key    string `json:"key"`
	Fields Fields `json:"fields"`
}

func (i *Issue) GetTitle() string {
	return fmt.Sprintf("%s %s", i.Key, i.Fields.Summary)
}
