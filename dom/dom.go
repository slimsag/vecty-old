package dom

import (
	"fmt"
	"io"
)

// EncodeOpts represents options for encoding the DOM to HTML.
//
// Note that although the addition of indention and newlines makes the generated
// HTML code look prettier, due to the way browsers handle whitespace in HTML
// code it also causes browsers to render your inline elements differently. For
// example see:
//
// https://jsfiddle.net/xwqp8wbw/2/
//
// In general, always stick to no indention/newlines except when trying to debug
// issues with the HTML generation process.
type EncodeOpts struct {
	// DocType is the document type string.
	DocType string

	// Prefix is a string prefixed onto each line.
	Prefix string

	// Indent is the indention string (e.g. "  " or "\t").
	Indent string

	// Newline is the newline string (e.g. "\n").
	Newline string
}

// DefaultEncodeOpts is the default encoding options used.
var DefaultEncodeOpts = &EncodeOpts{
	DocType: "<!DOCTYPE html>",
	Prefix:  "",
	Indent:  "",
	Newline: "",
}

// NewDOM returns a virtual DOM implementation.
func NewDOM(root *Elem) *DOM {
	return &DOM{
		root: root,
	}
}

// DOM virtually represents the document object model.
type DOM struct {
	root *Elem
}

// Encode encodes this DOM has HTML using the given options, writing the data
// out to the given writer and returning any IO errors.
func (d *DOM) Encode(w io.Writer, opts *EncodeOpts) error {
	// Handle default options.
	if opts == nil {
		cpy := *DefaultEncodeOpts
		opts = &cpy
	}

	// Write the DocType string.
	if _, err := fmt.Fprintf(w, "%s", opts.Prefix+opts.DocType+opts.Newline); err != nil {
		return err
	}

	// Encode the elements.
	return d.root.encode(0, w, opts)
}
