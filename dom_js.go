// +build js

package dom

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/gopherjs/gopherjs/js"
	"github.com/vecty/dom/unsafe"
)

// Apply applies the changes between the old DOM and new DOM directly to the
// browser's DOM.
func Apply(old, new *DOM) {
	if old == nil || new == nil {
		return
	}
	old.root.applyDiff(new.root, js.Global.Get("document").Get("documentElement"))
}

// applyDiff calculates and applies the difference between the old element, e,
// and the new element, n. The DOM element to which differences should be
// applied is the elem parameter.
func (e *Elem) applyDiff(n *Elem, elem *js.Object) {
	// Collect all element operations into groups.
	// TODO(slimsag): pass around a struct for these.
	oldAttrs, oldTags, oldTexts, oldUnsafeHTMLs, oldUnsafeTexts, oldChildren, oldIndexes := e.collect()
	newAttrs, newTags, newTexts, newUnsafeHTMLs, newUnsafeTexts, newChildren, newIndexes := n.collect()

	if !reflect.DeepEqual(oldIndexes, newIndexes) {
		panic("unexpected state for indexes")
	}
	indexes := oldIndexes

	// TODO(slimsag): move into collect
	if len(oldTags) != 1 || len(newTags) != 1 {
		panic("expected one tag only")
	}

	replacement := false
	if len(oldChildren) != len(newChildren) {
		// TODO(slimsag): do not resort to replacement here in all cases. Identify
		// situations where we can detect an element/node insert/append/delete
		// operation.
		replacement = true
	}

	if oldTags[0] != newTags[0] {
		replacement = true
	}

	if !replacement {
		e.applyAttrs(n, elem, oldAttrs, newAttrs)
		e.applyContent(n, elem, oldTexts, newTexts, oldUnsafeHTMLs, newUnsafeHTMLs, oldUnsafeTexts, newUnsafeTexts)
	}

	if replacement {
		panic("no replacement!")
		// Encode the element to HTML.
		buf := bytes.NewBuffer(nil)
		if err := n.encode(1, buf, DefaultEncodeOpts); err != nil {
			panic("unexpected: " + err.Error())
		}

		// Create a browser DOM element.
		tmpDiv := js.Global.Get("document").Call("createElement", "div")
		tmpDiv.Set("innerHTML", buf.String())
		newElem := tmpDiv.Get("children").Index(0)

		// Save current element state.
		//input := saveInput(elem, newElem) // TODO(slimsag): remove this.
		activeElement := saveActiveElement(elem, newElem)
		textSelection := saveTextSelection(elem, newElem, false)

		// Replace the child element.
		elem.Get("parentNode").Call("replaceChild", newElem, elem)

		// Restore previous element state.
		activeElement.restore()
		textSelection.restore()
		//input.restore() // TODO(slimsag): remove this.
		return
	}

	for i, oldChild := range oldChildren {
		newChild := newChildren[i]
		oldChild.applyDiff(newChild, elem.Get("children").Index(indexes[i]))
	}
	return
}

// applyAttrs calculates the different between the attributes and applies them.
func (oldElem *Elem) applyAttrs(newElem *Elem, elem *js.Object, oldAttrsSlice, newAttrsSlice []AttrOp) {
	// TODO(slimsag): move map composition into collect.
	oldAttrs := make(map[string]string)
	for _, oldAttr := range oldAttrsSlice {
		oldAttrs[oldAttr.Key] = oldAttr.Value
	}
	newAttrs := make(map[string]string)
	for _, newAttr := range newAttrsSlice {
		newAttrs[newAttr.Key] = newAttr.Value
	}

	// Handle the case where an existing attribute is given a new value, or a new
	// attribute is added.
	for newAttrKey, newAttrValue := range newAttrs {
		elem.Call("setAttribute", newAttrKey, newAttrValue)
	}

	// Handle the case where an old attribute is removed.
	for oldAttrKey := range oldAttrs {
		if _, ok := newAttrs[oldAttrKey]; !ok {
			// A old attribute is removed.
			elem.Call("removeAttribute", oldAttrKey)
		}
	}
}

