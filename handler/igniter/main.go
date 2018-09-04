package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/andgate-corp/awscli-commands/slack"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/andgate-corp/awscli-commands/commands"
)

const (
	// HELP show help (not implements)
	HELP = "help"
	// EC2 EC2 commands
	EC2 = "ec2"
	// DUMMY dummy commands
	DUMMY = "dummy"
)

// Handler Lambda main handler
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var response slack.Response
	var out = &bytes.Buffer{}

	// Parse request body
	values, err := url.ParseQuery(request.Body)
	if err != nil {
		response = slack.MessageResponse{
			Text:         "Invalid Request",
			ResponseType: slack.Ephemeral,
		}
		return events.APIGatewayProxyResponse{
			Body:       response.String(),
			StatusCode: 200,
		}, nil
	}

	// Verify tokens
	token := os.Getenv("REQUEST_VERIFICATION_TOKEN")
	if token != values.Get("token") {
		response = slack.MessageResponse{
			Text:         "Invalid Request",
			ResponseType: slack.Ephemeral,
		}
		return events.APIGatewayProxyResponse{
			Body:       response.String(),
			StatusCode: 200,
		}, nil
	}

	triggerID := values.Get("trigger_id")
	argv := strings.Split(values.Get("text"), " ")

	if len(argv) < 1 {
		response = slack.MessageResponse{
			Text:         "Invalid Request",
			ResponseType: slack.Ephemeral,
		}
		return events.APIGatewayProxyResponse{
			Body:       response.String(),
			StatusCode: 200,
		}, nil
	}

	var service commands.Command

	switch argv[0] {
	case EC2:
		service = &commands.EC2Command{
			OutStream: out,
			ErrStream: out,
		}
	case DUMMY:
		service = &commands.DummyDialogCommand{
			OutStream: out,
			ErrStream: out,
		}
	}

	if service != nil {
		err := service.Run(argv[1:])

		fmt.Println(out.String())

		if err != nil {
			fmt.Println(err.Error())
			response = slack.MessageResponse{
				Text:         fmt.Sprintf("Error: %s", err.Error()),
				ResponseType: slack.Ephemeral,
			}
		} else {
			if service.GetDataType() == commands.Message {

				data, _ := service.GetData().(*slack.MessageResponse)
				data.ResponseType = slack.Ephemeral
				response = *data
			}

			if service.GetDataType() == commands.Dialog {

				data, _ := service.GetData().(*slack.DialogBody)

				if len(data.Elements) == 0 {
					response = slack.MessageResponse{
						Text: "Instance is not found.",
					}
				} else {

					requestBody := slack.OpenDialogRequest{
						TriggerID: triggerID,
						Dialog:    *data,
					}

					fmt.Println(requestBody.String())

					req, err := http.NewRequest(
						"POST",
						"https://slack.com/api/dialog.open",
						bytes.NewBuffer([]byte(requestBody.String())),
					)

					if err != nil {
						d := slack.MessageResponse{
							Text: err.Error(),
						}
						fmt.Printf("BuildError, %s", err.Error())
						return events.APIGatewayProxyResponse{
							Body:       d.String(),
							StatusCode: 200,
						}, nil
					}

					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("BOT_OAUTH_TOKEN")))

					client := &http.Client{}
					resp, err := client.Do(req)

					if err != nil {
						d := slack.MessageResponse{
							Text: err.Error(),
						}
						fmt.Printf("RequestError, %s", err.Error())
						return events.APIGatewayProxyResponse{
							Body:       d.String(),
							StatusCode: 200,
						}, nil

					}

					defer resp.Body.Close()
					body, error := ioutil.ReadAll(resp.Body)
					if error != nil {
						log.Fatal(error)
					}
					fmt.Println("[body] " + string(body))

					response = slack.MessageResponse{
						Text: "Open dialog requested.",
					}
				}
			}
		}
	} else {
		response = &slack.MessageResponse{
			Text:         fmt.Sprintf("Not Supported: %s", argv[0]),
			ResponseType: slack.Ephemeral,
		}
	}

	fmt.Println("end")

	return events.APIGatewayProxyResponse{
		Body:       response.String(),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
