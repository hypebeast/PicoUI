package picoui

// Represents a Button element.
type Button struct {
	ui             *PicoUi
	id             string
	onClickHandler clickHandler
}

// NewButton creates a new Button and returns it.
func newButton(ui *PicoUi, text string, classAttributes []string, icon string, onClick clickHandler) *Button {
	id := "button_" + createId()
	button := Button{ui: ui, id: id, onClickHandler: onClick}

	msg := Message{Cmd: "addbutton"}
	attributes := make(map[string]interface{})
	attributes["eid"] = id
	attributes["txt"] = text
	if classAttributes != nil {
		attributes["classAttr"] = classAttributes
	}
	attributes["icon"] = icon
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

func (b *Button) SetClassAttributes(attr []string) {
	msg := Message{Cmd: "updateClassAttr"}
	attributes := make(map[string]interface{})
	attributes["eid"] = b.id
	attributes["classAttr"] = attr
	msg.Attributes = attributes
	b.ui.handlers.enqueue(msg)
}

func (b *Button) SetIcon(icon string) {
	msg := Message{Cmd: "setIcon"}
	attributes := make(map[string]interface{})
	attributes["eid"] = b.id
	attributes["icon"] = icon
	msg.Attributes = attributes
	b.ui.handlers.enqueue(msg)
}
