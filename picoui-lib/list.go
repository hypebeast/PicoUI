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
func (l *List) AddItem(text string, chevron bool, onClick clickHandler) *ListItem {
	item := newListItem(l.ui, l.id, text, chevron, onClick)
	l.page.elements = append(l.page.elements, item)
	if onClick != nil {
		l.page.clickables[item.id] = onClick
	}

	return item
}

func (l *List) AddToggle(text string, onToggle toggleHandler) *ListItem {
	id := "toggleitem_" + createId()
	tg_id := "tg_" + createId()
	item := ListItem{
		ui:              l.ui,
		id:              id,
		parentId:        l.id,
		toggleId:        tg_id,
		onToggleHandler: onToggle}

	if onToggle != nil {
		l.page.toggables[tg_id] = onToggle
	}

	msg := Message{Cmd: "addtoggleitem"}
	attributes := make(map[string]interface{})
	attributes["eid"] = id
	attributes["pid"] = l.id
	attributes["txt"] = text
	attributes["tid"] = tg_id
	msg.Attributes = attributes
	l.ui.handlers.enqueue(msg)

	return &item
}

func (l *List) AddCheckbox(text string, onToggle toggleHandler) *ListItem {
	id := "toggleitem_" + createId()
	tg_id := "tg_" + createId()
	item := ListItem{
		ui:              l.ui,
		id:              id,
		parentId:        l.id,
		toggleId:        tg_id,
		onToggleHandler: onToggle}

	if onToggle != nil {
		l.page.toggables[tg_id] = onToggle
	}

	msg := Message{Cmd: "addcheckboxitem"}
	attributes := make(map[string]interface{})
	attributes["eid"] = id
	attributes["pid"] = l.id
	attributes["txt"] = text
	attributes["tid"] = tg_id
	msg.Attributes = attributes
	l.ui.handlers.enqueue(msg)

	return &item
}

func (l *List) AddRange() {
	// TODO
}

// NewListItem creates a new ListItem element.
func newListItem(ui *PicoUi, parentId string, text string, chevron bool, onClick clickHandler) *ListItem {
	id := "listitem_" + createId()
	item := ListItem{
		ui:             ui,
		id:             id,
		parentId:       parentId,
		onClickHandler: onClick}

	msg := Message{Cmd: "addlistitem"}
	attributes := make(map[string]interface{})
	attributes["eid"] = id
	attributes["pid"] = parentId
	attributes["txt"] = text
	attributes["chevron"] = chevron
	msg.Attributes = attributes
	ui.handlers.enqueue(msg)

	return &item
}

func (l *ListItem) SetText(text string) {
	// TODO
}
