package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	DateTime      = "2006-01-02 15:04"
	MAX_RSS_ITEMS = 20
	RSS_FILE_NAME = "rss.xml"
	ITEM_TEMPLATE = `
  <item>
    <title>%s</title>
    <pubDate>%s</pubDate>
    <link>%s</link>
    <guid isPermaLink="true">%s</guid>
  </item>`
	RSS_TEMPLATE = `
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <title>%s</title>
    <description>%s</description>
    <link>%s</link>
    <atom:link href="%s/rss.xml" rel="self" type="application/rss+xml"/>
    <pubDate>%s</pubDate>
    <lastBuildDate>%s</lastBuildDate>
  </channel>
  %s
</rss>`
)

// SortedLimitedArray represents a slice that maintains a limited number of elements in a sorted order.
// The slice is constrained to the specified limit and sorted based on the provided comparison function.
type SortedLimitedArray[T any] struct {
	Items     []T               // Items holds the elements of the array.
	Limit     int               // Limit specifies the maximum number of elements the array can hold.
	CompareFn func(a, b T) bool // CompareFn is the function used to compare elements in the array.
}

// NewSortedLimitedArray creates and returns a new SortedLimitedArray with the specified limit and comparison function.
func NewSortedLimitedArray[T any](items []T, limit int, compareFn func(a, b T) bool) *SortedLimitedArray[T] {
	arr := &SortedLimitedArray[T]{
		Items:     make([]T, len(items)),
		Limit:     limit,
		CompareFn: compareFn,
	}
	copy(arr.Items, items)
	arr.sortAndTruncate()
	return arr
}

// Push adds elements to the array and then sorts and truncates it to ensure it doesn't exceed the specified limit.
// The elements are added to the array, and then the array is sorted and truncated according to the CompareFn and Limit.
func (s *SortedLimitedArray[T]) Push(items ...T) {
	s.Items = append(s.Items, items...)
	s.sortAndTruncate()
}

// sortAndTruncate sorts the array based on the CompareFn and truncates it to the size specified by Limit.
// This private method ensures that the array always stays within its size limit after any modifications.
func (s *SortedLimitedArray[T]) sortAndTruncate() {
	sort.Slice(s.Items, func(i, j int) bool {
		return s.CompareFn(s.Items[i], s.Items[j])
	})
	if len(s.Items) > s.Limit {
		s.Items = s.Items[:s.Limit]
	}
}

func compareFn(p1, p2 *Page) bool {
	t1, err := time.Parse(DateTime, p1.Meta.Date)
	if err != nil {
		log.Fatalf("%+v", errors.WithStack(err))
	}
	t2, err := time.Parse(DateTime, p2.Meta.Date)
	if err != nil {
		log.Fatalf("%+v", errors.WithStack(err))
	}
	return t1.After(t2) // 降順にソート
}

func CreateRssFileTask(pages []*Page, tz, outputDir, baseURL, title, description string) func() error {
	return func() error {
		location, err := time.LoadLocation(tz)
		if err != nil {
			return errors.WithStack(err)
		}
		pages, err = Filter(pages, func(p *Page, _ int) (bool, error) {
			return p.Meta.Date != "", nil
		})
		if err != nil {
			return err
		}
		sortedPages := NewSortedLimitedArray(pages, MAX_RSS_ITEMS, compareFn)

		var sb strings.Builder
		for _, page := range sortedPages.Items {
			pubDate, err := time.Parse(DateTime, page.Meta.Date)
			if err != nil {
				return errors.WithStack(err)
			}
			sb.WriteString(fmt.Sprintf(ITEM_TEMPLATE, page.Title, pubDate.UTC().In(location).Format(time.RFC1123), page.URL, page.URL))
		}

		items := sb.String()
		pubDate := time.Now().UTC().In(location).Format(time.RFC1123)
		rss := fmt.Sprintf(RSS_TEMPLATE, title, description, baseURL, baseURL, pubDate, pubDate, items)
		rss = strings.TrimSpace(rss)
		rssRSSFileName := filepath.Join(outputDir, RSS_FILE_NAME)
		err = os.WriteFile(rssRSSFileName, []byte(rss), 0644)
		if err != nil {
			return errors.WithStack(err)
		}
		return err
	}
}
