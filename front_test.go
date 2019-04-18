package front_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/deliveroo/assert-go"
	"github.com/jeffreylo/front"
)

type got struct {
	Language string   `yaml:"language"`
	Versions []string `yaml:"go"`
}

func TestFront(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/test.md")
	if err != nil {
		t.Fatal(err)
	}
	var v got
	body, err := front.Unmarshal(bytes.NewReader(data), &v)
	assert.Must(t, err)
	assert.Equal(t, v.Language, "go")
	assert.Equal(t, v.Versions, []string{"foo", "bar"})
	assert.Equal(t, body, []byte("hello"))
}
