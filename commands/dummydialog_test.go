package commands

import (
	"fmt"
	"testing"
)

func TestDummyDialog(t *testing.T) {

	c := &DummyDialogCommand{}

	c.Run([]string{})

	fmt.Print(c.GetData())
}
