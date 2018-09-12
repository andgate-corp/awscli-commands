package slack

import (
	"encoding/json"
	"fmt"
)

// Response Slackリクエスト返信用Interface
type Response interface {
	String() string
}

// ResponseType 返信時のメッセージ表示タイプ
type ResponseType int

const (
	// Ephemeral 発信者のみ
	Ephemeral ResponseType = iota
	// InChannel チャンネル全体に表示
	InChannel
)

func (t ResponseType) String() string {
	switch t {
	case Ephemeral:
		return "ephemeral"
	case InChannel:
		return "in_channel"
	default:
		return ""
	}
}

// MarshalJSON JSONパース用
func (t ResponseType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}

// MessageResponse テキストメッセージ用レスポンス
type MessageResponse struct {
	ResponseType ResponseType  `json:"response_type"`
	Text         string        `json:"text"`
	Channel      string        `json:"channel, omitempty"`
	Attachments  []interface{} `json:"attachments, omitempty"`
}

func (r MessageResponse) String() string {
	b, err := json.Marshal(r)

	if err != nil {
		return ""
	}

	fmt.Printf(string(b))

	return string(b)
}

// NullResponse メッセージ返答が必要ない場合のResponse
type NullResponse struct {
}

func (r NullResponse) String() string {
	return ""
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
