package main

import (
	"fmt"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/vecty/dom/internal/selection"
)

func main() {
	js.Global.Get("document").Call("write", `
<!DOCTYPE html>
<html>
  <head></head>
  <body id="body">
		<p>one two three four five six</p>
		<p id="para">some <strong>paragraph <em>text</em></strong></p>
		<p>seven eight nine ten eleven</p>
  </body>
</html>
`)
	js.Global.Get("document").Call("close")

	n := 0
	go func() {
		for {
			n++
			newElem := js.Global.Get("document").Call("createElement", "p")
			newElem.Call("setAttribute", "id", "para")
			oldElem := js.Global.Get("document").Call("getElementById", "para")
			sel := selection.Save(selection.Opts{
				NewElem: newElem,
				OldElem: oldElem,
			})
			newElem.Set("innerHTML", fmt.Sprintf("some <strong>paragraph <em>text(%d)</em></strong>", n))
			js.Global.Get("document").Call("getElementById", "body").Call("replaceChild", newElem, oldElem)
			sel.Restore()
			time.Sleep(1 * time.Second)
		}
	}()
	println("oh noe!")
}
