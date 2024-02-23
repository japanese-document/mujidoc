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

type PageMeta struct {
	Category string `json:"category,omitempty"`
	Order    int    `json:"order,omitempty"`
	Date     string `json:"date,omitempty"`
}

type Category struct {
	Name  string `json:"name,omitempty"`
	Order int    `json:"order,omitempty"`
}

type Meta struct {
	Category Category `json:"header,omitempty"`
	Order    int      `json:"order,omitempty"`
	Date     string   `json:"date,omitempty"`
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
func GetMetaAndMd(content string, categoryOrders map[string]int) (*Meta, string, error) {
	parts := SEPARATOR.Split(content, 2)
	if len(parts) != 2 {
		return nil, "", errors.WithStack(errors.New("invalid content format"))
	}

	pm := PageMeta{}
	err := json.Unmarshal([]byte(parts[0]), &pm)
	if err != nil {
		return nil, "", errors.WithStack(err)
	}

	categoryOrder, exist := categoryOrders[pm.Category]
	if !exist {
		return nil, "", errors.WithStack(errors.Errorf("%s does not exist in CATGEGORIES", pm.Category))
	}

	meta := &Meta{
		Category: Category{
			Name:  pm.Category,
			Order: categoryOrder,
		},
		Order: pm.Order,
		Date:  pm.Date,
	}
	md := strings.TrimSpace(parts[1])
	return meta, md, nil
}

func GetMd(content string) (string, error) {
	parts := SEPARATOR.Split(content, 2)
	if len(parts) != 2 {
		return "", errors.WithStack(errors.New("invalid content format"))
	}
	md := strings.TrimSpace(parts[1])
	return md, nil
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
func CreateURL(dir, name, sourceDir, baseURL string) string {
	trimmedDir := strings.TrimPrefix(dir, sourceDir)
	trimmedDir = strings.Trim(trimmedDir, "/")
	if trimmedDir == "" {
		return baseURL + "/" + name + ".html"
	}
	return baseURL + "/" + trimmedDir + "/" + name + ".html"
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
func CreatePageData(markDownFileName, sourceDir, baseURL string, categoryOrders map[string]int) (*Page, error) {
	content, err := os.ReadFile(markDownFileName)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta, md, err := GetMetaAndMd(string(content), categoryOrders)
	if err != nil {
		return nil, err
	}

	title := CreateTitle(md)
	dir, name := GetDirAndName(markDownFileName)
	url := CreateURL(dir, name, sourceDir, baseURL)

	page := &Page{
		Meta:  *meta,
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

	start := strings.Index(text, "\n")
	if start != -1 {
		text = text[start:]
	}

	result := strings.ReplaceAll(text, "\n", "")
	result = html.EscapeString(result)
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
		categoryOrder := page.Meta.Category.Order
		categoryName := page.Meta.Category.Name

		if _, exists := indexItemsMap[categoryOrder]; !exists {
			indexItemsMap[categoryOrder] = struct {
				Name  string
				Pages IndexItemPagesMap
			}{
				Name:  categoryName,
				Pages: IndexItemPagesMap{},
			}
		}

		if indexItemsMap[categoryOrder].Name != categoryName {
			return nil, errors.WithStack(fmt.Errorf("header already exists. existing item: %#v, current page: %#v", indexItemsMap[categoryOrder], page))
		}

		pageOrder := page.Meta.Order
		if _, exists := indexItemsMap[categoryOrder].Pages[pageOrder]; exists {
			return nil, errors.WithStack(fmt.Errorf("page already exists. existing page: %#v, current page: %#v", indexItemsMap[categoryOrder].Pages[pageOrder], page))
		}

		indexItemsMap[categoryOrder].Pages[pageOrder] = IndexItemPage{
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

// CreatePage generates the HTML for an individual page using the given parameters.
func CreatePage(layout, md, title, url, cssPath, indexMenu, headerList string) (string, error) {
	var buf bytes.Buffer
	if err := markdown.Convert([]byte(md), &buf); err != nil {
		return "", errors.WithStack(err)
	}
	body := buf.String()

	description, err := CreateDescription(body)
	if err != nil {
		return "", err
	}

	html := CreateHTML(layout, title, body, description, url, cssPath, indexMenu, headerList)
	return html, nil
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

// IsHeader determines if a given text line is a markdown code.
func IsCode(line string) bool {
	return strings.HasPrefix(line, "```")
}

// CreateHeaderList generates HTML for a header list from markdown text.
func CreateHeaderList(md string) (string, error) {
	lines := strings.Split(md, "\n")
	isInCode := false
	filtered, err := Filter(lines, func(line string, _ int) (bool, error) {
		if IsCode(line) {
			isInCode = !isInCode
		}
		// ignore `#` in code tag
		if isInCode {
			return false, nil
		}
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
func CreateIndexPage(layout, baseURL, header, title, description, cssPath string, indexItems []IndexItem) (string, error) {
	var builder strings.Builder

	// ヘッダーを追加
	builder.WriteString("# " + header + "\n")

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
	return CreateHTML(layout, title, body.String(), description, baseURL, cssPath, "", ""), nil
}

// createPageTask returns a task that generates page data from a specified markdown file.
func createPageTask(markDownFileName string, pages []*Page, sourceDir, baseURL string, categoryOrders map[string]int, index int) func() error {
	return func() error {
		page, err := CreatePageData(markDownFileName, sourceDir, baseURL, categoryOrders)
		if err != nil {
			return err
		}
		pages[index] = page
		return nil
	}
}

// CreateCategoryOrders takes a string of categories separated by commas, trims any trailing comma,
// and then splits the string into individual categories. It creates and returns a map where each
// category is a key with its value being the order (index) in which the category appears in the input string.
// If a category appears more than once, it is only added to the map once, with the index of its first occurrence.
//
// Parameters:
// - categories: A string of categories separated by commas. For example, "apple,orange,banana,apple".
//
// Returns:
//   - A map[string]int where keys are unique categories from the input string and values are the indexes
//     at which those categories first appear in the input string.
//
// Note:
// This function does not account for spaces around category names. For example, "apple, orange" will
// treat " orange" (with a leading space) as a distinct category from "orange".
func CreateCategoryOrders(categories string) map[string]int {
	categories = strings.TrimSuffix(categories, ",")
	cs := strings.Split(categories, ",")
	co := map[string]int{}
	for i, v := range cs {
		_, exists := co[v]
		if !exists && v != "" {
			co[v] = i
		}
	}
	return co
}

// CreatePages asynchronously generates a slice of Page data from multiple markdown files.
func CreatePages(markDownFileNames []string, sourceDir, baseURL, categories string) ([]*Page, error) {
	var g errgroup.Group
	pages := make([]*Page, len(markDownFileNames))
	categoryOrders := CreateCategoryOrders(categories)

	for i, fileName := range markDownFileNames {
		task := createPageTask(fileName, pages, sourceDir, baseURL, categoryOrders, i)
		g.Go(task)
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return pages, nil
}

// CreateHTMLFileDir constructs the file path for an HTML file based on the given directory, source directory, and output directory.
// It calculates the subdirectory path by removing the source directory prefix from the given directory path.
// This subdirectory path is then appended to the output directory to form the final file path.
// Parameters:
// - dir: The directory where the original Markdown file resides. It is expected to be a subdirectory of the source directory.
// - sourceDir: The root directory of all source Markdown files.
// - outputDir: The root directory where the resulting HTML files should be saved.
// Returns:
// The constructed file path for the HTML file, which combines the output directory and the subdirectory derived from the given directory, excluding the source directory prefix.
func CreateHTMLFileDir(dir, sourceDir, outputDir string) string {
	subDir := ""
	prefixDirCount := len(sourceDir) + len("/")
	if len(dir) > prefixDirCount {
		subDir = dir[prefixDirCount:]
	}
	return filepath.Join(outputDir, subDir)
}
