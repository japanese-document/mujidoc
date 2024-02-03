package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark"
	"golang.org/x/net/html"
	"golang.org/x/sync/errgroup"
)

type Header struct {
	Name  string `json:"name,omitempty"`
	Order int    `json:"order,omitempty"`
}

type Meta struct {
	Header Header `json:"header,omitempty"`
	Order  int    `json:"order,omitempty"`
	Date   string `json:"date,omitempty"`
}

type Page struct {
	Meta  Meta   `json:"meta,omitempty"`
	Title string `json:"title,omitempty"`
	URL   string `json:"url,omitempty"`
}

type IndexItem struct {
	Name  string          `json:"name,omitempty"`
	Pages []IndexItemPage `json:"pages,omitempty"`
}

type IndexItemPage struct {
	Title string `json:"title,omitempty"`
	URL   string `json:"url,omitempty"`
}

type IndexItemPagesMap map[int]IndexItemPage

var markdown goldmark.Markdown

func init() {
	markdown = NewMarkdown()
}

// CreateHash generates a hash value from the given text.
// It replaces specific characters with "_", "-_", and "_-" to create the hash.
func CreateHash(text string) string {
	text = hashRe.ReplaceAllString(text, "_")
	text = lessThanRe.ReplaceAllString(text, "-_")
	text = greaterThanRe.ReplaceAllString(text, "_-")
	return text
}

// GetMetaAndMd extracts metadata and markdown text from the content of a markdown file.
// The content is split by "---", where the first part is interpreted as JSON format metadata, and the second part as markdown text.
func GetMetaAndMd(content string) (Meta, string, error) {
	parts := SEPARATOR.Split(content, 2)
	if len(parts) != 2 {
		return Meta{}, "", errors.WithStack(errors.New("invalid content format"))
	}

	var meta Meta
	err := json.Unmarshal([]byte(parts[0]), &meta)
	if err != nil {
		return Meta{}, "", errors.WithStack(err)
	}

	md := strings.TrimSpace(parts[1])
	return meta, md, nil
}

// CreateTitle generates a title from the markdown text.
// It uses the first heading (# Title) as the title.
func CreateTitle(md string) string {
	start := len("# ")
	end := strings.Index(md, "\n")
	if end == -1 {
		return md[start:]
	}
	return md[start:end]
}

// CreateURL generates a URL from the specified directory and file name.
// It constructs the URL based on the source directory and base URL obtained from environment variables.
func CreateURL(dir, name string) string {
	trimmedDir := strings.TrimPrefix(dir, os.Getenv("SOURCE_DIR"))
	trimmedDir = strings.TrimPrefix(trimmedDir, "/")
	return os.Getenv("BASE_URL") + "/" + trimmedDir + "/" + name + ".html"
}

// GetDirAndName extracts the directory path and file name (without extension) from a file path.
func GetDirAndName(path string) (string, string) {
	base := filepath.Base(path)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	dir := filepath.Dir(path)
	return dir, name
}

// CreatePageData generates page data from a markdown file name.
// It parses the file content to extract metadata, title, and URL, and returns a Page struct containing these.
func CreatePageData(markDownFileName string) (*Page, error) {
	content, err := os.ReadFile(markDownFileName)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, md, err := GetMetaAndMd(string(content))
	if err != nil {
		return nil, err
	}

	title := CreateTitle(md)
	dir, name := GetDirAndName(markDownFileName)
	url := CreateURL(dir, name)

	page := &Page{
		Meta:  meta,
		Title: title,
		URL:   url,
	}

	return page, nil
}

// extractTextNodes extracts text nodes from an HTML string.
func extractTextNodes(htmlStr string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return "", errors.WithStack(err)
	}

	var buf bytes.Buffer
	var stack []*html.Node
	stack = append(stack, doc)

	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
		}

		// 子ノードを逆順でスタックに追加
		for c := n.LastChild; c != nil; c = c.PrevSibling {
			stack = append(stack, c)
		}
	}

	return buf.String(), nil
}

