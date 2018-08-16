package commands

import (
	"encoding/json"
)

type ResultType int

const (
	Message ResultType = iota
	ButtonAction
)

type CommandResult struct {
	Text        string
	ResultType  ResultType
	Color       string
	Attachments []interface{}
}

type AttachmentColor string

const ()

func (r CommandResult) String() string {
	b, err := json.Marshal(r)

	if err != nil {
		return ""
	}

	return string(b)
}

type ButtonActionAttachment struct {
	Title      string             `json:"title"`
	Text       string             `json:"text"`
	Fallback   string             `json:"fallback"`
	CallbackID string             `json:"callback_id"`
	Actions    []ButtonActionItem `json:"actions"`
	Fields     []AttachmentField  `json:"fields, omitempty"`
}

type ButtonActionItem struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Text  string `json:"text"`
	Value string `json:"value"`
	Style string `json:"style, omitempty"`
}

type AttachmentField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short, omitempty"`
}
