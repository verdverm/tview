package tview

import (
	"sync"

	"github.com/gdamore/tcell"
)

// page represents one page of a Pages object.
type Page struct {
	sync.RWMutex

	Name    string    // The page's name.
	Item    Primitive // The page's primitive.
	Resize  bool      // Whether or not to resize the page when it is drawn.
	Visible bool      // Whether or not this page is visible.
}

// Pages is a container for other primitives often used as the application's
// root primitive. It allows to easily switch the visibility of the contained
// primitives.
//
// See https://github.com/rivo/tview/wiki/Pages for an example.
type Pages struct {
	*Box

	// The contained pages.
	curr  *Page
	pages []*Page

	// We keep a reference to the function which allows us to set the focus to
	// a newly visible page.
	setFocus func(p Primitive)

	// An optional handler which is called whenever the visibility or the order of
	// pages changes.
	changed func()
}

// NewPages returns a new Pages object.
func NewPages() *Pages {
	p := &Pages{
		Box: NewBox(),
	}
	p.focus = p
	return p
}

// SetChangedFunc sets a handler which is called whenever the visibility or the
// order of any visible pages changes. This can be used to redraw the pages.
func (p *Pages) SetChangedFunc(handler func()) *Pages {
	//p.Lock()
	//defer p.Unlock()

	p.changed = handler
	return p
}

// AddPage adds a new page with the given name and primitive. If there was
// previously a page with the same name, it is overwritten. Leaving the name
// empty may cause conflicts in other functions.
//
// Visible pages will be drawn in the order they were added (unless that order
// was changed in one of the other functions). If "resize" is set to true, the
// primitive will be set to the size available to the Pages primitive whenever
// the pages are drawn.
func (p *Pages) AddPage(name string, item Primitive, resize, visible bool) *Pages {
	//p.Lock()
	//defer p.Unlock()

	for index, pg := range p.pages {
		if pg.Name == name {
			p.pages = append(p.pages[:index], p.pages[index+1:]...)
			break
		}
	}
	p.pages = append(p.pages, &Page{Item: item, Name: name, Resize: resize, Visible: visible})
	if p.changed != nil {
		p.changed()
	}
	if p.HasFocus() {
		p.Focus(p.setFocus)
	}
	return p
}

// AddAndSwitchToPage calls AddPage(), then SwitchToPage() on that newly added
// page.
func (p *Pages) AddAndSwitchToPage(name string, item Primitive, resize bool) *Pages {
	p.AddPage(name, item, resize, true)
	p.SwitchToPage(name, map[string]interface{}{"activation": "function"})
	return p
}

// RemovePage removes the page with the given name.
func (p *Pages) RemovePage(name string) *Pages {
	//p.Lock()
	//defer p.Unlock()

	hasFocus := p.HasFocus()
	for index, page := range p.pages {
		if page.Name == name {
			p.pages = append(p.pages[:index], p.pages[index+1:]...)
			if page.Visible && p.changed != nil {
				p.changed()
			}
			break
		}
	}
	if hasFocus {
		p.Focus(p.setFocus)
	}
	return p
}

// HasPage returns true if a page with the given name exists in this object.
func (p *Pages) HasPage(name string) bool {
	//p.RLock()
	//defer p.RUnlock()

	for _, page := range p.pages {
		if page.Name == name {
			return true
		}
	}
	return false
}

// HasPage returns true if a page with the given name exists in this object.
func (p *Pages) GetCurrentPage() *Page {
	//p.RLock()
	//defer p.RUnlock()

	return p.curr
}

// HasPage returns true if a page with the given name exists in this object.
func (p *Pages) GetPage(name string) *Page {
	//p.RLock()
	//defer p.RUnlock()

	for _, page := range p.pages {
		if page.Name == name {
			return page
		}
	}
	return nil
}

// ShowPage sets a page's visibility to "true" (in addition to any other pages
// which are already visible).
func (p *Pages) ShowPage(name string) *Pages {
	//p.Lock()
	//defer p.Unlock()

	for _, page := range p.pages {
		if page.Name == name {
			page.Visible = true
			p.curr = page
			if p.changed != nil {
				p.changed()
			}
			break
		}
	}
	if p.HasFocus() {
		p.Focus(p.setFocus)
	}
	return p
}