// CreateDescription generates a description from an HTML string.
// It extracts text nodes and removes newlines and specific characters, then uses the first 300 characters to generate the description.
func CreateDescription(htmlStr string) (string, error) {
	text, err := extractTextNodes(htmlStr)
	if err != nil {
		return "", err
	}

	result := strings.ReplaceAll(text, "\n", "")
	result = strings.ReplaceAll(result, "\"", "&quot;")
	// removing <a>#</a>
	result = strings.TrimPrefix(result, "#")
	if len(result) > 300 {
		result = result[:300]
	}

	return result, nil
}

// CreateHTML generates an HTML page using the given parameters.
// It applies the parameters to the layout template and returns the completed HTML string.
func CreateHTML(layout, title, body, description, url, cssPath, indexMenu, headerList string) string {
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("class").Matching(regexp.MustCompile(`^(anchor|Link)$`)).OnElements("a")
	p.AllowAttrs("class").Matching(regexp.MustCompile(`^(index-menu|header-list)$`)).OnElements("nav")
	p.AllowAttrs("class").Matching(regexp.MustCompile(`^(h1|h2|h3|h4)$`)).OnElements("p")
	html := TITLE.ReplaceAllString(layout, p.Sanitize(title))
	html = DESCRIPTION.ReplaceAllString(html, p.Sanitize(description))
	html = strings.Replace(html, URL, p.Sanitize(url), 1)
	html = strings.Replace(html, CSS, p.Sanitize(cssPath), 1)
	html = strings.Replace(html, INDEX, p.Sanitize(indexMenu), 1)
	html = strings.Replace(html, HEADER, p.Sanitize(headerList), 1)
	html = strings.Replace(html, BODY, p.Sanitize(body), 1)
	return html
}

// CreateIndexItems generates a slice of IndexItem from a slice of Pages.
// It organizes pages by header order and creates IndexItem structs containing each header and page titles and URLs.
func CreateIndexItems(pages []*Page) ([]IndexItem, error) {
	indexItemsMap := map[int]struct {
		Name  string
		Pages IndexItemPagesMap
	}{}
	for _, page := range pages {
		headerOrder := page.Meta.Header.Order
		headerName := page.Meta.Header.Name

		if _, exists := indexItemsMap[headerOrder]; !exists {
			indexItemsMap[headerOrder] = struct {
				Name  string
				Pages IndexItemPagesMap
			}{
				Name:  headerName,
				Pages: IndexItemPagesMap{},
			}
		}

		if indexItemsMap[headerOrder].Name != headerName {
			return nil, errors.WithStack(fmt.Errorf("Header already exists. Existing item: %#v, Current page: %#v", indexItemsMap[headerOrder], page))
		}

		pageOrder := page.Meta.Order
		if _, exists := indexItemsMap[headerOrder].Pages[pageOrder]; exists {
			return nil, errors.WithStack(fmt.Errorf("Page already exists. Existing page: %#v, Current page: %#v", indexItemsMap[headerOrder].Pages[pageOrder], page))
		}

		indexItemsMap[headerOrder].Pages[pageOrder] = IndexItemPage{
			Title: page.Title,
			URL:   page.URL,
		}
	}

	// マップをスライスに変換
	_indexItems, err := MapToSlice(indexItemsMap)
	if err != nil {
		return []IndexItem{}, err
	}
	indexItems, err := Map(_indexItems, func(item struct {
		Name  string
		Pages IndexItemPagesMap
	}, _ int) (IndexItem, error) {
		indexItem := IndexItem{
			Name:  item.Name,
			Pages: []IndexItemPage{},
		}
		pages, err := MapToSlice(item.Pages)
		if err != nil {
			return IndexItem{}, err
		}
		indexItem.Pages = pages
		return indexItem, nil
	})
	if err != nil {
		return []IndexItem{}, err
	}
	return indexItems, nil
}

