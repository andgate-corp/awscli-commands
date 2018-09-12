package slack

import (
	"encoding/json"
	"fmt"
	"strings"
)

type InteractionRequest struct {
	Type        string            `json:"type"`
	CallbackID  string            `json:"callback_id"`
	Team        map[string]string `json:"team"`
	Channel     map[string]string `json:"channel"`
	User        map[string]string `json:"user"`
	Token       string            `json:"token"`
	ResponseURL string            `json:"response_url"`
}

type ActionInteractionRequest struct {
	Type    string `json:"type"`
	Actions []struct {
		Name  string `json:"name"`
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"actions"`
	CallbackID      string            `json:"callback_id"`
	Team            map[string]string `json:"team"`
	Channel         map[string]string `json:"channel"`
	User            map[string]string `json:"user"`
	ActionTs        string            `json:"action_ts"`
	MessageTs       string            `json:"message_ts"`
	AttachmentID    string            `json:"attachment_id"`
	Token           string            `json:"token"`
	ResponseURL     string            `json:"response_url"`
	TriggerID       string            `json:"trigger_id"`
	OriginalMessage interface{}       `json:"original_message"`
}

type InteractionResponse struct {
	ResponseType    ResponseType             `json:"response_type"`
	ReplaceOriginal bool                     `json:"replace_original"`
	Text            string                   `json:"text"`
	Attachments     []ButtonActionAttachment `json:"attachment, omitempty"`
}

func (r InteractionResponse) String() string {
	b, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(b)
}

func (r ActionInteractionRequest) String() string {
	b, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(b)
}

type DialogInteractionRequest struct {
	Type        string            `json:"type"`
	Submission  map[string]string `json:"submission"`
	CallbackID  string            `json:"callback_id"`
	State       string            `json:"state"`
	Team        map[string]string `json:"team"`
	User        map[string]string `json:"user"`
	Channel     map[string]string `json:"channel"`
	ActionTS    string            `json:"action_ts"`
	Token       string            `json:"token"`
	ResponseURL string            `json:"response_url"`
}

func (r DialogInteractionRequest) String() string {
	b, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(b)
}

func (r DialogInteractionRequest) GetArgs() []string {

	if len(r.State) == 0 {
		return []string{}
	}

	args := strings.Split(r.State, " ")

	if len(r.Submission) > 0 {
		for k, v := range r.Submission {
			args = append(args, fmt.Sprintf("-%s", k))
			args = append(args, v)
		}
	}

	return args

}
