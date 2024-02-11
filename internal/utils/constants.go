package utils

import (
	"regexp"
)

var (
	TITLE         = regexp.MustCompile("__TITLE__")
	DESCRIPTION   = regexp.MustCompile("__DESCRIPTION__")
	SEPARATOR     = regexp.MustCompile(`---(.*)`)
	hashRe        = regexp.MustCompile(`(\s|\?|:|&|=|%|"|'|\/|@|\\)`)
	lessThanRe    = regexp.MustCompile("<")
	greaterThanRe = regexp.MustCompile(">")
)

const (
	BODY          = "__BODY__"
	HEADER        = "__HEADER__"
	INDEX         = "__INDEX__"
	CSS           = "__CSS__"
	URL           = "__URL__"
	IMAGE_DIR     = "images"
	CSS_FILE_NAME = "app.css"
)