// CreateIndexMenu generates an HTML navigation menu from a slice of IndexItem.
func CreateIndexMenu(items []IndexItem) string {
	var menu strings.Builder
	menu.WriteString("<nav class=\"index-menu\">")

	for _, item := range items {
		menu.WriteString(fmt.Sprintf("\n<details open>\n<summary>%s</summary>", html.EscapeString(item.Name)))
		for _, page := range item.Pages {
			menu.WriteString(fmt.Sprintf("\n<p><a href=\"%s\">%s</a></p>", html.EscapeString(page.URL), html.EscapeString(page.Title)))
		}
		menu.WriteString("\n</details>")
	}

	menu.WriteString("\n</nav>")
	return menu.String()
}

// IsHeader determines if a given text line is a markdown header.
func IsHeader(line string) bool {
	for i := 2; i <= 5; i++ {
		if strings.HasPrefix(line, strings.Repeat("#", i-1)+" ") {
			return true
		}
	}
	return false
}

// CreatePage generates the HTML for an individual page using the given parameters.
func CreatePage(layout, md, title, url, indexMenu, headerList string) (string, error) {
	var buf bytes.Buffer
	if err := markdown.Convert([]byte(md), &buf); err != nil {
		return "", errors.WithStack(err)
	}
	body := buf.String()

	description, err := CreateDescription(body)
	if err != nil {
		return "", err
	}

	html := CreateHTML(layout, title, body, description, url, os.Getenv("CSS_PATH"), indexMenu, headerList)
	return html, nil
}

// CreateHeaderList generates HTML for a header list from markdown text.
func CreateHeaderList(md string) (string, error) {
	lines := strings.Split(md, "\n")
	filtered, err := Filter(lines, func(line string, _ int) (bool, error) {
		return IsHeader(line), nil
	})
	if err != nil {
		return "", err
	}
	headers := []string{}
	for _, line := range filtered {
		for i := 2; i <= 5; i++ {
			if strings.HasPrefix(line, strings.Repeat("#", i-1)+" ") {
				headerContent := strings.TrimSpace(line[i:])
				href := CreateHash(headerContent)
				var html bytes.Buffer
				if err := markdown.Convert([]byte(headerContent), &html); err != nil {
					return "", errors.WithStack(err)
				}
				text, err := extractTextNodes(html.String())
				if err != nil {
					return "", err
				}
				headers = append(headers, fmt.Sprintf(`<p class="h%d"><a href="#%s">%s</a></p>`, i-1, href, text))
			}
		}
	}
	result := fmt.Sprintf(`<nav class="header-list">%s</nav>`, strings.Join(headers, "\n"))
	return result, nil
}

// CreateIndexPage generates the HTML for an index page from index items.
func CreateIndexPage(layout string, indexItems []IndexItem) (string, error) {
	var builder strings.Builder

	// ヘッダーを追加
	builder.WriteString("# " + os.Getenv("INDEX_PAGE_HEADER") + "\n")

	// 各IndexItemに対して処理
	for _, item := range indexItems {
		builder.WriteString("\n## " + item.Name + "\n")
		for _, page := range item.Pages {
			builder.WriteString(fmt.Sprintf("* [%s](%s)\n", page.Title, page.URL))
		}
	}

	var body bytes.Buffer
	if err := markdown.Convert([]byte(builder.String()), &body); err != nil {
		return "", errors.WithStack(err)
	}
	return CreateHTML(layout, os.Getenv("INDEX_PAGE_TITLE"), body.String(), os.Getenv("INDEX_PAGE_DESCRIPTION"), os.Getenv("BASE_URL"), os.Getenv("CSS_PATH"), "", ""), nil
}

// createPageTask returns a task that generates page data from a specified markdown file.
func createPageTask(markDownFileName string, pages []*Page, index int) func() error {
	return func() error {
		page, err := CreatePageData(markDownFileName)
		if err != nil {
			return err
		}
		pages[index] = page
		return nil
	}
}

// CreatePages asynchronously generates a slice of Page data from multiple markdown files.
func CreatePages(markDownFileNames []string) ([]*Page, error) {
	var g errgroup.Group
	pages := make([]*Page, len(markDownFileNames))

	for i, fileName := range markDownFileNames {
		task := createPageTask(fileName, pages, i)
		g.Go(task)
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return pages, nil
}
