// +build js

package selection

import "github.com/gopherjs/gopherjs/js"

// State represents the selection state.
type State struct {
	// Offsets into the containers at which point the selection starts / ends.
	StartOffset, EndOffset int

	// The container elements or nodes at which the selection starts / ends.
	StartContainer, EndContainer *js.Object
}

// Opts represents options to use when saving the selection state.
type Opts struct {
	// New and old elements for a "swapping" operation. If swapping is to be
	// performed, both elements must be specified.
	NewElem, OldElem *js.Object
}

// Save saves the current selection state for restoration later.
func Save(opt Opts) *State {
	sel := js.Global.Get("document").Call("getSelection")
	if sel.Get("rangeCount").Int() == 0 {
		// No selection.
		return nil
	}

	// We only care about the first selection range. Although the API supports
	// multiple, in practice most browsers only ever use one (dated API from the
	// Netscape-era).
	r := sel.Call("getRangeAt", 0)

	// Create and return the state.
	return &State{
		StartOffset:    r.Get("startOffset").Int(),
		EndOffset:      r.Get("endOffset").Int(),
		StartContainer: r.Get("startContainer"),
		EndContainer:   r.Get("endContainer"),
	}
}

// Restore restores this selection state.
func (s *State) Restore() {
	if s == nil {
		return
	}
	println("startOffset", s.StartOffset)
	println("endOffset", s.EndOffset)
	println("startContainer", s.StartContainer)
	println("endContainer", s.EndContainer)

	sel := js.Global.Get("document").Call("getSelection")
	sel.Call("removeAllRanges")
	r := js.Global.Get("document").Call("createRange")
	r.Call("setStart", s.StartContainer, s.StartOffset)
	r.Call("setEnd", s.EndContainer, s.EndOffset)
	sel.Call("addRange", r)
}

/*
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
*/
