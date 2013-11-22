package main

import (
	"fmt"
	"github.com/hypebeast/miui/miui-lib"
	"time"
)

var (
	ui    *miui.MiUi
	title *miui.MiUiTextBox
)

func staticPage() {
	page := ui.NewUiPage("Static Content", "Back", mainMenu)
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
	page := ui.NewUiPage("Buttons", "Back", mainMenu)
	title = page.AddTextbox("Buttons!", "h1")
	page.AddElement("hr")
	page.AddButton("Up Button &uarr;", onUp)
	page.AddButton("Down Button &darr;", onDown)
}

func togglesPage() {
	page := ui.NewUiPage("Toggles", "Back", mainMenu)
	title = page.AddTextbox("Home Automation Appliance", "h1")
	list := page.AddList()
	list.AddItem("Lights", false, true, nil, lightsHandler)
	list.AddItem("TV", false, true, nil, tvHandler)
	list.AddItem("Refrigerator", false, true, nil, refrigeratorHandler)
}

func inputPage() {
	page := ui.NewUiPage("Inputs", "Back", mainMenu)
	input1 := page.AddInput("text", "Input 1")
	input2 := page.AddInput("text", "Input 2")
	page.AddElement("hr")
	text := page.AddTextbox("Here goes the text from Input 1 + Input 2", "h3")
	page.AddElement("hr")

	buttonCallback := func() {
		text.SetText(input1.GetText() + input2.GetText())
	}

	page.AddButton("Get Text", buttonCallback)
}

// func imagesPage() {
// 	page := ui.NewUiPage("Images", "Back", mainMenu)
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
	page := ui.NewUiPage("PicoUi", "", nil)
	list := page.AddList()
	list.AddItem("Static Content", true, false, staticPage, nil)
	list.AddItem("Buttons", true, false, buttonsPage, nil)
	list.AddItem("Toggles", true, false, togglesPage, nil)
	list.AddItem("Inputs", true, false, inputPage, nil)
	//list.AddItem("Images", true, false, imagesPage, nil)
}

func main() {
	ui = miui.NewMiUi(1000)
	mainMenu()
	ui.Run()
}
