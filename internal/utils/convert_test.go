package utils

import (
	"log"
	"os"
	"reflect"
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

func TestGetMetaAndMd(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name    string
		args    args
		want    Meta
		want1   string
		wantErr bool
	}{
		{
			name: "get meta and md",
			args: args{
				content: "{\"header\": {\"name\": \"foo\", \"order\": 123}, \"order\": 3, \"date\": \"2023-01-01 01:02:03\"}\n---\n# Foo\nBar",
			},
			want: Meta{
				Header: Header{
					Name:  "foo",
					Order: 123,
				},
				Order: 3,
				Date:  "2023-01-01 01:02:03",
			},
			want1:   "# Foo\nBar",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetMetaAndMd(tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMetaAndMd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMetaAndMd() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetMetaAndMd() got1 = %v, want %v", got1, tt.want1)
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
		dir  string
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Root directory",
			args: args{dir: os.Getenv("SOURCE_DIR") + "/source", name: "index"},
			want: os.Getenv("BASE_URL") + "source/index.html",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateURL(tt.args.dir, tt.args.name); got != tt.want {
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
			args: args{htmlStr: "<p>Hello World</p>"},
			want: "Hello World",
		},
		{
			name: "Text with line breaks",
			args: args{htmlStr: "<p>Hello\nWorld</p>"},
			want: "HelloWorld",
		},
		{
			name: "Text with double quotes",
			args: args{htmlStr: `<p>Hello "World"</p>`},
			want: `Hello &quot;World&quot;`,
		},
		{
			name: "Complex HTML",
			args: args{htmlStr: "<div><p>Hello</p><p>World</p></div>"},
			want: "HelloWorld",
		},
		{
			name: "Empty string",
			args: args{htmlStr: ""},
			want: "",
		},
		{
			name: "Text exceeding 300 characters",
			args: args{htmlStr: "<p>" + strings.Repeat("a", 500) + "</p>"},
			want: strings.Repeat("a", 300),
		},
		{
			name: "Invalid HTML",
			args: args{htmlStr: "<p>Hello World"},
			want: "Hello World",
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
