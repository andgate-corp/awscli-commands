package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/andgate-corp/awscli-commands/commands"
	"github.com/andgate-corp/awscli-commands/slack"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	HELP = "help"
	EC2  = "ec2"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var response slack.Response

	values, err := url.ParseQuery(request.Body)

	if err != nil {
		fmt.Printf("Parse error, %s\n", err.Error())
		response = slack.SlackInteractionResponse{
			Text:            "Parse Error",
			ReplaceOriginal: true,
			ResponseType:    slack.Ephemeral,
		}
		return events.APIGatewayProxyResponse{
			Body:       response.String(),
			StatusCode: 200,
		}, nil
	}

	fmt.Println(request.Body)

	token := os.Getenv("REQUEST_VERIFICATION_TOKEN")

	if err != nil {
		response = slack.SlackInteractionResponse{
			Text:            "Invalid Request",
			ReplaceOriginal: true,
			ResponseType:    slack.Ephemeral,
		}
	} else {
		var req slack.SlackInteractionRequest
		err := json.Unmarshal([]byte(values.Get("payload")), &req)

		if err != nil {
			fmt.Println(err.Error())
			response = slack.SlackInteractionResponse{
				Text:            "Invalid Request",
				ReplaceOriginal: true,
				ResponseType:    slack.Ephemeral,
			}
		} else if token != req.Token {
			fmt.Printf("Token error, got %s\n", values.Get("token"))
			response = slack.SlackInteractionResponse{
				Text:            "Invalid Request",
				ReplaceOriginal: true,
				ResponseType:    slack.Ephemeral,
			}
		} else {

			fmt.Printf("Payload Type: %s\n", req.Type)
			switch req.Type {
			case "dialog_submission":
				response = &slack.MessageResponse{
					Text: "Dialog request received.",
				}
			case "interactive_message":
				response = ExecuteInteractiveMessages(req)
			default:
				response = slack.SlackInteractionResponse{
					Text:            fmt.Sprintf("%s is not supported.", req.Type),
					ReplaceOriginal: true,
					ResponseType:    slack.Ephemeral,
				}
			}
		}
	}

	fmt.Printf(response.String())

	return events.APIGatewayProxyResponse{
		Body:       response.String(),
		StatusCode: 200,
	}, nil
}

func ExecuteDialog(payload slack.SlackInteractionRequest) slack.SlackInteractionResponse {
	var response = slack.SlackInteractionResponse{
		Text:            "Dummy response",
		ResponseType:    slack.Ephemeral,
		ReplaceOriginal: false,
	}

	return response
}

func ExecuteInteractiveMessages(payload slack.SlackInteractionRequest) slack.SlackInteractionResponse {

	var response slack.SlackInteractionResponse

	var service commands.Command
	var out = &bytes.Buffer{}

	argv := strings.Split(payload.Actions[0].Value, " ")
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
			response = slack.SlackInteractionResponse{
				Text:            fmt.Sprintf("Error: %s", err.Error()),
				ResponseType:    slack.Ephemeral,
				ReplaceOriginal: true,
			}
		} else {
			result, _ := service.GetData().(slack.MessageResponse)
			response = slack.SlackInteractionResponse{
				Text:            fmt.Sprintf("Success: %s", result.Text),
				ResponseType:    slack.Ephemeral,
				ReplaceOriginal: true,
			}
		}
	} else {
		response = slack.SlackInteractionResponse{
			Text:            fmt.Sprintf("Not Supported: %s", argv[0]),
			ResponseType:    slack.Ephemeral,
			ReplaceOriginal: true,
		}
	}
	return response
}

func main() {
	lambda.Start(Handler)
}
