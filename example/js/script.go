package main

import (
	"fmt"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/vecty/dom"
	"github.com/vecty/dom/example/shared"
)

func main() {
	var (
		page = new(shared.Page)
		old  *dom.DOM
	)

	var render *js.Object
	render = js.MakeFunc(func(this *js.Object, arguments []*js.Object) interface{} {
		go func() {
			start := time.Now()

			new := dom.NewDOM(page.Render())
			dom.Apply(old, new)
			old = new

			page.LoadTime = time.Since(start).String()
			fmt.Println(time.Since(start))

			time.Sleep(3 * time.Second)
			js.Global.Call("requestAnimationFrame", render)
		}()
		return nil
	})
	js.Global.Call("requestAnimationFrame", render)
}
