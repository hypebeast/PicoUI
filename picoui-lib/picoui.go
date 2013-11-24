package picoui

import (
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

// Represents a PicoUi application.
type PicoUi struct {
	handlers *HandleController
}

// Function type for handling a click event.
type clickHandler func()

// Function type for handling a toggle event.
type toggleHandler func(value bool)

// NewMiUi creates a new PicoUi application and returns it.
func NewPicoUi(timeout int) *PicoUi {
	return &PicoUi{handlers: newHandleController(timeout)}
}

// Run starts the server.
func (ui *PicoUi) Run() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.ListenAndServe(fmt.Sprintf("%s:%d", SERVER_HOST, SERVER_PORT), nil)
}

// NewPage creates a new Page and returns it.
func (ui *PicoUi) NewPage(title string, prevText string, onPrevClick clickHandler) *Page {
	page := newPage(ui, title, prevText, onPrevClick)
	ui.handlers.newPage("ui", title, page)
	page.pagePost()
	return page
}

/////
// Helper functions
/////

func createId() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.FormatInt(r.Int63(), 10)
}
