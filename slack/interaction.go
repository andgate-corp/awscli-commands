package slack

import "encoding/json"

type SlackInteractionRequest struct {
	Type    string `json:"type"`
	Actions []struct {
		Name  string `json:"name"`
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"actions"`
	CallbackID      string      `json:"callback_id"`
	Team            interface{} `json:"team"`
	Channel         interface{} `json:"channel"`
	User            interface{} `json:"user"`
	ActionTs        string      `json:"action_ts"`
	MessageTs       string      `json:"message_ts"`
	AttachmentID    string      `json:"attachment_id"`
	Token           string      `json:"token"`
	ResponseURU     string      `json:"response_url"`
	TriggerID       string      `json:"trigger_id"`
	OriginalMessage interface{} `json:"original_message"`
}

type SlackInteractionResponse struct {
	ResponseType    ResponseType             `json:"response_type"`
	ReplaceOriginal bool                     `json:"replace_original"`
	Text            string                   `json:"text"`
	Attachments     []ButtonActionAttachment `json:"attachment, omitempty"`
}

func (r SlackInteractionResponse) String() string {
	b, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(b)
}
