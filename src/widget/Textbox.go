package widget

type TextBox struct {
	text string
	cursor int
}

func (box *TextBox) SetText(text string) {
	box.text = text
	Changed <- true
}

func (box *TextBox) AddText(text string) {
	box.text += text
	Changed <- true
}

func (box *TextBox) GetText() string {
	return box.text
}

func (box *TextBox) GetFormattedText() string {
	return box.text+"|"
}
