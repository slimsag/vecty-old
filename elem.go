package dom

import (
	"fmt"
	"html"
	"io"
	"strings"

	"github.com/vecty/dom/unsafe"
)

// Op is a DOM operation of any sort.
type Op interface{}

// AttrOp is an operation to set an attribute to a value.
type AttrOp struct {
	Key, Value string
}

// Attr returns an operation that sets an attribute to a value.
func Attr(key, value string) AttrOp {
	return AttrOp{key, value}
}

// TagOp is an operation to set the tag name of an element.
type TagOp string

// Tag returns an operation that sets the tag name of an element.
func Tag(name string) TagOp {
	return TagOp(name)
}

// TextOp is an operation to set the inner text of an element.
type TextOp string

// Text returns an operation that sets the inner text of an element.
func Text(text string) TextOp {
	return TextOp(text)
}

// Elem represents a DOM element.
type Elem struct {
	// ops is the operations used to compose this element.
	ops []Op
}

// New returns a new set of DOM operations.
func New(ops ...Op) *Elem {
	elem := &Elem{ops: ops}
	elem.validate()
	return elem
}

// findTag goes through each operation to find the element's tag name.
func (e *Elem) findTag() string {
	found := ""
	e.eachOp(func(op Op) bool {
		if tag, ok := op.(TagOp); ok {
			found = string(tag)
			return false
		}
		return true
	})
	return found
}

// findHasUnsafeHTML goes through each operation to find out if the element has
// unsafe HTML content or not.
func (e *Elem) findHasUnsafeHTML() bool {
	found := false
	e.eachOp(func(op Op) bool {
		if _, ok := op.(unsafe.HTMLOp); ok {
			found = true
			return false
		}
		return true
	})
	return found
}

func (e *Elem) eachOp(each func(op Op) bool) {
	for _, op := range e.ops {
		if ops, ok := op.([]Op); ok {
			for _, op := range ops {
				if !each(op) {
					return
				}
			}
		} else {
			if !each(op) {
				return
			}
		}
	}
}

// validate validates the element.
func (e *Elem) validate() {
	var (
		tag                                                    string
		haveText, haveUnsafeHTML, haveUnsafeText, haveChildren bool
	)
	e.eachOp(func(op Op) bool {
		switch v := op.(type) {
		case AttrOp:
			// TODO(slimsag): validate attribute keys/values?
		case TagOp:
			if tag != "" {
				panic("dom: must specify only one element tag")
			}
			if v == "" {
				panic("dom: element tag must be a non-empty string")
			}
			if strings.ToLower(string(v)) != string(v) {
				panic("dom: element tags must be lowercase")
			}
			tag = string(v)

		case unsafe.HTMLOp:
			haveUnsafeHTML = true
		case unsafe.TextOp:
			haveUnsafeText = true
		case TextOp:
			haveText = true
		case *Elem:
			haveChildren = true
			v.validate()
		case nil:
			// ignore it
		default:
			panic(fmt.Sprintf("dom: invalid op %T", v))
		}
		return true
	})
	if tag == "" {
		panic("dom: element has no tag name")
	}
	if isVoid(tag) && (haveText || haveUnsafeHTML || haveUnsafeText || haveChildren) {
		panic("dom: void elements may not have content")
	}
}

func (e *Elem) encode(depth int, w io.Writer, opts *EncodeOpts) error {
	// Build an indention string.
	indent := ""
	for i := 0; i < depth; i++ {
		indent += opts.Indent
	}

	// Write starting tag.
	tag := html.EscapeString(e.findTag())
	if _, err := fmt.Fprintf(w, "%s<%s", opts.Prefix+indent, tag); err != nil {
		return err
	}
	var attrs []string
	e.eachOp(func(op Op) bool {
		if attr, ok := op.(AttrOp); ok {
			attrs = append(attrs, fmt.Sprintf(`%s="%s"`, html.EscapeString(attr.Key), html.EscapeString(attr.Value)))
		}
		return true
	})
	if len(attrs) > 0 {
		if _, err := fmt.Fprintf(w, " %s", strings.Join(attrs, " ")); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintf(w, ">%s", opts.Newline); err != nil {
		return err
	}

	// Only handle content or write an ending tag if the element is not a void
	// one.
	if !isVoid(tag) {
		// Handle content.
		var err error
		e.eachOp(func(op Op) bool {
			switch v := op.(type) {
			case unsafe.HTMLOp:
				switch v.Display {
				case "block":
					if _, err = fmt.Fprintf(w, "%s<div>%s</div>%s", opts.Prefix+indent+opts.Indent, string(v.HTML), opts.Newline); err != nil {
						return false
					}
				case "inline":
					if _, err = fmt.Fprintf(w, "%s<span>%s</span>%s", opts.Prefix+indent+opts.Indent, string(v.HTML), opts.Newline); err != nil {
						return false
					}
				case "inline-block":
					if _, err = fmt.Fprintf(w, `%s<span style="display: inline-block;">%s</span>%s`, opts.Prefix+indent+opts.Indent, string(v.HTML), opts.Newline); err != nil {
						return false
					}
				default:
					panic("dom: unsafe.HTML with invalid display value")
				}

			case unsafe.TextOp:
				if _, err = fmt.Fprintf(w, "%s", opts.Prefix+indent+opts.Indent+string(v)+opts.Newline); err != nil {
					return false
				}

			case TextOp:
				if _, err = fmt.Fprintf(w, "%s", opts.Prefix+indent+opts.Indent); err != nil {
					return false
				}
				if _, err = fmt.Fprintf(w, "%s", html.EscapeString(string(v))); err != nil {
					return false
				}
				if _, err = fmt.Fprintf(w, "%s", opts.Newline); err != nil {
					return false
				}
			case *Elem:
				if err = v.encode(depth+1, w, opts); err != nil {
					return false
				}
			}
			return true
		})
		if err != nil {
			return err
		}

		// Write ending tag.
		if _, err = fmt.Fprintf(w, "%s%s</%s>", opts.Prefix, indent, tag); err != nil {
			return err
		}
	}

	if depth != 0 {
		if _, err := fmt.Fprintf(w, "%s", opts.Newline); err != nil {
			return err
		}
	}
	return nil
}

// isVoid tells if the given element tag (which must be lowercase) is considered
// a 'void element' (i.e. content-less) according to the HTML spec:
//
// http://www.w3.org/TR/html-markup/syntax.html#syntax-elements
//
func isVoid(tag string) bool {
	switch strings.ToLower(tag) {
	case "area", "base", "br", "col", "command", "embed", "hr", "img", "input", "keygen", "link", "meta", "param", "source", "track", "wbr":
		return true
	default:
		return false
	}
}
