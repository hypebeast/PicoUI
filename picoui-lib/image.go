package picoui

// Represents an Image element.
type Image struct {
	ui *PicoUi
	id string
}

/////
// MiUiImage
/////

func NewImage(ui *PicoUi, source string) *Image {
	return &Image{}
}
