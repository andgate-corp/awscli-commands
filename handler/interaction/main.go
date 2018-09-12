package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
		response = slack.InteractionResponse{
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
		response = slack.InteractionResponse{
			Text:            "Invalid Request",
			ReplaceOriginal: true,
			ResponseType:    slack.Ephemeral,
		}
	} else {

		payload := values.Get("payload")
		var req slack.InteractionRequest
		err := json.Unmarshal([]byte(payload), &req)

		if err != nil {
			fmt.Println(err.Error())
			response = slack.InteractionResponse{
				Text:            "Invalid Request",
				ReplaceOriginal: true,
				ResponseType:    slack.Ephemeral,
			}
		}

		fmt.Print(req)

		if token != req.Token {
			fmt.Printf("Token error, got %s\n", values.Get("token"))
			response = slack.InteractionResponse{
				Text:            "Invalid Request",
				ReplaceOriginal: true,
				ResponseType:    slack.Ephemeral,
			}
		} else {

			fmt.Printf("Payload Type: %s\n", req.Type)
			switch req.Type {
			case "dialog_submission":
				response = ExecuteDialog(payload)
			case "interactive_message":
				response = ExecuteInteractiveMessages(payload)
			default:
				response = slack.InteractionResponse{
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

func ExecuteDialog(payload string) slack.Response {
	var (
		service  commands.Command
		request  slack.DialogInteractionRequest
		response = slack.NullResponse{}
		out      = &bytes.Buffer{}
	)

	err := json.Unmarshal([]byte(payload), &request)

	if err != nil {
		fmt.Println(err.Error())
		return response
	}

	argv := request.GetArgs()

	switch argv[0] {
	case EC2:
		service = &commands.EC2Command{
			OutStream: out,
			ErrStream: out,
		}
	default:
		fmt.Printf("Command is not supported. => %s\n", argv[0])
	}

	if service != nil {
		err := service.Run(argv[1:])

		if err != nil {
			fmt.Println(err.Error())
			// TODO: Send Error Message
		} else {
			data, _ := service.GetData().(*slack.MessageResponse)
			data.Channel = request.Channel["id"]

			req, err := http.NewRequest(
				"POST",
				"https://slack.com/api/chat.postMessage",
				bytes.NewBuffer([]byte(data.String())),
			)

			if err != nil {
				fmt.Println(err.Error())
				return response
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("BOT_OAUTH_TOKEN")))

			client := &http.Client{}
			resp, err := client.Do(req)

			if err != nil {
				fmt.Println(err.Error())
				return response
			}

			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err.Error())
				return response
			}

			fmt.Printf("[body] %s\n", string(body))

		}
	}

	return response
}

func ExecuteInteractiveMessages(payload string) slack.Response {

	var (
		request  slack.ActionInteractionRequest
		response slack.InteractionResponse
		service  commands.Command
		out      = &bytes.Buffer{}
	)

	err := json.Unmarshal([]byte(payload), &request)

	if err != nil {
		response.Text = fmt.Sprintf("Error: %s", err.Error())

		return response
	}

	argv := strings.Split(request.Actions[0].Value, " ")
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
			response.Text = fmt.Sprintf("Error: %s", err.Error())
		} else {
			result, _ := service.GetData().(slack.MessageResponse)
			response.Text = fmt.Sprintf("Success: %s", result.Text)
		}
	} else {
		response.Text = fmt.Sprintf("Not Supported: %s", argv[0])
	}
	return response
}

func main() {
	lambda.Start(Handler)
}
