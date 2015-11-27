// Package unsafe allows placing unsafe content into the DOM.
package unsafe

// HTML is an operation for placing unsafe HTML into the output data stream.
type HTMLOp struct {
	HTML    string
	Display string
}

// HTML returns an operation which inserts arbitrarily unsafe HTML content into
// the output data stream.
//
// You are responsible for ensuring that the content has been safely escaped and
// validated.
//
// Due to the way the rendering algorithm operates, your arbitrary HTML will be
// placed under a singular element with the given CSS display value. For
// example:
//
//  // Block of HTML (wrapped by div):
//  unsafe.HTML("block", "<em>this could be</em><strong>dangerous!</strong>"),
//
//  // Inline HTML (wrapped by span):
//  unsafe.HTML("inline", "<em>this could be</em><strong>dangerous!</strong>"),
//
//  // Inline block of HTML (wrapped by span with CSS display value of
//  // inline-block):
//  unsafe.HTML("inline-block", "<em>this could be</em><strong>dangerous!</strong>"),
//
func HTML(display, html string) HTMLOp {
	return HTMLOp{
		Display: display,
		HTML:    html,
	}
}

// Text is an operation for placing unsafe text content into the output data
// stream.
type TextOp string

// Text returns an operation which inserts arbitrarily unsafe text content into
// the output data stream.
//
// Unlike unsafe.HTML which validates that the HTML is at least under a singular
// DOM element (which is compatible with the rendering algorithm), Text provides
// no such validation. If you insert multiple elements, the rendering algorithm
// will break and may panic.
//
// You are responsible for ensuring that the content has been safely escaped and
// validated.
//
// It can be used to insert HTML entity codes, for example:
//
//  unsafe.Text("&lt;")
//  unsafe.Text("&#60;")
//
func Text(text string) TextOp {
	return TextOp(text)
}
