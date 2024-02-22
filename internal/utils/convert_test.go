package utils

import (
	"log"
	"path/filepath"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../../.env.mujidoc")
	if err != nil {
		log.Fatal(err)
	}
}

func TestCreateHash(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "create hash",
			args: args{
				text: `<a>?:&foo=%"\'@><`,
			},
			want: "-_a_-___foo_______--_",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateHash(tt.args.text); got != tt.want {
				t.Errorf("CreateHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateTitle(t *testing.T) {
	type args struct {
		md string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Standard title",
			args: args{md: "# Sample Title\nThis is a Markdown file."},
			want: "Sample Title",
		},
		{
			name: "No newline",
			args: args{md: "# Title Without Newline"},
			want: "Title Without Newline",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateTitle(tt.args.md); got != tt.want {
				t.Errorf("CreateTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateURL(t *testing.T) {
	type args struct {
		dir       string
		name      string
		sourceDir string
		baseURL   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Sub directory",
			args: args{dir: "src/source", name: "index", sourceDir: "src", baseURL: "https://example.com"},
			want: "https://example.com/source/index.html",
		},
		{
			name: "Root directory",
			args: args{dir: "src", name: "index", sourceDir: "src", baseURL: "https://example.com"},
			want: "https://example.com/index.html",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateURL(tt.args.dir, tt.args.name, tt.args.sourceDir, tt.args.baseURL); got != tt.want {
				t.Errorf("CreateURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateDescription(t *testing.T) {
	type args struct {
		htmlStr string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Simple text",
			args: args{htmlStr: "1\n<p>Hello World</p>"},
			want: "Hello World",
		},
		{
			name: "Text with line breaks",
			args: args{htmlStr: "1\n<p>Hello\nWorld</p>"},
			want: "HelloWorld",
		},
		{
			name: "Text with double quotes",
			args: args{htmlStr: "1\n<p>Hello \"World\"</p>"},
			want: `Hello &#34;World&#34;`,
		},
		{
			name: "Complex HTML",
			args: args{htmlStr: "1\n<div><p>Hello</p><p>World</p></div>"},
			want: "HelloWorld",
		},
		{
			name: "Empty string",
			args: args{htmlStr: ""},
			want: "",
		},
		{
			name: "Text exceeding 300 characters",
			args: args{htmlStr: "1\n<p>" + strings.Repeat("a", 500) + "</p>"},
			want: strings.Repeat("a", 300),
		},
		{
			name: "Invalid HTML",
			args: args{htmlStr: "1\n<p>Hello World"},
			want: "Hello World",
		},
		{
			name: "1 line",
			args: args{htmlStr: "foo"},
			want: "foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateDescription(tt.args.htmlStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDescription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreateDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateIndexMenu(t *testing.T) {
	type args struct {
		items []IndexItem
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Empty Items",
			args: args{items: []IndexItem{}},
			want: "<nav class=\"index-menu\">\n</nav>",
		},
		{
			name: "Single Item No Pages",
			args: args{items: []IndexItem{{Name: "Item1"}}},
			want: "<nav class=\"index-menu\">\n<details open>\n<summary>Item1</summary>\n</details>\n</nav>",
		},
		{
			name: "Multiple Items",
			args: args{items: []IndexItem{
				{Name: "Item1", Pages: []IndexItemPage{
					{Title: "Page1", URL: "https://example.com/page1"},
				}},
				{Name: "Item2", Pages: []IndexItemPage{
					{Title: "Page2", URL: "https://example.com/page2"},
				}},
			}},
			want: "<nav class=\"index-menu\">\n<details open>\n<summary>Item1</summary>\n<p><a href=\"https://example.com/page1\">Page1</a></p>\n</details>\n<details open>\n<summary>Item2</summary>\n<p><a href=\"https://example.com/page2\">Page2</a></p>\n</details>\n</nav>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateIndexMenu(tt.args.items); got != tt.want {
				t.Errorf("CreateIndexMenu() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateHTMLFileDir(t *testing.T) {
	type args struct {
		dir       string
		sourceDir string
		outputDir string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestSubdirectory",
			args: args{
				dir:       "/home/user/source/articles",
				sourceDir: "/home/user/source",
				outputDir: "/home/user/output",
			},
			want: filepath.Join("/home/user/output", "articles"),
		},
		{
			name: "TestNoSubdirectory",
			args: args{
				dir:       "/home/user/source",
				sourceDir: "/home/user/source",
				outputDir: "/home/user/output",
			},
			want: "/home/user/output",
		},
		{
			name: "TestNestedSubdirectory",
			args: args{
				dir:       "/home/user/source/articles/2020",
				sourceDir: "/home/user/source",
				outputDir: "/home/user/output",
			},
			want: filepath.Join("/home/user/output", "articles/2020"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateHTMLFileDir(tt.args.dir, tt.args.sourceDir, tt.args.outputDir); got != tt.want {
				t.Errorf("CreateHTMLFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsHeader(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "# is a header",
			args: args{
				line: "# This a header",
			},
			want: true,
		},
		{
			name: "## is a header",
			args: args{
				line: "## This is a header",
			},
			want: true,
		},
		{
			name: "### is a header",
			args: args{
				line: "### This is also a header",
			},
			want: true,
		},
		{
			name: "#### is a header",
			args: args{
				line: "#### This is still a header",
			},
			want: true,
		},
		{
			name: "##### is a header but not detected",
			args: args{
				line: "##### This is not detected as a header by IsHeader",
			},
			want: false,
		},
		{
			name: "Non-header text",
			args: args{
				line: "This is just some text.",
			},
			want: false,
		},
		{
			name: "Empty string",
			args: args{
				line: "",
			},
			want: false,
		},
		{
			name: "There is no space",
			args: args{
				line: "#This a header",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsHeader(tt.args.line); got != tt.want {
				t.Errorf("IsHeader() = %v, want %v for line %v", got, tt.want, tt.args.line)
			}
		})
	}
}

func TestIsCode(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "only code",
			args: args{line: "```"},
			want: true,
		},
		{
			name: "with language",
			args: args{line: "```go"},
			want: true,
		},
		{
			name: "start space",
			args: args{line: " ```"},
			want: false,
		},
		{
			name: "others",
			args: args{line: "foobar"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCode(tt.args.line); got != tt.want {
				t.Errorf("IsCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
