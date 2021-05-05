package server

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/yura-under-review/port-api/models"
	"io"
	"strings"
)

type FileParser struct {
	r *bufio.Reader
}

func NewFileParser(r io.Reader) *FileParser {
	return &FileParser{
		r: bufio.NewReader(r),
	}
}

func (p *FileParser) Read() (*models.PortInfo, error) {

	// looking for for the opening quotes
	if err := seekRune(p.r, '"'); err != nil {
		return nil, err
	}

	// reading port symbol

	var symbolBuilder strings.Builder

	for {
		r, _, err := p.r.ReadRune()

		if errors.Is(err, io.EOF) {
			return nil, err
		}

		if err != nil {
			return nil, fmt.Errorf("failed to parse port symbol: %w", err)
		}

		if (r == '{') || (r == '}') || (r == '[') || (r == ']') {
			return nil, fmt.Errorf("failed to find port symbol end: %w", err)
		}

		if r == '"' {
			break
		}

		symbolBuilder.WriteRune(r)
	}

	portInfo := models.PortInfo{
		Symbol: symbolBuilder.String(),
	}

	// seeking to the object beginning
	if err := seekRune(p.r, '{'); err != nil {
		return nil, err
	}

	// reading object
	NBrackets := 1
	var objectBuilder strings.Builder
	objectBuilder.WriteRune('{')

	for {
		r, _, err := p.r.ReadRune()

		if errors.Is(err, io.EOF) {
			return nil, err
		}

		if err != nil {
			return nil, fmt.Errorf("failed to parse port symbol: %w", err)
		}

		switch r {
		case '{':
			NBrackets++
		case '}':
			NBrackets--
		}

		objectBuilder.WriteRune(r)

		if NBrackets == 0 {
			break
		}
	}

	if err := json.Unmarshal([]byte(objectBuilder.String()), &portInfo); err != nil {
		return nil, fmt.Errorf("failed to parse object: %w", err)
	}

	return &portInfo, nil
}

func seekRune(reader *bufio.Reader, r rune) error {
	for {
		s, _, err := reader.ReadRune()
		if err != nil {
			return io.EOF
		}

		if s == r {
			break
		}
	}

	return nil
}
