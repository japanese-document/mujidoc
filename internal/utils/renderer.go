package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

var parentPathRe = regexp.MustCompile(`^(\.\./)+`)

func isLocalImage(imageDir, destination string) bool {
	destination = parentPathRe.ReplaceAllString(destination, "")
	return strings.HasPrefix(destination, imageDir+"/")
}

type customRenderer struct{}

func (r customRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindLink, r.renderLink)
	reg.Register(ast.KindHeading, r.renderHeading)
	reg.Register(ast.KindImage, r.renderImage)
}

func (r customRenderer) renderLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*ast.Link)
		_, err := w.WriteString(`<a href="`)
		if err != nil {
			return 0, err
		}
		_, err = w.Write(util.EscapeHTML(n.Destination))
		if err != nil {
			return 0, err
		}
		_, err = w.WriteString(`" class="Link">`)
		if err != nil {
			return 0, err
		}
	} else {
		_, err := w.WriteString("</a>")
		if err != nil {
			return 0, err
		}
	}
	return ast.WalkContinue, nil
}

func (r customRenderer) renderHeading(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	heading := node.(*ast.Heading)
	if entering {
		headingID := CreateHash(string(heading.Text(source)))
		startTag := fmt.Sprintf(`<h%d id="%s"><a href="#%s">`, heading.Level, headingID, headingID)
		_, err := w.WriteString(startTag)
		if err != nil {
			return 0, err
		}
	} else {
		endTag := fmt.Sprintf("</a></h%d>\n", heading.Level)
		_, err := w.WriteString(endTag)
		if err != nil {
			return 0, err
		}
	}
	return ast.WalkContinue, nil
}

func (r customRenderer) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.Image)
	destination := string(util.EscapeHTML(n.Destination))
	img := ""
	if isLocalImage(IMAGE_DIR, destination) {
		imagePath := parentPathRe.ReplaceAllString(destination, "")
		path := filepath.Join(os.Getenv("SOURCE_DIR"), imagePath)
		width, height, err := ImageSize(path)
		if err != nil {
			fmt.Printf("%s does not exist\n", path)
			img = fmt.Sprintf(`<img src="%s" alt="%s">`, destination, destination)
		} else {
			img = fmt.Sprintf(`<img loading="lazy" src="%s" alt="%s" width="%d" height="%d">`, destination, destination, width, height)
		}
	} else {
		img = fmt.Sprintf(`<img src="%s" alt="%s">`, destination, destination)
	}
	_, err := w.WriteString(img)
	if err != nil {
		return 0, err
	}
	return ast.WalkSkipChildren, nil
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
