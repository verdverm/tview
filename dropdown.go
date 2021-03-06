package tview

import (
	"github.com/gdamore/tcell"
)

// dropDownOption is one option that can be selected in a drop-down primitive.
type DropDownOption struct {
	Text     string      // The text to be displayed in the drop-down.
	Value    interface{} // The value associated with this item, passed on submit
	Selected func()      // The (optional) callback for when this option was selected.
}

// DropDown is a one-line box (three lines if there is a title) where the
// user can enter text.
//
// See https://github.com/rivo/tview/wiki/DropDown for an example.
type DropDown struct {
	*Box

	// The options from which the user can choose.
	options []*DropDownOption

	// The index of the currently selected option. Negative if no option is
	// currently selected.
	currentOption int

	// Set to true if the options are visible and selectable.
	open bool

	// The list element for the options.
	list *List

	// The text to be displayed before the input area.
	label string

	// The label color.
	labelColor tcell.Color

	// The background color of the input area.
	fieldBackgroundColor tcell.Color

	// The text color of the input area.
	fieldTextColor tcell.Color

	// The screen width of the input area. A value of 0 means extend as much as
	// possible.
	fieldWidth int

	// An optional function which is called when the user indicated that they
	// are done selecting options. The key which was pressed is provided (tab,
	// shift-tab, or escape).
	done func(tcell.Key)
}

// NewDropDown returns a new drop-down.
func NewDropDown() *DropDown {
	list := NewList().ShowSecondaryText(false)
	list.SetMainTextColor(Styles.PrimitiveBackgroundColor).
		SetSelectedTextColor(Styles.PrimitiveBackgroundColor).
		SetSelectedBackgroundColor(Styles.PrimaryTextColor).
		SetBackgroundColor(Styles.MoreContrastBackgroundColor)

	d := &DropDown{
		Box:                  NewBox(),
		currentOption:        -1,
		list:                 list,
		labelColor:           Styles.SecondaryTextColor,
		fieldBackgroundColor: Styles.ContrastBackgroundColor,
		fieldTextColor:       Styles.PrimaryTextColor,
	}

	d.focus = d

	return d
}

func (d *DropDown) GetValues() map[string]interface{} {
	_, c := d.GetCurrentOption()
	return map[string]interface{}{
		d.name: c.Value,
	}
}

func (d *DropDown) SetValues(values map[string]interface{}) {
	list := NewList().ShowSecondaryText(false)
	list.SetMainTextColor(Styles.PrimitiveBackgroundColor).
		SetSelectedTextColor(Styles.PrimitiveBackgroundColor).
		SetSelectedBackgroundColor(Styles.PrimaryTextColor).
		SetBackgroundColor(Styles.MoreContrastBackgroundColor)

	d.list = list

	d.options = []*DropDownOption{}
	for k, v := range values {
		d.AddOption(k, v, nil)
	}

}

// SetCurrentOption sets the index of the currently selected option.
func (d *DropDown) SetCurrentOption(index int) *DropDown {
	d.currentOption = index
	d.list.SetCurrentItem(index)
	return d
}

// GetCurrentOption returns the index of the currently selected option as well
// as its text. If no option was selected, -1 and an empty string is returned.
func (d *DropDown) GetCurrentOption() (int, *DropDownOption) {
	var opt *DropDownOption
	if d.currentOption >= 0 && d.currentOption < len(d.options) {
		opt = d.options[d.currentOption]
	}
	return d.currentOption, opt
}

// SetLabel sets the text to be displayed before the input area.
func (d *DropDown) SetLabel(label string) *DropDown {
	d.label = label
	return d
}

// GetLabel returns the text to be displayed before the input area.
func (d *DropDown) GetLabel() string {
	return d.label
}

// SetLabelColor sets the color of the label.
func (d *DropDown) SetLabelColor(color tcell.Color) *DropDown {
	d.labelColor = color
	return d
}

// SetFieldBackgroundColor sets the background color of the options area.
func (d *DropDown) SetFieldBackgroundColor(color tcell.Color) *DropDown {
	d.fieldBackgroundColor = color
	return d
}

// SetFieldTextColor sets the text color of the options area.
func (d *DropDown) SetFieldTextColor(color tcell.Color) *DropDown {
	d.fieldTextColor = color
	return d
}

// SetFormAttributes sets attributes shared by all form items.
func (d *DropDown) SetFormAttributes(label string, labelColor, bgColor, fieldTextColor, fieldBgColor tcell.Color) FormItem {
	d.label = label
	d.labelColor = labelColor
	d.backgroundColor = bgColor
	d.fieldTextColor = fieldTextColor
	d.fieldBackgroundColor = fieldBgColor
	return d
}

// SetFieldWidth sets the screen width of the options area. A value of 0 means
// extend to as long as the longest option text.
func (d *DropDown) SetFieldWidth(width int) *DropDown {
	d.fieldWidth = width
	return d
}

// GetFieldWidth returns this primitive's field screen width.
func (d *DropDown) GetFieldWidth() int {
	if d.fieldWidth > 0 {
		return d.fieldWidth
	}
	fieldWidth := 0
	for _, option := range d.options {
		width := StringWidth(option.Text)
		if width > fieldWidth {
			fieldWidth = width
		}
	}
	return fieldWidth
}

