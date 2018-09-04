package commands

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/andgate-corp/awscli-commands/slack"
)

func TestDescribeInstance(t *testing.T) {

	command := EC2Command{
		OutStream: os.Stdout,
		ErrStream: os.Stderr,
	}

	input := `describe-instances -Region ap-northeast-1 -Name Indie-us-wpdev`

	err := command.Run(strings.Split(input, " "))

	if err != nil {
		t.Error(err.Error())
	}

	if command.GetDataType() != Message {
		t.Errorf("command.GetDataType() is %s, got=%s", command.GetDataType(), Message)
	}

	if command.GetData() == nil {
		t.Error("command.GetResult() returned nil.")
	}

	d, ok := command.GetData().(*slack.MessageResponse)

	if !ok {
		t.Errorf("command.GetData() is not slack.MessageResponse, got=%T", command.GetData())
	}

	t.Log(d.String())
}

func TestStartInstances(t *testing.T) {
	command := EC2Command{
		OutStream: os.Stdout,
		ErrStream: os.Stderr,
	}

	input := `start-instances -InstanceID i-0d139b4f371f0d0b7 -Region ap-northeast-1`

	err := command.Run(strings.Split(input, " "))

	if err != nil {
		t.Error(err.Error())
	}

	if command.GetDataType() != Message {
		t.Errorf("command.GetDataType() is %s, got=%s", command.GetDataType(), Message)
	}

	if command.GetData() == nil {
		t.Error("command.GetResult() returned nil.")
	}

	t.Log(command.GetData())
}

func TestStopInstances(t *testing.T) {
	command := EC2Command{
		OutStream: os.Stdout,
		ErrStream: os.Stderr,
	}

	input := `stop-instances -InstanceID i-0d139b4f371f0d0b7 -Region ap-northeast-1`

	err := command.Run(strings.Split(input, " "))

	if err != nil {
		t.Error(err.Error())
	}

	if command.GetDataType() != Message {
		t.Errorf("command.GetDataType() is %s, got=%s", command.GetDataType(), Message)
	}

	if command.GetData() == nil {
		t.Error("command.GetResult() returned nil.")
	}

	t.Log(command.GetData())
}

func TestForceStopInstances(t *testing.T) {
	command := EC2Command{
		OutStream: os.Stdout,
		ErrStream: os.Stderr,
	}

	input := `stop-instances -Force -InstanceID i-0d139b4f371f0d0b7 -Region ap-northeast-1`

	err := command.Run(strings.Split(input, " "))

	if err != nil {
		t.Error(err.Error())
	}

	if command.GetDataType() != Message {
		t.Errorf("command.GetDataType() is %s, got=%s", command.GetDataType(), Message)
	}

	if command.GetData() == nil {
		t.Error("command.GetResult() returned nil.")
	}

	t.Log(command.GetData())
}

func TestStartInstancesDialog(t *testing.T) {

	command := EC2Command{
		OutStream: os.Stdout,
		ErrStream: os.Stderr,
	}

	input := `start-instances-dialog -Region ap-northeast-1`

	err := command.Run(strings.Split(input, " "))

	if err != nil {
		t.Error(err.Error())
	}

	if command.GetDataType() != Dialog {
		t.Errorf("command.GetDataType() is %s, got=%s", command.GetDataType(), Dialog)
	}

	if command.GetData() == nil {
		t.Error("command.GetResult() returned nil.")
	}

	b, err := json.Marshal(command.GetData())

	if err != nil {
		t.Error("Json Marshal error.")
	}

	t.Log(string(b))
}

func TestStopInstancesDialog(t *testing.T) {

	command := EC2Command{
		OutStream: os.Stdout,
		ErrStream: os.Stderr,
	}

	input := `stop-instances-dialog -Region ap-northeast-1`

	err := command.Run(strings.Split(input, " "))

	if err != nil {
		t.Error(err.Error())
	}

	if command.GetDataType() != Dialog {
		t.Errorf("command.GetDataType() is %s, got=%s", command.GetDataType(), Dialog)
	}

	if command.GetData() == nil {
		t.Error("command.GetResult() returned nil.")
	}

	t.Log(command.GetData())

}
