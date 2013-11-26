package main

import (
	"fmt"
	"github.com/hypebeast/picoui/picoui-lib"
	"time"
)

var (
	ui    *picoui.PicoUi
	title *picoui.TextBox
)

func staticPage() {
	page := ui.NewPage("Static Content", "Back", mainMenu)
	page.AddTextbox("Header 1 Text", "h1")
	page.AddTextbox("Header 2 Text", "h2")
	page.AddTextbox("Header 3 Text", "h3")
	page.AddTextbox("Normal Text", "")
	page.AddElement("hr")
	page.AddTextbox("Lorem ipsum dolor sit amet, consectetuer adipiscing elit, sed diam nonummy nibh euismod tincidunt ut laoreet dolore magna aliquam erat volutpat. Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat.", "")
	page.AddElement("hr")
	page.AddTextbox("Your code can update page contents any time.", "h3")
	update := page.AddTextbox("Like this one...", "")
	for i := 1; i <= 10; i++ {
		time.Sleep(1 * time.Second)
		update.SetText(fmt.Sprintf("Seconds: %d", i))
	}
}

func buttonsPage() {
	page := ui.NewPage("Buttons", "Back", mainMenu)
	title = page.AddTextbox("Buttons!", "h1")
	page.AddElement("hr")
	page.AddTextbox("Buttons with click handlers", "h2")
	page.AddButton("Up Button &uarr;", nil, "", onUp)
	page.AddElement("p")
	page.AddButton("Down Button &darr;", nil, "", onDown)
	page.AddElement("hr")

	page.AddTextbox("Different Colors", "h2")
	page.AddButton("Default", nil, "", nil)
	page.AddElement("p")
	page.AddButton("button-light", []string{"button-light"}, "", nil)
	page.AddElement("p")
	page.AddButton("button-stable", []string{"button-stable"}, "", nil)
	page.AddElement("p")
	page.AddButton("button-positive", []string{"button-positive"}, "", nil)
	page.AddElement("p")
	page.AddButton("button-balanced", []string{"button-balanced"}, "", nil)
	page.AddElement("p")
	page.AddButton("button-energized", []string{"button-energized"}, "", nil)
	page.AddElement("p")
	page.AddButton("button-assertive", []string{"button-assertive"}, "", nil)
	page.AddElement("p")
	page.AddButton("button-royal", []string{"button-royal"}, "", nil)
	page.AddElement("p")
	page.AddButton("button-dark", []string{"button-dark"}, "", nil)
	page.AddElement("hr")

	page.AddTextbox("Block Buttons", "h2")
	page.AddButton("button-light", []string{"button-block", "button-light"}, "", nil)
	page.AddElement("p")
	page.AddButton("button-energized", []string{"button-block", "button-energized"}, "", nil)
	page.AddElement("p")
	page.AddButton("button-assertive", []string{"button-block", "button-assertive"}, "", nil)
	page.AddElement("p")

	page.AddTextbox("Buttons with Icons", "h2")
	page.AddButton("button-light", []string{"button-positive"}, "ion-navicon", nil)
	page.AddElement("p")
	page.AddButton("button-light", []string{"button-royal"}, "ion-email", nil)
}

func togglesPage() {
	page := ui.NewPage("Toggles", "Back", mainMenu)
	title = page.AddTextbox("Home Automation Appliance", "h1")
	list := page.AddList()
	list.AddToggle("Lights", lightsHandler)
	list.AddToggle("TV", tvHandler)
	list.AddToggle("Refrigerator", refrigeratorHandler)
}

func checkboxesPage() {
	page := ui.NewPage("Checkboxes", "Back", mainMenu)
	title = page.AddTextbox("Home Automation Appliance", "h1")
	list := page.AddList()
	list.AddCheckbox("Lights", lightsHandler)
	list.AddCheckbox("TV", tvHandler)
	list.AddCheckbox("Refrigerator", refrigeratorHandler)
}

func inputPage() {
	page := ui.NewPage("Inputs", "Back", mainMenu)
	input1 := page.AddInput("text", "Input 1")
	input2 := page.AddInput("text", "Input 2")
	page.AddElement("hr")
	text := page.AddTextbox("Here goes the text from Input 1 + Input 2", "h3")
	page.AddElement("hr")

	buttonCallback := func() {
		text.SetText(input1.GetText() + input2.GetText())
	}

	page.AddButton("Get Text", nil, "", buttonCallback)
}

// func imagesPage() {
// 	page := ui.NewPage("Images", "Back", mainMenu)
// 	page.AddImage("nature3.png")
// 	page.AddElement("p")
// 	page.AddImage("Beauty-of-nature.jpg")
// }

func onUp() {
	title.SetText("Up!")
}

func onDown() {
	title.SetText("Down!")
}

func lightsHandler(v bool) {
	title.SetText("Toggled Lights: " + fmt.Sprintf("%t", v))
}

func tvHandler(v bool) {
	title.SetText("Toggled TV: " + fmt.Sprintf("%t", v))
}

func refrigeratorHandler(v bool) {
	title.SetText("Toggled Refrigerator: " + fmt.Sprintf("%t", v))
}

func mainMenu() {
	page := ui.NewPage("PicoUi", "", nil)
	list := page.AddList()
	list.AddItem("Static Content", false, staticPage)
	list.AddItem("Buttons", false, buttonsPage)
	list.AddItem("Toggles", false, togglesPage)
	list.AddItem("Checkboxes", false, checkboxesPage)
	list.AddItem("Inputs", false, inputPage)
	// list.AddItem("Images", false, false, imagesPage, nil)
}

func main() {
	ui = picoui.NewPicoUi(1000)
	mainMenu()
	ui.Run()
}
