// Package feed renders a list of audiobooks as an iTunes-compatible RSS 2.0
// podcast feed.
package feed

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/CallumKerson/podcasts"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

var unixEpoch = time.Unix(0, 0).UTC()

// Renderer holds everything about a feed that does not vary between feeds: the
// host the media is served from and the channel-level iTunes metadata.
type Renderer struct {
	// Host is the external base URL, without a trailing slash.
	Host string
	// MediaPath is the URL prefix the media root is served under, without slashes.
	MediaPath string
	ImageLink string
	Explicit  bool
	Language  string
	Author    string
	Email     string
	// HandlePreUnixEpoch clamps release dates before 1970 up to the epoch.
	// Podcast clients vary in how they treat negative timestamps.
	HandlePreUnixEpoch bool
}

// Render produces the complete feed document for a set of books.
//
// Output is fully deterministic — there is no build timestamp or generator
// stamp anywhere in it — so identical input always yields identical bytes.
func (r *Renderer) Render(books []audiobooks.Audiobook, title, description string) ([]byte, error) {
	pod := &podcasts.Podcast{
		Title:       title,
		Description: description,
		Language:    r.Language,
		Link:        r.host(),
	}

	for index := range books {
		item, err := r.item(&books[index])
		if err != nil {
			return nil, err
		}
		pod.AddItem(item)
	}

	feed, err := pod.Feed(podcasts.Block)
	if err != nil {
		return nil, err
	}

	if r.Author != "" {
		if err := feed.SetOptions(podcasts.Author(r.Author), podcasts.Owner(r.Author, r.Email)); err != nil {
			return nil, err
		}
	}
	if r.Explicit {
		if err := feed.SetOptions(podcasts.Explicit); err != nil {
			return nil, err
		}
	}
	if r.ImageLink != "" {
		if err := feed.SetOptions(podcasts.Image(r.ImageLink)); err != nil {
			return nil, err
		}
	}

	var buf bytes.Buffer
	if err := feed.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (r *Renderer) item(book *audiobooks.Audiobook) (*podcasts.Item, error) {
	hostedFile, err := r.mediaURL(book.Path)
	if err != nil {
		return nil, err
	}

	item := &podcasts.Item{
		Title: book.Title,
		Description: &podcasts.CDATAText{
			Value: fmt.Sprintf("%s by %s", book.Title, book.GetAuthor()),
		},
		PubDate:  podcasts.NewPubDate(r.pubDate(book)),
		Duration: podcasts.NewDuration(book.Duration),
		// The GUID is the enclosure URL. Changing how either is built makes every
		// subscriber treat every book as new and re-download the whole library.
		GUID: hostedFile,
		Enclosure: &podcasts.Enclosure{
			URL:    hostedFile,
			Length: fmt.Sprintf("%d", book.FileSize),
			Type:   book.MIMEType,
		},
		Subtitle:       book.GetAuthor(),
		ContentEncoded: &podcasts.CDATAText{Value: summaryHTML(book)},
	}

	if book.ImagePath != "" {
		if imageURL, imageErr := r.mediaURL(book.ImagePath); imageErr == nil {
			item.Image = &podcasts.ItunesImage{Href: imageURL}
		}
	}

	return item, nil
}

// host strips surrounding slashes so that a configured host with or without a
// trailing slash produces the same URLs — and therefore the same GUIDs.
func (r *Renderer) host() string {
	return strings.Trim(r.Host, "/")
}

// mediaURL builds the public URL of a file inside the media root.
//
// The file path is assigned to url.URL.Path and escaped by String(), rather
// than being pasted into a string and parsed back out. Parsing a concatenation
// hands the filename to the URL grammar, where three characters do damage:
//
//	?  starts the query string   — the library really does contain "What If?2"
//	#  starts the fragment       — everything after it silently disappears
//	%  begins a percent-escape   — "%20" decodes to a space, and a bare "%"
//	                               fails to parse, which fails the whole build
//
// Assigning to Path keeps the filename as data, so each of those is escaped
// into the URL rather than interpreted as part of it.
func (r *Renderer) mediaURL(filePath string) (string, error) {
	base, err := url.Parse(r.host())
	if err != nil {
		return "", err
	}
	base.Path = strings.TrimSuffix(base.Path, "/") + "/" + strings.Trim(r.MediaPath, "/") + filePath
	// Path has been rewritten, so any escaping the host carried is now stale.
	base.RawPath = ""
	return base.String(), nil
}

// pubDate offsets release dates by eight hours so that a book released on a
// given day does not appear on the previous day in western timezones.
func (r *Renderer) pubDate(book *audiobooks.Audiobook) time.Time {
	var pubDate time.Time
	if book.ReleaseDate != nil {
		pubDate = book.ReleaseDate.AsTime(time.UTC)
	}
	if r.HandlePreUnixEpoch && unixEpoch.After(pubDate) {
		pubDate = unixEpoch
	}
	return pubDate.Add(8 * time.Hour)
}