// HidePage sets a page's visibility to "false".
func (p *Pages) HidePage(name string) *Pages {
	//p.Lock()
	//defer p.Unlock()

	for _, page := range p.pages {
		if page.Name == name {
			page.Visible = false
			if p.changed != nil {
				p.changed()
			}
			break
		}
	}
	if p.HasFocus() {
		p.Focus(p.setFocus)
	}
	return p
}

// SendToFront changes the order of the pages such that the page with the given
// name comes last, causing it to be drawn last with the next update (if
// visible).
func (p *Pages) SendToFront(name string) *Pages {
	{
		//p.Lock()
		//defer p.Unlock()

		for index, page := range p.pages {
			if page.Name == name {
				if index < len(p.pages)-1 {
					p.pages = append(append(p.pages[:index], p.pages[index+1:]...), page)
				}
				if page.Visible && p.changed != nil {
					p.changed()
				}
				break
			}
		}
	}

	if p.HasFocus() {
		p.Focus(p.setFocus)
	}
	return p
}

// SendToBack changes the order of the pages such that the page with the given
// name comes first, causing it to be drawn first with the next update (if
// visible).
func (p *Pages) SendToBack(name string) *Pages {
	//p.Lock()
	//defer p.Unlock()

	for index, pg := range p.pages {
		if pg.Name == name {
			if index > 0 {
				p.pages = append(append([]*Page{pg}, p.pages[:index]...), p.pages[index+1:]...)
			}
			if pg.Visible && p.changed != nil {
				p.changed()
			}
			break
		}
	}
	if p.HasFocus() {
		p.Focus(p.setFocus)
	}
	return p
}

// HasFocus returns whether or not this primitive has focus.
func (p *Pages) HasFocus() bool {
	//p.RLock()
	//defer p.RUnlock()

	for _, page := range p.pages {
		if page.Item.GetFocusable().HasFocus() {
			return true
		}
	}

	/*
		if p.curr != nil && p.curr.Item.GetFocusable().HasFocus() {
			return true
		}
	*/
	return false
}

// Focus is called by the application when the primitive receives focus.
func (p *Pages) Focus(delegate func(p Primitive)) {
	p.Lock()
	p.setFocus = delegate
	p.Unlock()

	p.RLock()
	defer p.RUnlock()

	var topItem Primitive
	for _, page := range p.pages {
		page.RLock()
		if page.Visible {
			topItem = page.Item
		}
		page.RUnlock()
	}

	if topItem != nil {
		delegate(topItem)
	}
}

// SwitchToPage sets a page's visibility to "true" and all other pages'
// visibility to "false".
func (p *Pages) SwitchToPage(name string, context map[string]interface{}) *Pages {
	{ // lock scope
		//p.RLock()
		//defer p.RUnlock()

		//p.curr.RLock()
		//defer p.curr.RUnlock()

		if p.curr != nil && p.curr.Name == name {
			p.curr.Item.Refresh(context)
			return p
		}
	}

	for _, page := range p.pages {
		page.Lock()

		if page.Name == name {
			page.Visible = true
			page.Item.Mount(context)

			p.Lock()
			if p.curr != nil {
				p.curr.Item.Unmount()
			}
			p.curr = page
			p.Unlock()
			context["currPage"] = p.curr
		} else {
			page.Visible = false
		}

		page.Unlock()
	}

	p.RLock()
	if p.changed != nil {
		p.changed()
	}
	p.RUnlock()

	if p.HasFocus() {
		p.Focus(p.setFocus)
	}
	return p
}

// Draw draws this primitive onto the screen.
func (p *Pages) Draw(screen tcell.Screen) {
	p.RLock()
	defer p.RUnlock()

	for _, page := range p.pages {
		page.RLock()

		if !page.Visible {
			page.RUnlock()
			continue
		}
		if page.Resize {
			x, y, width, height := p.GetInnerRect()
			page.Item.SetRect(x, y, width, height)
		}
		page.Item.Draw(screen)
		page.RUnlock()
	}
}
