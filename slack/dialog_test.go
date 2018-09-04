package slack

import "testing"

func TestAddTextElement(t *testing.T) {
	input := &OpenDialogInput{}

	elem := input.AddTextElement()

	if elem.Type != Text {
		t.Errorf("result.Type should be equals %s, got=%s", Text, elem.Type)
	}
}

func TestAddTestareaElement(t *testing.T) {
	input := &OpenDialogInput{}

	elem := input.AddTextareaElement()

	if elem.Type != TextArea {
		t.Errorf("result.Type should be equals %s, got=%s", TextArea, elem.Type)
	}
}

func TestAddSelectElement(t *testing.T) {
	input := &OpenDialogInput{}

	elem := input.AddSelectElement()

	if elem.Type != Select {
		t.Errorf("result.Type should be equals %s, got=%s", Select, elem.Type)
	}
}
