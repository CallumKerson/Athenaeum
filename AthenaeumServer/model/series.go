package model

import (
	"fmt"
	"strconv"
)

type Series struct {
	Entry float32 `json:"entry"`
	Title string  `json:"title"`
}

func (s *Series) String() string {
	return fmt.Sprintf("%s Book %s", s.Title, s.EntryString())
}

func (s *Series) EntryString() string {
	return strconv.FormatFloat(float64(s.Entry), 'f', -1, 64)
}
