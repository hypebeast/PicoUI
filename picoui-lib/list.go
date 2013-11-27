package picoui

import (
	"strconv"
)

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

// Represents an Input element.
type Input struct {
	ui       *PicoUi
	id       string
	parentId string
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
func (l *List) AddItem(text string, leftIcon string, rightIcon string, onClick clickHandler) *ListItem {
	item := newListItem(l.ui, l.id, text, leftIcon, rightIcon, onClick)
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

	l.page.elements = append(l.page.elements, item)

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

	l.page.elements = append(l.page.elements, item)

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

func (l *List) AddRange(min int, max int, iconLeft string, iconRight string, onSlide slideHandler) *Range {
	id := "range_" + createId()
	slide_id := "slide_" + createId()
	item := Range{ui: l.ui, page: l.page, id: id, slideId: slide_id, onSlide: onSlide}

	if onSlide != nil {
		l.page.slideHandlers[slide_id] = onSlide
	}

	l.page.elements = append(l.page.elements, item)

	msg := Message{Cmd: "addrangeitem"}
	attributes := make(map[string]interface{})
	attributes["eid"] = id
	attributes["pid"] = l.id
	attributes["slideid"] = slide_id
	attributes["min"] = strconv.Itoa(min)
	attributes["max"] = strconv.Itoa(max)
	attributes["iconleft"] = iconLeft
	attributes["iconright"] = iconRight
	msg.Attributes = attributes
	l.ui.handlers.enqueue(msg)

	return &item
}

func (l *List) AddDivider(text string) {
	id := "divider_" + createId()
	msg := Message{Cmd: "adddivider"}
	attributes := make(map[string]interface{})
	attributes["id"] = id
	attributes["pid"] = l.id
	attributes["txt"] = text
	msg.Attributes = attributes
	l.ui.handlers.enqueue(msg)
}

func (l *List) AddInput(inputType string, placeholder string) *Input {
	item := newInput(l.ui, l.id, inputType, placeholder)
	l.page.elements = append(l.page.elements, item)

	return item
}

// NewListItem creates a new ListItem element and returns it.
func newListItem(ui *PicoUi, parentId string, text string, leftIcon string, rightIcon string, onClick clickHandler) *ListItem {
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
	attributes["lefticon"] = leftIcon
	attributes["righticon"] = rightIcon
	msg.Attributes = attributes
	ui.handlers.enqueue(msg)

	return &item
}

func (l *ListItem) SetText(text string) {
	// TODO
}

// NewInput creates and returns a new MiUiInput.
func newInput(ui *PicoUi, parentId string, inputType string, placeholder string) *Input {
	id := "input_" + createId()
	input := Input{ui: ui, id: id, parentId: parentId}

	msg := Message{Cmd: "addinputitem"}
	attributes := make(map[string]interface{})
	attributes["eid"] = id
	attributes["type"] = inputType
	attributes["placeholder"] = placeholder
	attributes["pid"] = parentId
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
