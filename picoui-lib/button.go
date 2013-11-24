package picoui

// Represents a UI button.
type Button struct {
	ui             *PicoUi
	id             string
	onClickHandler clickHandler
}

// NewButton creates a new Button and returns it.
func NewButton(text string, ui *PicoUi, onClick clickHandler) *Button {
	id := "button_" + createId()
	button := Button{ui: ui, id: id, onClickHandler: onClick}

	msg := Message{Cmd: "addbutton"}
	attributes := make(map[string]interface{})
	attributes["eid"] = id
	attributes["txt"] = text
	msg.Attributes = attributes
	ui.handlers.enqueue(msg)

	return &button
}

// SetText sets the given text.
func (b *Button) SetText(text string) {
	msg := Message{Cmd: "updateinner"}
	attributes := make(map[string]interface{})
	attributes["eid"] = b.id
	attributes["txt"] = text
	msg.Attributes = attributes
	b.ui.handlers.enqueue(msg)
}
