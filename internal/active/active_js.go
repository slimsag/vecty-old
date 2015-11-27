// +build js

package active

import "github.com/gopherjs/gopherjs/js"

type State struct {
	// Elem is the active element.
	Elem *js.Object
}

type Opts struct {
	// New and old elements for a "swapping" operation. If swapping is to be
	// performed, both elements must be specified.
	NewElem, OldElem *js.Object
}

func Element() *State {
	if opt.NewElem != nil && opt.OldElem != nil {
		if js.Global.Get("document").Get("activeElement") != opt.NewElem {
			return nil
		}
		panic("swapping not implemented")
	}
	return &activeElementState{newElem: newElem}
}

func (s *activeElementState) restore() {
	if s == nil {
		return
	}
	s.newElem.Call("focus")
}
