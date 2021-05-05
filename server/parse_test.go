package server

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yura-under-review/port-api/models"
)

func Test_Parse(t *testing.T) {

	parser := NewFileParser(strings.NewReader(data))

	var items []*models.PortInfo
	for {
		item, err := parser.Read()

		if item == nil {
			break
		}

		items = append(items, item)

		assert.Truef(t, errors.Is(err, io.EOF), "bad error")

		if err != nil {
			break
		}
	}

	for _, item := range items {
		fmt.Printf("%#v\n", item)
	}

}
