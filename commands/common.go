package commands

import (
	"encoding/json"
)

// Command コマンド用共通インターフェース
type Command interface {
	Run([]string) error
	GetResult() CommandResult
}

// CommandResult Commandインターフェースの結果出力用構造体
type CommandResult struct {
	Text        string
	Color       string
	Attachments []interface{}
}

func (r CommandResult) String() string {
	b, err := json.Marshal(r)

	if err != nil {
		return ""
	}

	return string(b)
}

// ButtonActionAttachment ButtonActions用のAttachment構造体
type ButtonActionAttachment struct {
	Text       string             `json:"text"`
	Fallback   string             `json:"fallback"`
	CallbackID string             `json:"callback_id"`
	Actions    []ButtonActionItem `json:"actions"`
	Fields     []AttachmentField  `json:"fields, omitempty"`
}

// ButtonActionItem Action用ボタン情報の構造体
type ButtonActionItem struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Text  string `json:"text"`
	Value string `json:"value"`
	Style string `json:"style, omitempty"`
}

// AttachmentField Attachmentの汎用フィールド構造体
type AttachmentField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short, omitempty"`
}
