package shared

import (
	"fmt"

	v "github.com/vecty/dom"
	"github.com/vecty/dom/unsafe"
)

type Page struct {
	Count    int
	LoadTime string
}

func (p *Page) Render() *v.Elem {
	p.Count++

	var greenText v.Op
	var spaces = []v.Op{unsafe.Text("&nbsp;")}
	if (p.Count % 2) == 0 {
		greenText = v.Attr("style", "background-color: green")
	} else {
		spaces = append(spaces, unsafe.Text("&nbsp;"))
		spaces = append(spaces, unsafe.Text("&nbsp;"))
		spaces = append(spaces, unsafe.Text("&nbsp;"))
		spaces = append(spaces, unsafe.Text("&nbsp;"))
		spaces = append(spaces, unsafe.Text("&nbsp;"))
		spaces = append(spaces, unsafe.Text("&nbsp;"))
		spaces = append(spaces, unsafe.Text("&nbsp;"))
		spaces = append(spaces, unsafe.Text("&nbsp;"))
		spaces = append(spaces, unsafe.Text("&nbsp;"))
		spaces = append(spaces, unsafe.Text("&nbsp;"))
	}

	title := v.New(
		v.Tag("title"),
		v.Text("Hello Vecty!"),
	)
	head := v.New(
		v.Tag("head"),
		title,
	)
	body := v.New(
		v.Tag("body"),
		v.Attr("style", "font-size: 16px;"),
		v.New(
			v.New(
				v.Tag("strong"),
				v.Text("text above"),
			),
			v.Tag("p"),
			unsafe.Text("&nbsp;"),
			v.Text("Some paragraph text! It's <strong>escaped!</strong>"),
			spaces,
			v.New(
				v.Tag("strong"),
				v.Text("text below"),
			),
		),
		unsafe.HTML("block", "But <em>we still have <strong>raw HTML power</strong></em> "+p.LoadTime),
		v.New(
			v.Tag("div"),
			v.New(
				v.Tag("div"),
				v.New(v.Tag("hr")),
				v.New(
					v.Tag("strong"),
					greenText,
					v.Text(fmt.Sprintf("Render: %v", p.Count)),
				),
				v.New(v.Tag("br")),
				v.New(
					v.Tag("button"),
					v.Attr("type", "button"),
					v.Text(fmt.Sprintf("Render: %v", p.Count)),
				),
				v.New(
					v.Tag("input"),
					v.Attr("type", "text"),
					v.Attr("value", fmt.Sprintf("Render: %v", p.Count)),
				),
				v.New(
					v.Tag("span"),
					v.New(
						v.Tag("input"),
						v.Attr("type", "checkbox"),
						v.Attr("value", "car"),
					),
					v.Text(fmt.Sprintf("I have a car: %v", p.Count)),
				),
			),
		),
		v.New(
			v.Tag("script"),
			v.Attr("src", "/assets/script.js"),
		),
	)
	return v.New(
		v.Tag("html"),
		head,
		body,
	)
}
