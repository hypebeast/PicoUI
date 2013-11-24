package picoui

// Represents a Textbox element.
type TextBox struct {
	ui *PicoUi
	id string
}

// NewTextBox creates a new TextBox and returns it.
func NewTextBox(text string, element string, ui *PicoUi) *TextBox {
	id := "textbox_" + createId()
	box := TextBox{ui: ui, id: id}

	msg := Message{Cmd: "addelement"}
	attributes := make(map[string]interface{})
	attributes["e"] = element
	attributes["eid"] = id
	attributes["txt"] = text
	msg.Attributes = attributes
	ui.handlers.enqueue(msg)

	return &box
}

// SetText sets a new text to the textbox.
func (b *TextBox) SetText(text string) {
	msg := Message{Cmd: "updateinner"}
	attributes := make(map[string]interface{})
	attributes["eid"] = b.id
	attributes["txt"] = text
	msg.Attributes = attributes
	b.ui.handlers.enqueue(msg)
}
