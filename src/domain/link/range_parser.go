package link

import (
	"errors"
	"strconv"
	"strings"
)

// ParseRange парсит диапазон из строки в формат Range.
//
// Поддерживаемые форматы:
// - "[0,4]" (react-admin query)
// - "0-4"
// - "resource=0-4" / "bytes=0-4"
// - "10" (end=10, start=0)
func ParseRange(strRange string) (*Range, error) {
	var start, end int
	var err error

	strRange = strings.TrimSpace(strRange)
	strRange = strings.ReplaceAll(strRange, " ", "")

	if strRange == "" {
		return nil, errors.New("invalid range header")
	}

	// Поддержка форматов вида: "resource=0-10" или "bytes=0-10"
	if i := strings.IndexByte(strRange, '='); i != -1 {
		strRange = strRange[i+1:]
	}

	// Поддержка формата: "0-10"
	if !strings.HasPrefix(strRange, "[") && strings.Contains(strRange, "-") {
		parts := strings.SplitN(strRange, "-", 2)
		if len(parts) != 2 {
			return nil, errors.New("invalid range header")
		}
		start, err = strconv.Atoi(parts[0])
		if err != nil {
			return nil, errors.New("invalid range header")
		}
		end, err = strconv.Atoi(parts[1])
		if err != nil {
			return nil, errors.New("invalid range header")
		}
		if end < start {
			return nil, errors.New("invalid range header (end < start)")
		}
		return &Range{Start: start, End: end}, nil
	}

	// Поддержка формата: "[0,10]"
	if strings.HasPrefix(strRange, "[") {
		strRange = strings.TrimPrefix(strRange, "[")
		strRange = strings.TrimSuffix(strRange, "]")

		parts := strings.Split(strRange, ",")
		if len(parts) != 2 {
			return nil, errors.New("invalid range header (len part)")
		}

		start, err = strconv.Atoi(parts[0])
		if err != nil {
			return nil, errors.New("invalid range header (atoi part one)")
		}

		end, err = strconv.Atoi(parts[1])
		if err != nil {
			return nil, errors.New("invalid range header (atoi part two)")
		}

		if end < start {
			return nil, errors.New("invalid range header (end < start)")
		}
		return &Range{Start: start, End: end}, nil
	}

	// Поддержка формата: "10" (end=10, start=0)
	end, err = strconv.Atoi(strRange)
	if err != nil {
		return nil, errors.New("invalid range header")
	}
	if end < start {
		return nil, errors.New("invalid range header (end < start)")
	}
	return &Range{Start: start, End: end}, nil
}