func (oldElem *Elem) applyContent(newElem *Elem, elem *js.Object, oldTexts, newTexts []TextOp, oldUnsafeHTMLs, newUnsafeHTMLs []unsafe.HTMLOp, oldUnsafeTexts, newUnsafeTexts []unsafe.TextOp) {
	equal := reflect.DeepEqual(oldTexts, newTexts) && reflect.DeepEqual(oldUnsafeHTMLs, newUnsafeHTMLs) && reflect.DeepEqual(oldUnsafeTexts, newUnsafeTexts)
	if equal {
		return
	}

	//textSelection := saveTextSelection(elem, elem, false)
	//defer textSelection.restore()

	// TODO(slimsag): can't properly identify nodes when page is indented.
	childNodeIndex := 0
	lastChildNodeIndex := -1
	var last interface{}

	newElem.eachOp(func(op Op) bool {
		//fmt.Println("..")
		//fmt.Printf("[op:%d] %T %+v\n", childNodeIndex, op, op)

		childNode := elem.Get("childNodes").Index(childNodeIndex)
		if lastChildNodeIndex != childNodeIndex {
			lastChildNodeIndex = childNodeIndex
			//println(elem)
			//println(childNode)
			//println("before clear", childNode.Get("nodeValue"))
			childNode.Set("nodeValue", "")
		}

		switch v := op.(type) {
		case unsafe.HTMLOp:
			textSelection := saveTextSelection(childNode, childNode, true)
			childNode.Set("innerHTML", v.HTML)
			textSelection.restore()
			if last != nil {
				if _, ok := last.(unsafe.HTMLOp); !ok {
					childNodeIndex++
				}
			}
			last = v

		case unsafe.TextOp:
			// We have a DOM Text node which needs plain text, so we must decode the
			// HTML entities.
			//
			// TODO(slimsag): write our own decoder for HTML entities.
			tmp := js.Global.Get("document").Call("createElement", "textarea")
			tmp.Set("innerHTML", string(v))

			childNode.Set("nodeValue", childNode.Get("nodeValue").String()+tmp.Get("value").String())
			if last != nil {
				if _, ok := last.(unsafe.TextOp); !ok {
					if _, ok := last.(TextOp); !ok {
						childNodeIndex++
					}
				}
			}
			last = v

		case TextOp:
			childNode.Set("nodeValue", childNode.Get("nodeValue").String()+string(v))
			if last != nil {
				if _, ok := last.(unsafe.TextOp); !ok {
					if _, ok := last.(TextOp); !ok {
						childNodeIndex++
					}
				}
			}
			last = v

		case *Elem:
			// TODO(slimsag): anything to handle here? Aren't elems handled by applyDiff?
			childNodeIndex++
		}
		return true
	})
}

// collect goes over each operation in the element and collects them into typed
// groups.
func (e *Elem) collect() (attrs []AttrOp, tags []TagOp, texts []TextOp, unsafeHTMLs []unsafe.HTMLOp, unsafeTexts []unsafe.TextOp, elems []*Elem, elemIndexes []int) {
	childIndex := 0
	e.eachOp(func(op Op) bool {
		switch v := op.(type) {
		case AttrOp:
			attrs = append(attrs, v)
		case TagOp:
			tags = append(tags, v)
		case TextOp:
			texts = append(texts, v)
		case unsafe.HTMLOp:
			unsafeHTMLs = append(unsafeHTMLs, v)
			childIndex++
		case unsafe.TextOp:
			unsafeTexts = append(unsafeTexts, v)
		case *Elem:
			elems = append(elems, v)
			elemIndexes = append(elemIndexes, childIndex)
			childIndex++
		case nil:
			// ignore it
		default:
			panic(fmt.Sprintf("dom: invalid op %T", v))
		}
		return true
	})
	return
}