// AddOption adds a new selectable option to this drop-down. The "selected"
// callback is called when this option was selected. It may be nil.
func (d *DropDown) AddOption(text string, value interface{}, selected func()) *DropDown {
	d.options = append(d.options, &DropDownOption{Text: text, Value: value, Selected: selected})
	d.list.AddItem(text, "", 0, selected)
	return d
}

// SetOptions replaces all current options with the ones provided and installs
// one callback function which is called when one of the options is selected.
// It will be called with the option's text and its index into the options
// slice. The "selected" parameter may be nil.
func (d *DropDown) SetOptions(texts []string, values []interface{}, selected func(text string, value interface{}, index int)) *DropDown {
	d.list.Clear()
	d.options = nil
	for index, text := range texts {
		func(t string, v interface{}, i int) {
			d.AddOption(t, v, func() {
				if selected != nil {
					selected(t, v, i)
				}
			})
		}(text, values[index], index)
	}
	return d
}

// SetDoneFunc sets a handler which is called when the user is done selecting
// options. The callback function is provided with the key that was pressed,
// which is one of the following:
//
//   - KeyEscape: Abort selection.
//   - KeyTab: Move to the next field.
//   - KeyBacktab: Move to the previous field.
func (d *DropDown) SetDoneFunc(handler func(key tcell.Key)) *DropDown {
	d.done = handler
	return d
}

// SetFinishedFunc calls SetDoneFunc().
func (d *DropDown) SetFinishedFunc(handler func(key tcell.Key)) FormItem {
	return d.SetDoneFunc(handler)
}

// SetFinishedFunc calls SetDoneFunc().
func (d *DropDown) SetFinishedFunction(handler func(key tcell.Key)) {
	d.SetDoneFunc(handler)
}

// Draw draws this primitive onto the screen.
func (d *DropDown) Draw(screen tcell.Screen) {
	d.Box.Draw(screen)

	d.RLock()
	defer d.RUnlock()

	// Prepare.
	x, y, width, height := d.GetInnerRect()
	rightLimit := x + width
	if height < 1 || rightLimit <= x {
		return
	}

	// Draw label.
	_, drawnWidth := Print(screen, d.label, x, y, rightLimit-x, AlignLeft, d.labelColor)
	x += drawnWidth

	// What's the longest option text?
	maxWidth := 0
	for _, option := range d.options {
		strWidth := StringWidth(option.Text)
		if strWidth > maxWidth {
			maxWidth = strWidth
		}
	}

	// Draw selection area.
	fieldWidth := d.fieldWidth
	if fieldWidth == 0 {
		fieldWidth = maxWidth
	}
	if rightLimit-x < fieldWidth {
		fieldWidth = rightLimit - x
	}
	fieldStyle := tcell.StyleDefault.Background(d.fieldBackgroundColor)
	if d.GetFocusable().HasFocus() && !d.open {
		fieldStyle = fieldStyle.Background(d.fieldTextColor)
	}
	for index := 0; index < fieldWidth; index++ {
		screen.SetContent(x+index, y, ' ', nil, fieldStyle)
	}

	// Draw selected text.
	if d.currentOption >= 0 && d.currentOption < len(d.options) {
		color := d.fieldTextColor
		if d.GetFocusable().HasFocus() && !d.open {
			color = d.fieldBackgroundColor
		}
		Print(screen, d.options[d.currentOption].Text, x, y, fieldWidth, AlignLeft, color)
	}

	// Draw options list.
	if d.HasFocus() && d.open {
		// We prefer to drop down but if there is no space, maybe drop up?
		lx := x
		ly := y + 1
		lwidth := maxWidth
		lheight := len(d.options)
		_, sheight := screen.Size()
		if ly+lheight >= sheight && ly-lheight-1 >= 0 {
			ly = y - lheight
		}
		d.list.SetRect(lx, ly, lwidth, lheight)
		d.list.Draw(screen)
	}
}

// InputHandler returns the handler for this primitive.
func (d *DropDown) InputHandler() func(tcell.Event, func(Primitive)) {
	return d.wrapInputHandler(func(event tcell.Event, setFocus func(p Primitive)) {
		switch evt := event.(type) {

		// Process key event.
		case *tcell.EventKey:
			switch key := evt.Key(); key {
			case tcell.KeyEnter, tcell.KeyRune, tcell.KeyDown:
				if key == tcell.KeyRune && evt.Rune() != ' ' {
					break
				}
				d.open = true
				d.list.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
					// An option was selected. Close the list again.
					d.open = false
					setFocus(d)
					d.currentOption = index

					// Trigger "selected" event.
					if d.options[d.currentOption].Selected != nil {
						d.options[d.currentOption].Selected()
					}
				})
				d.list.SetDoneFunc(func() {
					d.open = false
					setFocus(d)
				})
				setFocus(d.list)
			case tcell.KeyEscape, tcell.KeyTab, tcell.KeyBacktab:
				if d.done != nil {
					d.done(key)
				}
			}
		}
	})
}

// Focus is called by the application when the primitive receives focus.
func (d *DropDown) Focus(delegate func(p Primitive)) {
	d.Box.Focus(delegate)
	if d.open {
		delegate(d.list)
	}
}

// HasFocus returns whether or not this primitive has focus.
func (d *DropDown) HasFocus() bool {
	if d.open {
		return d.list.HasFocus()
	}
	return d.hasFocus
}
