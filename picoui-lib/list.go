package picoui

// Represents a List element.
type List struct {
	ui   *PicoUi
	page *Page
	id   string
}

// Represents a list item.
type ListItem struct {
	ui              *PicoUi
	id              string
	parentId        string
	onClickHandler  clickHandler
	onToggleHandler toggleHandler
	toggleId        string
}

func NewList(ui *PicoUi, page *Page) *List {
	id := "list_" + createId()
	list := List{ui: ui, page: page, id: id}

	msg := Message{Cmd: "addlist"}
	attributes := make(map[string]interface{})
	attributes["eid"] = id
	msg.Attributes = attributes
	ui.handlers.enqueue(msg)

	return &list
}

// AddItem adds a ListItem to the list.
func (l *List) AddItem(text string, chevron bool, toggle bool, onClick clickHandler, onToggle toggleHandler) *ListItem {
	item := NewListItem(l.ui, l.id, text, chevron, toggle, onClick, onToggle)
	l.page.elements = append(l.page.elements, item)
	if onClick != nil {
		l.page.clickables[item.id] = onClick
	}
	if onToggle != nil {
		l.page.toggables[item.toggleId] = onToggle
	}

	return item
}

// NewListItem creates a new ListItem element.
func NewListItem(ui *PicoUi, parentId string, text string, chevron bool, toggle bool, onClick clickHandler, onToggle toggleHandler) *ListItem {
	id := "listitem_" + createId()
	tg_id := "tg_" + createId()
	item := ListItem{
		ui:              ui,
		id:              id,
		parentId:        parentId,
		onClickHandler:  onClick,
		onToggleHandler: onToggle,
		toggleId:        tg_id}

	msg := Message{Cmd: "addlistitem"}
	attributes := make(map[string]interface{})
	attributes["eid"] = id
	attributes["pid"] = parentId
	attributes["txt"] = text
	attributes["chevron"] = chevron
	attributes["toggle"] = toggle
	attributes["tid"] = tg_id
	msg.Attributes = attributes
	ui.handlers.enqueue(msg)

	return &item
}
