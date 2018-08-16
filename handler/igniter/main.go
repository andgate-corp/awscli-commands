package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/t_sumisaki/awscli-commands/commands"
)

type ResponseType int

const (
	Ephemeral ResponseType = iota
	InChannel
)

const (
	HELP = "help"
	EC2  = "ec2"
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
func (t ResponseType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}

type SlackResponse struct {
	ResponseType ResponseType  `json:"response_type"`
	Text         string        `json:"text"`
	Attachments  []interface{} `json:"attachments, omitempty"`
}

func (r *SlackResponse) String() string {
	b, err := json.Marshal(r)

	if err != nil {
		return ""
	}

	fmt.Printf(string(b))

	return string(b)
}

type Command interface {
	Run([]string) error
	GetResult() commands.CommandResult
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var response SlackResponse
	var out = &bytes.Buffer{}

	// Parse request body
	values, err := url.ParseQuery(request.Body)
	if err != nil {
		response = SlackResponse{
			Text:         "Invalid Request",
			ResponseType: Ephemeral,
		}
		return events.APIGatewayProxyResponse{
			Body:       response.String(),
			StatusCode: 200,
		}, nil
	}

	// Verify tokens
	token := os.Getenv("REQUEST_VERIFICATION_TOKEN")
	if token != values.Get("token") {
		response = SlackResponse{
			Text:         "Invalid Request",
			ResponseType: Ephemeral,
		}
		return events.APIGatewayProxyResponse{
			Body:       response.String(),
			StatusCode: 200,
		}, nil
	}

	fmt.Printf(values.Get("text"))
	argv := strings.Split(values.Get("text"), " ")

	if len(argv) < 1 {
		response = SlackResponse{
			Text:         "Invalid Request",
			ResponseType: Ephemeral,
		}
		return events.APIGatewayProxyResponse{
			Body:       response.String(),
			StatusCode: 200,
		}, nil
	}

	var service Command

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
			response = SlackResponse{
				Text:         fmt.Sprintf("Error: %s", err.Error()),
				ResponseType: Ephemeral,
			}
		} else {

			result := service.GetResult()
			response = SlackResponse{
				Text:         result.Text,
				ResponseType: Ephemeral,
			}

			if len(result.Attachments) > 0 {
				response.Attachments = result.Attachments
			}
		}
	} else {
		response = SlackResponse{
			Text:         fmt.Sprintf("Not Supported: %s", argv[0]),
			ResponseType: Ephemeral,
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
