// Package front parses a body of text for yaml-compatible frontmatter.
package front

import (
	"bufio"
	"bytes"
	"io"

	"github.com/go-yaml/yaml"
	"github.com/pkg/errors"
)

const splitToken = "---"

// Unmarshal unmarshals input into v, returning the body value.
func Unmarshal(input io.Reader, v interface{}) ([]byte, error) {
	bufsize := 1024 * 1024
	buf := make([]byte, bufsize)

	s := bufio.NewScanner(input)
	s.Buffer(buf, bufsize)

	var frontMatter []byte
	var body []byte
	s.Split(splitFunc)
	n := 0
	for s.Scan() {
		if n == 0 {
			frontMatter = s.Bytes()
		} else if n == 1 {
			body = s.Bytes()
		}
		n++
	}
	if err := s.Err(); err != nil {
		return nil, errors.Wrap(err, "front: failed to scan text")
	}
	if err := yaml.Unmarshal(frontMatter, v); err != nil {
		return nil, errors.Wrap(err, "front: failed to parse yaml")
	}
	return body, nil
}

func splitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	delim, err := sniffDelim(data)
	if err != nil {
		return 0, nil, err
	}
	if delim != splitToken {
		return 0, nil, errors.Errorf("splitFunc: %s is not a supported delimiter", delim)
	}
	if x := bytes.Index(data, []byte(delim)); x >= 0 {
		if next := bytes.Index(data[x+len(delim):], []byte(delim)); next > 0 {
			return next + len(delim), bytes.TrimSpace(data[:next+len(delim)]), nil
		}
		return len(data), bytes.TrimSpace(data[x+len(delim):]), nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

func sniffDelim(input []byte) (string, error) {
	if len(input) < 4 {
		return "", errors.New("sniffDelim: input is empty")
	}
	return string(input[:3]), nil
}
