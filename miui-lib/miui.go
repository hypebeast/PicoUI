package miui

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	SERVER_PORT  = 9999
	SERVER_HOST  = "0.0.0.0"
	MAX_MESSAGES = 250
)

// Represents a MiUi application.
type MiUi struct {
	handlers *Handlers
}

// Represents the handlers for MiUi. A Handlers is responsible for the
// communication with a client application.
type Handlers struct {
	messages                []Message
	messagesSinceLastReload []Message
	currentPage             string
	currentPageTitle        string
	currentPageObj          *MiUiPage
	timeout                 int
	inBuffer                []string
}

// Represents a JSON message that is used to communicate with the client application.
type Message struct {
	Cmd        string                 `json:"cmd"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// Represents an user interface page.
type MiUiPage struct {
	ui          *MiUi
	id          string
	title       string
	prevText    string
	onPrevClick clickHandler
	prevId      string
	elements    []interface{}
	clickables  map[string]clickHandler
	toggables   map[string]toggleHandler
	inputs      map[string]interface{}
}

// Represents an UI button.
type MiUiButton struct {
	ui             *MiUi
	id             string
	onClickHandler clickHandler
}

// Represents an UI textbox.
type MiUiTextBox struct {
	ui *MiUi
	id string
}

// Represents an UI List.
type MiUiList struct {
	ui   *MiUi
	page *MiUiPage
	id   string
}

// Represents an UI list item.
type MiUiListItem struct {
	ui              *MiUi
	id              string
	parentId        string
	onClickHandler  clickHandler
	onToggleHandler toggleHandler
	toggleId        string
}

// Represents an input UI element.
type MiUiInput struct {
	ui *MiUi
	id string
}

// Represents an image UI element.
type MiUiImage struct {
	ui *MiUi
	id string
}

// Function type for handling a click event.
type clickHandler func()

// Function type for handling a toggle event.
type toggleHandler func(value bool)

/////
// MiUi
/////

// NewMiUi creates a new MiUi and returns it.
func NewMiUi(timeout int) *MiUi {
	return &MiUi{handlers: NewHandlers(timeout)}
}

// Run starts the server.
func (ui *MiUi) Run() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.ListenAndServe(fmt.Sprintf("%s:%d", SERVER_HOST, SERVER_PORT), nil)
}

// NewPage creates a new UI page and returns it.
func (ui *MiUi) NewUiPage(title string, prevText string, onPrevClick clickHandler) *MiUiPage {
	page := newMiUiPage(ui, title, prevText, onPrevClick)
	ui.handlers.NewPage("ui", title, page)
	page.pagePost()
	return page
}

/////
// Handlers
/////

// NewHandlers creates a new Handlers and returns it.
func NewHandlers(timeout int) *Handlers {
	handlers := &Handlers{timeout: timeout, currentPage: "/", currentPageTitle: ""}

	// Add all handlers
	http.HandleFunc("/ping", handlers.pingHandler)
	http.HandleFunc("/init", handlers.initHandler)
	http.HandleFunc("/poll", handlers.pollHandler)
	http.HandleFunc("/click", handlers.clickHandler)
	http.HandleFunc("/toggle", handlers.toggleHandler)
	http.HandleFunc("/state", handlers.stateHandler)

	return handlers
}

// NewPages sets the given page as the current active page for this handler.
func (hd *Handlers) NewPage(name string, title string, page *MiUiPage) {
	hd.currentPage = "/" + name
	hd.currentPageTitle = title
	hd.currentPageObj = page
	hd.flushQueue()
	attributes := make(map[string]interface{})
	attributes["page"] = name
	attributes["title"] = title
	attributes["eid"] = page.id
	hd.enqueue(Message{Cmd: "newpage", Attributes: attributes})
}

// Implements the ping handler.
func (hd *Handlers) pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

// Implements the init handler. This handler should be called by the client if a
// page reload happens.
func (hd *Handlers) initHandler(w http.ResponseWriter, r *http.Request) {
	hd.pageReload()
	fmt.Fprintf(w, "ok")
}

// Implements the poll handler. It returns messages until a timeout occurs or
// all messages are returned.
func (hd *Handlers) pollHandler(w http.ResponseWriter, r *http.Request) {
	var waiting int = 0
	enc := json.NewEncoder(w)

	for waiting < hd.timeout {
		// If a message is available encode it and return it
		if len(hd.messages) > 0 {
			msg := hd.messages[0]
			hd.messages = hd.messages[1:]

			err := enc.Encode(&msg)
			if err != nil {
				fmt.Println("encoding error", err)
			}

			return
		}
		// Sleep for some time
		time.Sleep(10 * time.Millisecond)
		waiting += 10
	}

	timeoutMsg := Message{Cmd: "timeout", Attributes: nil}
	err := enc.Encode(&timeoutMsg)
	if err != nil {
		fmt.Println("encoding error", err)
	}
}

func (hd *Handlers) stateHandler(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	text := vals["msg"][0]
	hd.inBuffer = append(hd.inBuffer, text)
}

func (hd *Handlers) clickHandler(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	eid := vals["eid"][0]
	hd.currentPageObj.handleClick(eid)
}

func (hd *Handlers) toggleHandler(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	eid := vals["eid"][0]
	v, err := strconv.ParseBool(vals["v"][0])
	if err != nil {
		return
	}
	hd.currentPageObj.handleToogle(eid, v)
}

// pageReload handles a page reload.
func (hd *Handlers) pageReload() {
	if len(hd.messagesSinceLastReload) > len(hd.messages) {
		hd.messages = nil
		for _, msg := range hd.messagesSinceLastReload {
			hd.messages = append(hd.messages, msg)
		}
	}
}

// enqueue adds the given message to the message queue.
func (hd *Handlers) enqueue(msg Message) {
	hd.messages = append(hd.messages, msg)

	// Save the message for the next page reload
	if id, ok := msg.Attributes["eid"]; ok {
		// Check if a msg for the given UI element is already saved
		found := false
		for _, m := range hd.messagesSinceLastReload {
			if m.Attributes["eid"] == id {
				found = true
				break
			}
		}

		// A msg with the given id was not found; add it
		if !found {
			hd.messagesSinceLastReload = append(hd.messagesSinceLastReload, msg)
		}
	}

	if len(hd.messages) > MAX_MESSAGES {
		hd.messages = hd.messages[:len(hd.messages)-1]
	}
}

func (hd *Handlers) enqueue_and_wait_for_result(msg Message) string {
	// enqueue the message
	hd.messages = append(hd.messages, msg)

	if len(hd.messages) > MAX_MESSAGES {
		hd.messages = hd.messages[:len(hd.messages)-1]
	}

	// Wait until all messages are received by the app and the message queue is empty
	done := false
	for !done {
		// Sleep for some time
		time.Sleep(10 * time.Millisecond)

		if len(hd.messages) == 0 {
			done = true
		}
	}

	// Wait for the response
	received := false
	var result string
	for !received {
		// Sleep for some time
		time.Sleep(10 * time.Millisecond)

		if len(hd.inBuffer) > 0 {
			received = true
		}

		if received {
			// Get the last message
			result = hd.inBuffer[len(hd.inBuffer)-1]
			hd.inBuffer = hd.inBuffer[:len(hd.inBuffer)-1]
		}
	}

	return result
}

// flushQueue removes all items from the message queues.
func (hd *Handlers) flushQueue() {
	hd.messages = nil
	hd.messagesSinceLastReload = nil
}

/////
// MiUiPage
/////

// newMiUiPage creates a new MiUiPage and returns it.
func newMiUiPage(ui *MiUi, title string, prevText string, onPrevClick clickHandler) *MiUiPage {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := "page_" + strconv.FormatInt(r.Int63(), 10)
	page := &MiUiPage{
		ui:          ui,
		id:          id,
		title:       title,
		prevText:    prevText,
		onPrevClick: onPrevClick,
		elements:    make([]interface{}, 10),
		clickables:  make(map[string]clickHandler),
		toggables:   make(map[string]toggleHandler)}
	return page
}

// pagePost enqueues a new 'pagepost' message. This message must be send when a
// new page was created.
func (p *MiUiPage) pagePost() {
	msg := Message{Cmd: "pagepost"}
	attributes := make(map[string]interface{})
	attributes["title"] = p.title
	if p.prevText != "" && p.onPrevClick != nil {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		p.prevId = "button_" + strconv.FormatInt(r.Int63(), 10)
		p.clickables[p.prevId] = p.onPrevClick
		attributes["prevtxt"] = p.prevText
		attributes["previd"] = p.prevId
		attributes["eid"] = p.prevId
	}
	msg.Attributes = attributes
	p.ui.handlers.enqueue(msg)
}

func (p *MiUiPage) printLine(line string) {
	// TODO
}

// AddTextbox creates a new textbox. The argument text sets the text for the textbox
// and the argument element specifies the HTML element type. If an empty string is
// given for element, then the 'p' element is used.
func (p *MiUiPage) AddTextbox(text string, element string) *MiUiTextBox {
	var box *MiUiTextBox
	if element == "" {
		box = NewMiUiTextBox(text, "p", p.ui)
	} else {
		box = NewMiUiTextBox(text, element, p.ui)
	}
	p.elements = append(p.elements, box)
	return box
}

func (p *MiUiPage) AddButton(text string, onClick clickHandler) *MiUiButton {
	button := NewMiUiButton(text, p.ui, onClick)
	p.elements = append(p.elements, button)
	p.clickables[button.id] = onClick
	return button
}

func (p *MiUiPage) AddElement(element string) *MiUiTextBox {
	ele := NewMiUiTextBox("", element, p.ui)
	p.elements = append(p.elements, ele)
	return ele
}

func (p *MiUiPage) AddInput(inputType string, placeholder string) *MiUiInput {
	input := NewMiUiInput(p.ui, inputType, placeholder)
	p.elements = append(p.elements, input)
	return input
}

func (p *MiUiPage) AddImage() {
	// TODO
}

func (p *MiUiPage) AddList() *MiUiList {
	list := NewMiUiList(p.ui, p)
	p.elements = append(p.elements, list)
	return list
}

// HandleClick handles a click event.
func (p *MiUiPage) handleClick(id string) {
	if handler, ok := p.clickables[id]; ok {
		handler()
	}
}

// HandleToogle handles a toggle event.
func (p *MiUiPage) handleToogle(id string, v bool) {
	if handler, ok := p.toggables[id]; ok {
		handler(v)
	}
}

/////
// MiUiButton
/////

func NewMiUiButton(text string, ui *MiUi, onClick clickHandler) *MiUiButton {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := "button_" + strconv.FormatInt(r.Int63(), 10)
	button := MiUiButton{ui: ui, id: id, onClickHandler: onClick}

	msg := Message{Cmd: "addbutton"}
	attributes := make(map[string]interface{})
	attributes["eid"] = id
	attributes["txt"] = text
	msg.Attributes = attributes
	ui.handlers.enqueue(msg)

	return &button
}

// SetText sets the given text.
func (b *MiUiButton) SetText(text string) {
	msg := Message{Cmd: "updateinner"}
	attributes := make(map[string]interface{})
	attributes["eid"] = b.id
	attributes["txt"] = text
	msg.Attributes = attributes
	b.ui.handlers.enqueue(msg)
}

/////
// MiUiList
/////

func NewMiUiList(ui *MiUi, page *MiUiPage) *MiUiList {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := "list_" + strconv.FormatInt(r.Int63(), 10)
	list := MiUiList{ui: ui, page: page, id: id}

	msg := Message{Cmd: "addlist"}
	attributes := make(map[string]interface{})
	attributes["eid"] = id
	msg.Attributes = attributes
	ui.handlers.enqueue(msg)

	return &list
}

// AddItem adds an MiUiListItem to the list.
func (l *MiUiList) AddItem(text string, chevron bool, toggle bool, onClick clickHandler, onToggle toggleHandler) *MiUiListItem {
	item := NewMiUiListItem(l.ui, l.id, text, chevron, toggle, onClick, onToggle)
	l.page.elements = append(l.page.elements, item)
	if onClick != nil {
		l.page.clickables[item.id] = onClick
	}
	if onToggle != nil {
		l.page.toggables[item.toggleId] = onToggle
	}

	return item
}

/////
// MiUiListItem
/////

// NewMiUiListItem creates a new MiUiListItem.
func NewMiUiListItem(ui *MiUi, parentId string, text string, chevron bool, toggle bool, onClick clickHandler, onToggle toggleHandler) *MiUiListItem {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := "listitem_" + strconv.FormatInt(r.Int63(), 10)
	tg_id := "tg_" + strconv.FormatInt(r.Int63(), 10)
	item := MiUiListItem{
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

/////
// MiUiTextBox
/////

func NewMiUiTextBox(text string, element string, ui *MiUi) *MiUiTextBox {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := "textbox_" + strconv.FormatInt(r.Int63(), 10)
	box := MiUiTextBox{ui: ui, id: id}

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
func (b *MiUiTextBox) SetText(text string) {
	msg := Message{Cmd: "updateinner"}
	attributes := make(map[string]interface{})
	attributes["eid"] = b.id
	attributes["txt"] = text
	msg.Attributes = attributes
	b.ui.handlers.enqueue(msg)
}

/////
// MiUiInput
/////

// NewMiUiInput creates and returns a new MiUiInput.
func NewMiUiInput(ui *MiUi, inputType string, placeholder string) *MiUiInput {
	id := "input_" + createId()
	input := MiUiInput{ui: ui, id: id}

	msg := Message{Cmd: "addinput"}
	attributes := make(map[string]interface{})
	attributes["eid"] = id
	attributes["type"] = inputType
	attributes["placeholder"] = placeholder
	msg.Attributes = attributes
	ui.handlers.enqueue(msg)

	return &input
}

func (i *MiUiInput) GetText() string {
	msg := Message{Cmd: "getinput"}
	attributes := make(map[string]interface{})
	attributes["eid"] = i.id
	msg.Attributes = attributes
	return i.ui.handlers.enqueue_and_wait_for_result(msg)
}

/////
// MiUiImage
/////

/////
// Helper functions
/////

func createId() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.FormatInt(r.Int63(), 10)
}
