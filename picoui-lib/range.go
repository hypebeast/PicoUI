package picoui

import (
	"strconv"
)

// Represents a Range (http://ionicframework.com/docs/components/#range) element.
type Range struct {
	ui      *PicoUi
	page    *Page
	id      string
	slideId string
	onSlide slideHandler
}

// newRange creates a new Range and returns it.
func newRange(ui *PicoUi, page *Page, min int, max int, iconLeft string, iconRight string, onSlide slideHandler) *Range {
	id := "range_" + createId()
	slide_id := "slide_" + createId()
	item := Range{ui: ui, page: page, id: id, slideId: slide_id, onSlide: onSlide}

	msg := Message{Cmd: "addrange"}
	attributes := make(map[string]interface{})
	attributes["eid"] = id
	attributes["slideid"] = slide_id
	attributes["min"] = strconv.Itoa(min)
	attributes["max"] = strconv.Itoa(max)
	attributes["iconleft"] = iconLeft
	attributes["iconright"] = iconRight
	msg.Attributes = attributes
	ui.handlers.enqueue(msg)

	return &item
}

func (r *Range) GetValue() int {
	// TODO
	return 0
}

func (r *Range) SetValue(val int) {
	// TODO
}
