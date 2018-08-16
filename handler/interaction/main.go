package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/t_sumisaki/awscli-commands/commands"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	HELP = "help"
	EC2  = "ec2"
)

type SlackInteractionRequest struct {
	Type    string `json:"type"`
	Actions []struct {
		Name  string `json:"name"`
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"actions"`
	CallbackId      string      `json:"callback_id"`
	Team            interface{} `json:"team"`
	Channel         interface{} `json:"channel"`
	User            interface{} `json:"user"`
	ActionTs        string      `json:"action_ts"`
	MessageTs       string      `json:"message_ts"`
	AttachmentId    string      `json:"attachment_id"`
	Token           string      `json:"token"`
	ResponseUrl     string      `json:"response_url"`
	TriggerId       string      `json:"trigger_id"`
	OriginalMessage interface{} `json:"original_message"`
}

type ResponseType string

const (
	Ephemeral ResponseType = "ephemeral"
	InChannel ResponseType = "in_channel"
)

type SlackInteractionResponse struct {
	ResponseType    ResponseType                      `json:"response_type"`
	ReplaceOriginal bool                              `json:"replace_original"`
	Text            string                            `json:"text"`
	Attachments     []commands.ButtonActionAttachment `json:"attachment, omitempty"`
}

func (r SlackInteractionResponse) String() string {
	b, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(b)
}

type Command interface {
	Run([]string) error
	GetResult() commands.CommandResult
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var response SlackInteractionResponse

	values, err := url.ParseQuery(request.Body)

	fmt.Println(request.Body)

	token := os.Getenv("REQUEST_VERIFICATION_TOKEN")

	if err != nil {
		response = SlackInteractionResponse{
			Text:            "Invalid Request",
			ReplaceOriginal: true,
			ResponseType:    Ephemeral,
		}
	} else {
		var req SlackInteractionRequest
		err := json.Unmarshal([]byte(values.Get("payload")), &req)

		if err != nil {
			fmt.Println(err.Error())
			response = SlackInteractionResponse{
				Text:            "Invalid Request",
				ReplaceOriginal: true,
				ResponseType:    Ephemeral,
			}
		} else if token != req.Token {
			fmt.Printf("Token error, got %s\n", values.Get("token"))
			response = SlackInteractionResponse{
				Text:            "Invalid Request",
				ReplaceOriginal: true,
				ResponseType:    Ephemeral,
			}
		} else {

			var service Command
			var out = &bytes.Buffer{}

			argv := strings.Split(req.Actions[0].Value, " ")
			switch argv[0] {
			case EC2:
				service = &commands.EC2Command{
					OutStream: out,
					ErrStream: out,
				}
			}

			if service != nil {
				err := service.Run(argv[1:])

				fmt.Println(out.String())

				if err != nil {
					fmt.Println(err.Error())
					response = SlackInteractionResponse{
						Text:            fmt.Sprintf("Error: %s", err.Error()),
						ResponseType:    Ephemeral,
						ReplaceOriginal: true,
					}
				} else {
					result := service.GetResult()
					response = SlackInteractionResponse{
						Text:            fmt.Sprintf("Success: %s", result.Text),
						ResponseType:    Ephemeral,
						ReplaceOriginal: true,
					}
				}
			} else {
				response = SlackInteractionResponse{
					Text:            fmt.Sprintf("Not Supported: %s", argv[0]),
					ResponseType:    Ephemeral,
					ReplaceOriginal: true,
				}
			}
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
