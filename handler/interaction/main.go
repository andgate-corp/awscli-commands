package main

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type SlackInteractionRequest struct {
	Type            string        `json:"type"`
	Actions         []interface{} `json:"actions"`
	CallbackId      string        `json:"callback_id"`
	Team            interface{}   `json:"team"`
	Channel         interface{}   `json:"channel"`
	User            interface{}   `json:"user"`
	ActionTs        string        `json:"action_ts"`
	MessageTs       string        `json:"message_ts"`
	AttachmentId    string        `json:"attachment_id"`
	Token           string        `json:"token"`
	ResponseUrl     string        `json:"response_url"`
	TriggerId       string        `json:"trigger_id"`
	OriginalMessage interface{}   `json:"original_message"`
}

type ResponseType string

const (
	Ephemeral ResponseType = "ephemeral"
	InChannel ResponseType = "in_channel"
)

type SlackInteractionResponse struct {
	ResponseType    ResponseType `json:"response_type"`
	ReplaceOriginal bool         `json:"replace_original"`
	Text            string       `json:"text"`
}

func (r SlackInteractionResponse) String() string {
	b, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(b)
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var response SlackInteractionResponse

	values, err := url.ParseQuery(request.Body)

	if err != nil {
		response = SlackInteractionResponse{
			Text:            "Error",
			ReplaceOriginal: false,
			ResponseType:    Ephemeral,
		}
	} else {
		fmt.Println(values)
		response = SlackInteractionResponse{
			Text:            "Success",
			ReplaceOriginal: false,
			ResponseType:    Ephemeral,
		}
	}

	return events.APIGatewayProxyResponse{
		Body:       response.String(),
		StatusCode: 200,
	}, nil
}
func main() {
	lambda.Start(Handler)
}
