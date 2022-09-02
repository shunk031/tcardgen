package hugo

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gohugoio/hugo/parser/pageparser"
	"github.com/mattn/go-runewidth"
)

const (
	fmTitle      = "title"
	fmAuthors    = "authors"
	fmCategories = "categories"
	fmTags       = "tags"

	fmDate        = "date"        // priority high
	fmLastmod     = "lastmod"     // priority middle
	fmPublishDate = "publishDate" // priority low
)

type FrontMatter struct {
	Title    string
	Authors  string
	Category string
	Tags     []string
	Date     time.Time
}

// ParseFrontMatter parses the frontmatter of the specified Hugo content.
func ParseFrontMatter(filename string) (*FrontMatter, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseFrontMatter(file)
}

func parseFrontMatter(r io.Reader) (*FrontMatter, error) {
	cfm, err := pageparser.ParseFrontMatterAndContent(r)
	if err != nil {
		return nil, err
	}

	fm := &FrontMatter{}
	if fm.Title, err = getString(&cfm, fmTitle); err != nil {
		return nil, err
	}
	if isArray := isArray(&cfm, fmAuthors); isArray {
		if fm.Authors, err = getConcatenatedStringItem(&cfm, fmAuthors); err != nil {
			return nil, err
		}
	} else {
		if fm.Authors, err = getString(&cfm, fmAuthors); err != nil {
			return nil, err
		}
	}
	if fm.Category, err = getFirstStringItem(&cfm, fmCategories); err != nil {
		return nil, err
	}
	if fm.Tags, err = getTags(&cfm, fmTags); err != nil {
		return nil, err
	}
	if fm.Date, err = getContentDate(&cfm); err != nil {
		return nil, err
	}

	return fm, nil
}

func getContentDate(cfm *pageparser.ContentFrontMatter) (time.Time, error) {
	for _, key := range []string{fmDate, fmLastmod, fmPublishDate} {
		t, err := getTime(cfm, key)
		if err != nil {
			switch err.(type) {
			case *FMNotExistError:
				continue
			}
		}
		return t, err
	}
	return time.Now(), NewFMNotExistError(
		strings.Join([]string{fmDate, fmLastmod, fmPublishDate}, ", "))
}

func getTime(cfm *pageparser.ContentFrontMatter, fmKey string) (time.Time, error) {
	v, ok := cfm.FrontMatter[fmKey]
	if !ok {
		return time.Now(), NewFMNotExistError(fmKey)
	}
	switch t := v.(type) {
	case string:
		return time.Parse(time.RFC3339, t)
	case time.Time:
		return t, nil
	default:
		return time.Now(), NewFMInvalidTypeError(fmKey, "time.Time or string", t)
	}
}

func makeFixedWidthString(str string, length int) string {
	var buffer bytes.Buffer
	l := 0
	for _, c := range str {
		cl := runewidth.RuneWidth(c)
		if l+cl > length {
			for _, c = range "..." {
				buffer.WriteRune(c)
			}
			break
		}
		buffer.WriteRune(c)
		l += cl
	}

	// for i := 0; i < length-l; i++ {
	// 	buffer.WriteRune(' ')
	// }

	return buffer.String()
}

func getString(cfm *pageparser.ContentFrontMatter, fmKey string) (string, error) {
	v, ok := cfm.FrontMatter[fmKey]
	if !ok {
		return "", NewFMNotExistError(fmKey)
	}

	switch s := v.(type) {
	case string:
		if s == "" {
			return "", NewFMNotExistError(fmKey)
		}
		s = makeFixedWidthString(s, 89)
		fmt.Println(len(s), utf8.RuneCountInString(s))

		return s, nil
	default:
		return "", NewFMInvalidTypeError(fmKey, "string", s)
	}
}

func getAllStringItems(cfm *pageparser.ContentFrontMatter, fmKey string) ([]string, error) {
	v, ok := cfm.FrontMatter[fmKey]
	if !ok {
		return nil, NewFMNotExistError(fmKey)
	}

	switch arr := v.(type) {
	case []interface{}:
		var strarr []string
		for _, item := range arr {
			switch s := item.(type) {
			case string:
				if s != "" {
					strarr = append(strarr, s)
				}
			default:
				return nil, NewFMInvalidTypeError(fmKey, "string", s)
			}
		}
		if len(strarr) < 1 {
			return nil, NewFMNotExistError(fmKey)
		}
		return strarr, nil

	default:
		return nil, NewFMInvalidTypeError(fmKey, "[]interface{}", arr)
	}
}

func getConcatenatedStringItem(cfm *pageparser.ContentFrontMatter, fmKey string) (string, error) {
	arr, err := getAllStringItems(cfm, fmKey)
	if err != nil {
		return "", err
	}
	if len(arr) > 2 {
		if len(arr) > 4 {
			arr = arr[:3]
		}
		return strings.Join(arr, ", "), nil
	} else {
		return arr[0], nil
	}
}

func getTags(cfm *pageparser.ContentFrontMatter, fmKey string) ([]string, error) {
	arr, err := getAllStringItems(cfm, fmKey)
	if err != nil {
		return nil, err
	}
	if len(arr) > 3 {
		arr = arr[:3]
		arr = append(arr, "...")
		return arr, nil
	} else {
		return arr, nil
	}
}

func getFirstStringItem(cfm *pageparser.ContentFrontMatter, fmKey string) (string, error) {
	arr, err := getAllStringItems(cfm, fmKey)
	if err != nil {
		return "", err
	}
	return arr[0], nil
}

func isArray(cfm *pageparser.ContentFrontMatter, fmKey string) bool {
	switch v := cfm.FrontMatter[fmKey]; v.(type) {
	case []interface{}:
		return true
	default:
		return false
	}
}
