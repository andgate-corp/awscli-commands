package commands

import (
	"io"

	"github.com/andgate-corp/awscli-commands/slack"
)

type DummyDialogCommand struct {
	OutStream, ErrStream io.Writer
	DialogBody           slack.DialogBody
}

func (c *DummyDialogCommand) GetDataType() DataType {
	return Dialog
}

func (c *DummyDialogCommand) GetData() interface{} {
	return c.DialogBody
}

func (c *DummyDialogCommand) Run(argv []string) error {
	dialog := slack.DialogBody{}

	dialog.Title = "Dummy dialog"
	dialog.SubmitLabel = "Submit"
	dialog.CallbackID = "dummy_callback"
	dialog.Elements = []interface{}{
		slack.DialogTextElement{
			Type:        slack.Text,
			Label:       "DummyText",
			Name:        "DummyText",
			PlaceHolder: "Dummy text",
		},
	}

	c.DialogBody = dialog

	return nil
}
