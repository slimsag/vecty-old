// +build js

package dom

import (
	"fmt"

	"github.com/gopherjs/gopherjs/js"
)
import "github.com/gopherjs/jsbuiltin"

// TODO(slimsag): remove this / require that inputs be modified and not rebuilt?

type inputState struct {
	newElem                                          *js.Object
	value                                            *js.Object
	selectionStart, selectionEnd, selectionDirection *js.Object
}

func saveInput(elem, newElem *js.Object) *inputState {
	if !jsbuiltin.InstanceOf(elem, js.Global.Get("HTMLInputElement")) {
		return nil
	}

	println(elem, newElem)
	fmt.Println(elem.Get("value").Interface())
	fmt.Println(elem.Get("selectionStart").Interface())
	fmt.Println(elem.Get("selectionEnd").Interface())
	fmt.Println(elem.Get("selectionDirection").Interface())

	//println("sleep 20")
	//time.Sleep(20 * time.Second)

	return &inputState{
		newElem:            newElem,
		value:              elem.Get("value"),
		selectionStart:     elem.Get("selectionStart"),
		selectionEnd:       elem.Get("selectionEnd"),
		selectionDirection: elem.Get("selectionDirection"),
	}
}

func (s *inputState) restore() {
	if s == nil {
		return
	}
	s.newElem.Set("value", s.value)
	s.newElem.Set("selectionStart", s.selectionStart)
	s.newElem.Set("selectionEnd", s.selectionEnd)
	s.newElem.Set("selectionDirection", s.selectionDirection)
}
