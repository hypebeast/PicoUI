package picoui

// Represents an Input element.
type Input struct {
	ui *PicoUi
	id string
}

// NewInput creates and returns a new MiUiInput.
func NewInput(ui *PicoUi, inputType string, placeholder string) *Input {
	id := "input_" + createId()
	input := Input{ui: ui, id: id}

	msg := Message{Cmd: "addinput"}
	attributes := make(map[string]interface{})
	attributes["eid"] = id
	attributes["type"] = inputType
	attributes["placeholder"] = placeholder
	msg.Attributes = attributes
	ui.handlers.enqueue(msg)

	return &input
}

func (i *Input) GetText() string {
	msg := Message{Cmd: "getinput"}
	attributes := make(map[string]interface{})
	attributes["eid"] = i.id
	msg.Attributes = attributes
	return i.ui.handlers.enqueue_and_wait_for_result(msg)
}
