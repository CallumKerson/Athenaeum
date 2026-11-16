// Package feed renders a list of audiobooks as an iTunes-compatible RSS 2.0
// podcast feed.
package feed

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/CallumKerson/podcasts/v2"

	"github.com/CallumKerson/Athenaeum/pkg/audiobooks"
)

// mediaPath is the URL prefix the media root is served under. It must keep
// matching the web server's configuration: it forms part of every item's GUID.
const mediaPath = "/media"

// artworkPath is where the site serves the channel artwork, relative to the host.
const artworkPath = "/static/itunes_image.jpg"

var unixEpoch = time.Unix(0, 0).UTC()

// ErrInvalidHost reports a host that cannot prefix the enclosure URLs.
var ErrInvalidHost = errors.New("host must be an absolute http or https URL, such as https://books.example.com")

// Channel is the channel-level iTunes metadata shared by every feed.
type Channel struct {
	Explicit bool
	Language string
	Author   string
	Email    string
	// HandlePreUnixEpoch clamps release dates before 1970 up to the epoch.
	// Podcast clients vary in how they treat negative timestamps.
	HandlePreUnixEpoch bool
	// Artwork links the channel to the artwork the site serves under the host.
	// It is optional because the server's feeds, which the parity goldens
	// capture, were rendered without it.
	Artwork bool
}

// Renderer renders feeds for one host. Its fields are private so that every
// published URL — and so every GUID — comes from the one host New validated.
type Renderer struct {
	base    url.URL
	channel Channel
}

// New validates host and returns a Renderer that publishes every URL under it.
//
// A default or relative host would publish GUIDs that subscribers can never
// fetch, so it is refused, as is one with userinfo, a query or a fragment,
// which would otherwise be baked into every GUID. The trailing slash is
// trimmed from the parsed path so that a host configured with or without one
// produces the same URLs, and therefore the same GUIDs.
func New(host string, channel Channel) (*Renderer, error) {
	base, err := url.Parse(host)
	if err != nil || (base.Scheme != "http" && base.Scheme != "https") || base.Host == "" ||
		base.User != nil || base.RawQuery != "" || base.ForceQuery || base.Fragment != "" {
		return nil, fmt.Errorf("%w: got %q", ErrInvalidHost, host)
	}
	base.Path = strings.TrimRight(base.Path, "/")
	// Every published URL rewrites the path, so the host's own escaping cannot
	// be kept for some URLs and not others; dropping it here keeps them alike.
	base.RawPath = ""
	return &Renderer{base: *base, channel: channel}, nil
}

// Host is the validated host with its trailing slash trimmed, which is the
// prefix every published URL starts with.
func (r *Renderer) Host() string {
	return r.base.String()
}

// Render produces the complete feed document for a set of books.
//
// Output is fully deterministic — there is no build timestamp or generator
// stamp anywhere in it — so identical input always yields identical bytes.
func (r *Renderer) Render(books []audiobooks.Audiobook, title, description string) ([]byte, error) {
	pod := &podcasts.Podcast{
		Title:       title,
		Description: description,
		Language:    r.channel.Language,
		Link:        r.Host(),
	}

	for index := range books {
		pod.AddItem(r.item(&books[index]))
	}

	feed, err := pod.Feed(podcasts.Block)
	if err != nil {
		return nil, err
	}

	if r.channel.Author != "" {
		author := r.channel.Author
		if err := feed.SetOptions(podcasts.Author(author), podcasts.Owner(author, r.channel.Email)); err != nil {
			return nil, err
		}
	}
	if r.channel.Explicit {
		if err := feed.SetOptions(podcasts.Explicit); err != nil {
			return nil, err
		}
	}
	if r.channel.Artwork {
		if err := feed.SetOptions(podcasts.Image(r.publishedURL(artworkPath))); err != nil {
			return nil, err
		}
	}

	var buf bytes.Buffer
	if err := feed.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (r *Renderer) item(book *audiobooks.Audiobook) *podcasts.Item {
	hostedFile := r.mediaURL(book.Path)

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
		item.Image = &podcasts.ItunesImage{Href: r.mediaURL(book.ImagePath)}
	}

	return item
}

// publishedURL builds the public URL of a path under the host. Every URL the
// feed publishes goes through it, so they all share the one prefix New parsed.
//
// The path is appended to url.URL.Path and escaped by String(), rather than
// being pasted into a string and parsed back out. Parsing a concatenation
// hands the filename to the URL grammar, where three characters do damage:
//
//	?  starts the query string   — the library really does contain "What If?2"
//	#  starts the fragment       — everything after it silently disappears
//	%  begins a percent-escape   — "%20" decodes to a space, and a bare "%"
//	                               fails to parse, which fails the whole build
//
// Keeping the filename in Path keeps it as data, so each of those is escaped
// into the URL rather than interpreted as part of it.
func (r *Renderer) publishedURL(path string) string {
	published := r.base
	published.Path += path
	return published.String()
}

// mediaURL builds the public URL of a file inside the media root, which is
// also the item's GUID.
func (r *Renderer) mediaURL(filePath string) string {
	return r.publishedURL(mediaPath + filePath)
}

// pubDate offsets release dates by eight hours so that a book released on a
// given day does not appear on the previous day in western timezones.
func (r *Renderer) pubDate(book *audiobooks.Audiobook) time.Time {
	var pubDate time.Time
	if book.ReleaseDate != nil {
		pubDate = book.ReleaseDate.AsTime(time.UTC)
	}
	if r.channel.HandlePreUnixEpoch && unixEpoch.After(pubDate) {
		pubDate = unixEpoch
	}
	return pubDate.Add(8 * time.Hour)
}
