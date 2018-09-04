package slack

import (
	"encoding/json"
)

// OpenDialogRequest SlackAPI request data for `dialog.open`
type OpenDialogRequest struct {
	TriggerID string     `json:"trigger_id"`
	Dialog    DialogBody `json:"dialog"`
}

func (t OpenDialogRequest) String() string {
	b, err := json.Marshal(t)

	if err != nil {
		return ""
	}

	return string(b)
}

// DialogBody Dialog body
type DialogBody struct {
	CallbackID     string        `json:"callback_id"`
	Title          string        `json:"title"`
	SubmitLabel    string        `json:"submit_label"`
	NotifyOnCancel bool          `json:"notify_on_cancel"`
	State          string        `json:"state"`
	Elements       []interface{} `json:"elements"`
}

// DialogElementType Slack Dialog element type
type DialogElementType int

const (
	// Text TEXT input field
	Text DialogElementType = iota
	// TextArea TEXTAREA input field
	TextArea
	// Select SELECT input field
	Select
)

func (t DialogElementType) String() string {
	switch t {
	case Text:
		return "text"
	case TextArea:
		return "textarea"
	case Select:
		return "select"
	default:
		return ""
	}
}

// MarshalJSON JSON Marshalize interface method
func (t DialogElementType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}

// TextSubtype Text and Textarea format type
type TextSubtype int

const (
	// Email Email format
	Email TextSubtype = iota
	// Number Number format
	Number
	// Tel Tel-Number format
	Tel
	// URL URL format
	URL
)

func (t TextSubtype) String() string {
	switch t {
	case Email:
		return "email"
	case Number:
		return "number"
	case Tel:
		return "tel"
	case URL:
		return "url"
	default:
		return ""
	}
}

// MarshalJSON JSON marshalize interface method
func (t TextSubtype) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}

// DialogSelectElement SELECT element in Slack Dialog
type DialogSelectElement struct {
	Label   string                      `json:"label"`
	Type    DialogElementType           `json:"type"`
	Name    string                      `json:"name"`
	Options []DialogSelectElementOption `json:"options"`
}

// DialogSelectElementOption OPTION element for DialogSelectElement
type DialogSelectElementOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// AddOption Add option to select element
func (elem *DialogSelectElement) AddOption(label string, value string) {
	elem.Options = append(elem.Options, DialogSelectElementOption{Label: label, Value: value})
}

// AddSelectElement Create element for select
func (d *DialogBody) AddSelectElement(name string, label string) *DialogSelectElement {
	elem := &DialogSelectElement{
		Type:  Select,
		Name:  name,
		Label: label,
	}

	d.Elements = append(d.Elements, elem)

	return elem
}

// AddTextElement Create element for text
func (d *DialogBody) AddTextElement(name string, label string) *DialogTextElement {
	elem := &DialogTextElement{
		Type:  Text,
		Name:  name,
		Label: label,
	}

	d.Elements = append(d.Elements, elem)

	return elem
}

// AddTextareaElement Create element for textarea
func (d *DialogBody) AddTextareaElement(name string, label string) *DialogTextElement {
	elem := &DialogTextElement{
		Type:  TextArea,
		Name:  name,
		Label: label,
	}

	d.Elements = append(d.Elements, elem)

	return elem
}

// DialogTextElement Slack Dialog text and textarea field element
type DialogTextElement struct {
	Label       string            `json:"label"`
	Name        string            `json:"name"`
	Type        DialogElementType `json:"type"`
	MaxLength   int               `json:"max_length,omitempty"`
	MinLength   int               `json:"min_length,omitempty"`
	Optional    bool              `json:"optional"`
	Hint        string            `json:"hint,omitempty"`
	SubType     TextSubtype       `json:"subtype,omitempty"`
	Value       string            `json:"value,omitempty"`
	PlaceHolder string            `json:"placeholder,omitempty"`
}
