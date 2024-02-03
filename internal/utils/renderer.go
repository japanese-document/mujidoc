package utils

import (
	"fmt"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type customRenderer struct{}

func (r customRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindLink, r.renderLink)
	reg.Register(ast.KindHeading, r.renderHeading)
	reg.Register(ast.KindImage, r.renderImage)
}

func (r customRenderer) renderLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*ast.Link)
		w.WriteString(`<a href="`)
		w.Write(util.EscapeHTML(n.Destination))
		w.WriteString(`" class="Link">`)
	} else {
		w.WriteString("</a>")
	}
	return ast.WalkContinue, nil
}

func (r customRenderer) renderHeading(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	heading := node.(*ast.Heading)
	if entering {
		headingID := CreateHash(string(heading.Text(source)))
		startTag := fmt.Sprintf(`<h%d id="%s"><a href="#%s" class="anchor">#</a>`, heading.Level, headingID, headingID)
		w.WriteString(startTag)
	} else {
		endTag := fmt.Sprintf(`</h%d>`, heading.Level)
		w.WriteString(endTag)
	}
	return ast.WalkContinue, nil
}

func (r customRenderer) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*ast.Image)
		destination := string(n.Destination)
		w.WriteString(`<img src="`)
		w.WriteString(destination)
		w.WriteString(`" alt="`)
		w.WriteString(destination)
		w.WriteString(`">`)
	}
	return ast.WalkContinue, nil
}

// NewMarkdown initializes a new goldmark.Markdown instance with custom rendering logic.
// It includes GitHub Flavored Markdown (GFM) extensions and sets the custom renderer with high priority.
func NewMarkdown() goldmark.Markdown {
	option := goldmark.WithRendererOptions(renderer.WithNodeRenderers(
		util.Prioritized(customRenderer{}, 200),
	))
	markdown := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		option,
	)
	return markdown
}
