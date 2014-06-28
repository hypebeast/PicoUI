package picoui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Represents a handle controller for PicoUi. It is responsible for the
// communication with a client application.
type HandleController struct {
	messages                []Message
	messagesSinceLastReload []Message
	currentPage             string
	currentPageTitle        string
	currentPageObj          *Page
	timeout                 int
	inBuffer                []string
}

// Represents a JSON message that is used to communicate with the client application.
type Message struct {
	Cmd        string                 `json:"cmd"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// NewHandleController creates a new HandleController and returns it.
func newHandleController(timeout int) *HandleController {
	HandleController := &HandleController{timeout: timeout, currentPage: "/", currentPageTitle: ""}

	// Add all handler functions
	http.HandleFunc("/ping", HandleController.pingHandler)
	http.HandleFunc("/init", HandleController.initHandler)
	http.HandleFunc("/poll", HandleController.pollHandler)
	http.HandleFunc("/click", HandleController.clickHandler)
	http.HandleFunc("/toggle", HandleController.toggleHandler)
	http.HandleFunc("/slide", HandleController.slideHandler)
	http.HandleFunc("/state", HandleController.stateHandler)

	return HandleController
}

// newPage sets the given page as the current active page for this handler.
func (hd *HandleController) newPage(name string, title string, page *Page) {
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
func (hd *HandleController) pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

// Implements the init handler. This handler should be called by the client if a
// page reload happens.
func (hd *HandleController) initHandler(w http.ResponseWriter, r *http.Request) {
	hd.pageReload()
	fmt.Fprintf(w, "ok")
}

// Implements the poll handler. It returns messages until a timeout occurs or
// all messages are returned.
func (hd *HandleController) pollHandler(w http.ResponseWriter, r *http.Request) {
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

func (hd *HandleController) stateHandler(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	text := vals["msg"][0]
	hd.inBuffer = append(hd.inBuffer, text)
}

func (hd *HandleController) clickHandler(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	eid := vals["eid"][0]
	hd.currentPageObj.handleClick(eid)
}

func (hd *HandleController) toggleHandler(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	eid := vals["eid"][0]
	v, err := strconv.ParseBool(vals["v"][0])
	if err != nil {
		return
	}
	hd.currentPageObj.handleToogle(eid, v)
}

func (hd *HandleController) slideHandler(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	eid := vals["eid"][0]
	v, err := strconv.Atoi(vals["v"][0])
	if err != nil {
		return
	}
	hd.currentPageObj.handleSlide(eid, v)
}

// pageReload handles a page reload.
func (hd *HandleController) pageReload() {
	if len(hd.messagesSinceLastReload) > len(hd.messages) {
		hd.messages = nil
		for _, msg := range hd.messagesSinceLastReload {
			hd.messages = append(hd.messages, msg)
		}
	}
}

// enqueue adds the given message to the message queue.
func (hd *HandleController) enqueue(msg Message) {
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

func (hd *HandleController) enqueue_and_wait_for_result(msg Message) string {
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
func (hd *HandleController) flushQueue() {
	hd.messages = nil
	hd.messagesSinceLastReload = nil
}
