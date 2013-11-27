package picoui

// Represents a page element.
type Page struct {
	ui            *PicoUi
	id            string
	title         string
	prevText      string
	onPrevClick   clickHandler
	prevId        string
	elements      []interface{}
	clickables    map[string]clickHandler
	toggables     map[string]toggleHandler
	slideHandlers map[string]slideHandler
	inputs        map[string]interface{}
}

// newPage creates a new Page and returns it.
func newPage(ui *PicoUi, title string, prevText string, onPrevClick clickHandler) *Page {
	id := "page_" + createId()
	page := &Page{
		ui:            ui,
		id:            id,
		title:         title,
		prevText:      prevText,
		onPrevClick:   onPrevClick,
		elements:      make([]interface{}, 10),
		clickables:    make(map[string]clickHandler),
		toggables:     make(map[string]toggleHandler),
		slideHandlers: make(map[string]slideHandler)}
	return page
}

// pagePost enqueues a new 'pagepost' message. This message must be send when a
// new page was created.
func (p *Page) pagePost() {
	msg := Message{Cmd: "pagepost"}
	attributes := make(map[string]interface{})
	attributes["title"] = p.title
	if p.prevText != "" && p.onPrevClick != nil {
		p.prevId = "button_" + createId()
		p.clickables[p.prevId] = p.onPrevClick
		attributes["prevtxt"] = p.prevText
		attributes["previd"] = p.prevId
		attributes["eid"] = p.prevId
	}
	msg.Attributes = attributes
	p.ui.handlers.enqueue(msg)
}

func (p *Page) printLine(line string) {
	// TODO
}

// AddTextbox creates a new textbox. The argument text sets the text for the textbox
// and the argument element specifies the HTML element type. If an empty string is
// given for element, then the 'p' element is used.
func (p *Page) AddTextbox(text string, element string) *TextBox {
	var box *TextBox
	if element == "" {
		box = NewTextBox(text, "p", p.ui)
	} else {
		box = NewTextBox(text, element, p.ui)
	}
	p.elements = append(p.elements, box)
	return box
}

// AddButton creates and returns a new Button element.
func (p *Page) AddButton(text string, classAttributes []string, icon string, onClick clickHandler) *Button {
	button := newButton(p.ui, text, classAttributes, icon, onClick)
	p.elements = append(p.elements, button)
	p.clickables[button.id] = onClick
	return button
}

func (p *Page) AddElement(element string) *TextBox {
	ele := NewTextBox("", element, p.ui)
	p.elements = append(p.elements, ele)
	return ele
}

func (p *Page) AddImage() {
	// TODO
}

func (p *Page) AddList() *List {
	list := NewList(p.ui, p)
	p.elements = append(p.elements, list)
	return list
}

// AddRange creates a new Range and returns it.
func (p *Page) AddRange(min int, max int, leftIcon string, rightIcon string, onSlide slideHandler) *Range {
	item := newRange(p.ui, p, min, max, leftIcon, rightIcon, onSlide)
	p.elements = append(p.elements, item)

	if onSlide != nil {
		p.slideHandlers[item.slideId] = onSlide
	}

	return item
}

// handleClick handles a click event.
func (p *Page) handleClick(id string) {
	if handler, ok := p.clickables[id]; ok {
		if handler != nil {
			handler()
		}
	}
}

// handleToogle handles a toggle event.
func (p *Page) handleToogle(id string, v bool) {
	if handler, ok := p.toggables[id]; ok {
		if handler != nil {
			handler(v)
		}
	}
}

// handleSlide handles a slide event.
func (p *Page) handleSlide(id string, v int) {
	if handler, ok := p.slideHandlers[id]; ok {
		if handler != nil {
			handler(v)
		}
	}
}
