package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/pkg/htmlbuilder"
)

type Book struct {
	Id              string        `json:"id"`
	Title           string        `json:"title"`
	Subtitle        string        `json:"subtitle,omitempty"`
	Author          []Person      `json:"author"`
	Summary         string        `json:"summary,omitempty"`
	ReleaseDateTime time.Time     `json:"releaseDate"`
	Genre           []string      `json:"genre,omitempty"`
	Series          *Series       `json:"series,omitempty"`
	File            AudiobookFile `json:"file"`
}

func (b *Book) ReleaseDate() string {
	return b.ReleaseDateTime.Format("2006-01-02")
}

func (b *Book) AuthorString() string {
	if len(b.Author) == 1 {
		return b.Author[0].String()
	} else {
		andAuthor := b.Author[len(b.Author)-1]
		commaSeparatedAuthorStrings := []string{}
		for _, author := range b.Author[:len(b.Author)-1] {
			commaSeparatedAuthorStrings = append(commaSeparatedAuthorStrings, author.String())
		}
		return fmt.Sprintf("%s and %s", strings.Join(commaSeparatedAuthorStrings[:], ", "), andAuthor.String())
	}
}

func (b *Book) Description() string {
	if b.Series == nil {
		return fmt.Sprintf("%s by %s", b.Title, b.AuthorString())
	} else {
		return fmt.Sprintf("%s by %s, %s", b.Title, b.AuthorString(), b.Series.String())
	}
}

func (b *Book) FullHTMLDescription() string {
	var sb strings.Builder
	sb.WriteString(htmlbuilder.H1(b.Title))
	if b.Subtitle != "" {
		sb.Write([]byte(htmlbuilder.H4(b.Subtitle)))
	}
	sb.WriteString(htmlbuilder.H2(fmt.Sprintf("By %s", b.AuthorString())))
	if b.Series != nil {
		sb.WriteString(htmlbuilder.H4(b.Series.String()))
	}
	if b.Summary != "" {
		sb.WriteString(htmlbuilder.P(b.Summary))
	}
	return sb.String()
}

type AudiobookFile struct {
	FileExtension string        `json:"extension"`
	FileLocation  string        `json:"location"`
	FileHash      []byte        `json:"hash"`
	FileDuration  time.Duration `json:"duration"`
	FileSize      int64         `json:"size"`
}
