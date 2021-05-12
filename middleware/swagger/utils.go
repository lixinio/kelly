package swagger

import (
	"fmt"
	"regexp"
	"strings"
)

func parseFileNode(entry string) (filepath string, node string, err error) {
	entrys := strings.Split(entry, ":")
	if len(entrys) != 2 {
		err = fmt.Errorf("invalid swagger entry (%s)", entry)
		return
	}

	filepath = entrys[0]
	node = entrys[1]
	if len(filepath) == 0 {
		err = fmt.Errorf("invalid swagger entry (%s),invalid filepath", entry)
		return
	}
	if len(node) == 0 {
		err = fmt.Errorf("invalid swagger entry (%s),invalid node", entry)
		return
	}
	return
}

type pathEditor struct {
	regex  string
	engine *regexp.Regexp
}

func newPathEditor() *pathEditor {
	regex := `:([0-9a-zA-Z]+)`
	re := regexp.MustCompile(regex)
	return &pathEditor{
		regex:  regex,
		engine: re,
	}
}

func (pe *pathEditor) update(path string) string {
	return pe.engine.ReplaceAllString(path, "{$1}")
}

type tagOptions string

func (o tagOptions) contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}

	s := string(o)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if s == optionName {
			return true
		}
		s = next
	}
	return false
}

func (o tagOptions) getValue(optionName string) (string, bool) {
	if len(o) == 0 {
		return "", false
	}

	s := string(o)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if strings.ContainsAny(s, "=") {
			var strs = strings.Split(s, "=")
			if len(strs) == 2 && strs[0] == optionName {
				return strs[1], true
			}
		}
		s = next
	}
	return "", false
}
