// +build js

package dom

import "github.com/gopherjs/gopherjs/js"

type activeElementState struct {
	newElem *js.Object
}

func saveActiveElement(elem, newElem *js.Object) *activeElementState {
	if js.Global.Get("document").Get("activeElement") != elem {
		return nil
	}
	return &activeElementState{newElem: newElem}
}

func (s *activeElementState) restore() {
	if s == nil {
		return
	}
	s.newElem.Call("focus")
}
