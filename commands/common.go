package commands

// Command コマンド用共通インターフェース
type Command interface {
	Run([]string) error
	GetDataType() DataType
	GetData() interface{}
}

type DataType int

const (
	None DataType = iota
	Message
	Dialog
)

func (t DataType) String() string {
	switch t {
	case Message:
		return "message"
	case Dialog:
		return "dialog"
	default:
		return "undefined"
	}
}
