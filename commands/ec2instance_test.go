package commands

import (
	"os"
	"strings"
	"testing"
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

	t.Log(command.Result.String())
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

	t.Log(command.Result.String())
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

	t.Log(command.Result.String())
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

	t.Log(command.Result.String())
}
