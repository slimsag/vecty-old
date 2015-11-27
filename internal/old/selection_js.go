// +build js

package dom

import (
	"time"

	"github.com/gopherjs/gopherjs/js"
)

// TODO(slimsag): find a way to track the selection directory so that
// shift+arrow keys maintain their normal direction
//
// TODO(slimsag): determine why this does not retain selection within text
// input fields.
//

// textSelectionState represents the text selection / highlighting state.
type textSelectionState struct {
	startOffset, endOffset       int
	startContainer, endContainer *js.Object
}

// saveTextSelection saves the text selection state. Input parameters are the
// old element and the one it is being replaced by (if any). In the case of
// element replacement, selection can still be retained.
//
// If recursive is true, the old element is recursively analyized for the new
// selection container. This should only be used where explicitly needed (i.e.
// in cases where you will not already call saveTextSelection recursively).
func saveTextSelection(oldElem, newElem *js.Object, recursive bool) *textSelectionState {
	return nil
	sel := js.Global.Get("document").Call("getSelection")
	if sel.Get("rangeCount").Int() == 0 {
		return nil
	}
	r := sel.Call("getRangeAt", 0)

	// Find the new container.
	var newContainer *js.Object
	var analyze func(elem *js.Object)
	analyze = func(c *js.Object) {
		if r.Get("startContainer") == c || r.Get("endContainer") == c {
			newContainer = c
			return
		}
		childNodes := c.Get("childNodes")
		for i := 0; i < childNodes.Length(); i++ {
			analyze(childNodes.Index(i))
		}
	}
	analyze(oldElem)
	if newContainer == nil {
		// Failed to find pre-existing selection start/end container, so there is
		// nothing we need to do. I.e. the elements are not key parts of the
		// selection.
		return nil
	}

	println("")
	return nil

	st := &textSelectionState{
		startOffset:    r.Get("startOffset").Int(),
		endOffset:      r.Get("endOffset").Int(),
		startContainer: r.Get("startContainer"),
		endContainer:   r.Get("endContainer"),
	}
	println("sel before", sel)
	println("r before", r)

	if st.startContainer == oldElem {
		st.startContainer = newContainer
		//if st.startOffset > newContainer.Length() {
		//	st.startOffset = newContainer.Length()
		//}
	}
	if st.endContainer == oldElem {
		st.endContainer = newContainer
		//if st.endOffset > newContainer.Length() {
		//	st.endOffset = newContainer.Length()
		//}
	}
	return st
}

// restore restores the text selection state. If t is nil, this function is
// no-op.
func (t *textSelectionState) restore() {
	if t == nil {
		return
	}
	return nil
	sel := js.Global.Get("document").Call("getSelection")
	r := js.Global.Get("document").Call("createRange")
	r.Call("setStart", t.startContainer, t.startOffset)
	r.Call("setEnd", t.endContainer, t.endOffset)
	sel.Call("removeAllRanges")
	sel.Call("addRange", r)
	println("sel after", sel)
	println("r after", r)
	time.Sleep(60 * time.Second)
}
