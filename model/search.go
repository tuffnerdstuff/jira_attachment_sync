package model

type Search struct {
	Offset       int     `json:"startAt"`
	MaxResults   int     `json:"maxResults"`
	TotalResults int     `jason:"total"`
	Issues       []Issue `json:"issues"`
}
