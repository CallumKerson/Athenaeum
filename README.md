# Athenaeum

[![GitHub release](https://img.shields.io/github/v/release/CallumKerson/Athenaeum?display_name=release&style=flat-square)](https://github.com/CallumKerson/Athenaeum/releases/latest)
![Build status](https://img.shields.io/github/actions/workflow/status/CallumKerson/Athenaeum/main.yaml?style=flat-square)
[![Go Reference](https://pkg.go.dev/badge/github.com/CallumKerson/Athenaeum.svg)](https://pkg.go.dev/github.com/CallumKerson/Athenaeum)
[![Go Report Card](https://goreportcard.com/badge/github.com/CallumKerson/Athenaeum?style=flat-square)](https://goreportcard.com/report/github.com/CallumKerson/Athenaeum)

Turns a collection of `.m4b` audiobooks into podcast feeds you can subscribe to.

![Athenaeum](docs/athenaeum.jpg)

## Basics

`athenaeum build` scans a library of `.m4b` audiobooks and writes a static site of podcast feeds.
There is no server and no database — serve the output directory with any web server and re-run the command whenever the library changes.

A config file at `~/.config/athenaeum/athenaeum.toml` tells it where the library
is, where to write the site, and the public URL the site is served from:

```toml
Host = "https://athenaeum.testserver.net"

[Media]
Root = "~/audiobooks"

[Site]
Root = "~/Sites/athenaeum"
```

`Host` is used to build the enclosure URLs, so it must match the address the site is actually served from.
Every setting can also be given as a flag; run `athenaeum build --help` for the full list.

That config produces a feed at `https://athenaeum.testserver.net/podcast/feed.rss`, which can be added to your favourite podcast player.
I use [Overcast](https://overcast.fm).

Alongside the main feed, the build writes one feed per author, narrator, genre and
tag, under `/podcast/authors/`, `/podcast/narrators/`, `/podcast/genre/` and
`/podcast/tags/`.

### Browsing the library

The build also writes a plain HTML site for finding your way around those feeds:

- `/` links to each way of browsing, and to the main feed
- `/books/` lists every book, with a search box that filters as you type — `?q=` prefills it, so a search can be bookmarked
- `/authors/`, `/narrators/`, `/genre/` and `/tags/` list every name, with the feed URL alongside
- `/authors/<name>/` and friends show that name's books and the URL to subscribe to

Links between pages are relative, so the site works wherever it is mounted, and
every page is readable without JavaScript — the search box and the copy-URL
buttons appear only when there is JavaScript to run them.

### Serving the site

The site references the audiobook files under `/media/`, so the web server needs to serve the library at that prefix alongside the generated tree.
With [Caddy](https://caddyserver.com):

```caddyfile
athenaeum.testserver.net {
	handle_path /media/* {
		root * /home/you/audiobooks
		file_server
	}
	handle {
		root * /home/you/Sites/athenaeum
		encode zstd gzip
		file_server
	}
}
```

### Rebuilding

`athenaeum build` is safe to run repeatedly.
It only writes files whose contents have changed, which keeps modification times — and so the web server's ETags — stable, and it caches the durations it reads out of each `.m4b` so that a rebuild only has to parse newly added books.

Files left over from a previous build are removed.
As a guard against pointing `Site.Root` at the wrong directory, the build refuses to write into a non-empty directory it did not create, which it recognises by the `.athenaeum-site` marker file it leaves there.

If `ThirdParty.NotifyOvercast` is set, a build that changed something pings Overcast to re-fetch the feeds.

### Audiobook Media Layout

The layout of the audiobooks in the root media folder is quite flexible, but the following format is recommended:

```shell
$MEDIA_ROOT/Author/Audiobook/Audiobook.m4b
```

To detect a `.m4b` audiobook, a corresponding `.toml` file must exist in the same directory and with the same name as the audiobook file.
This file provides metadata for Athenaeum to use when constructing a podcast feed.
For example, if `A Wizard of Earthsea.m4b` exists in the media root, then `A Wizard of Earthsea.toml` must exist in the same directory for Athenaeum to discover it and serve it in the podcast feed.

The following an example of the format for a `.toml` metadata file:

```toml
Title = "A Wizard of Earthsea"
Authors = ["Ursula K. Le Guin"]
ReleaseDate = 1968-11-01
Genres = ["Children's", "Fantasy"]
Narrators = ["Kobna Holdbrook-Smith"]

[Description]
Text = '''Ged, the greatest sorcerer in all Earthsea, was called Sparrowhawk in his reckless youth.

Hungry for power and knowledge, Sparrowhawk tampered with long-held secrets and loosed a terrible shadow upon the world.

This is the tale of his testing, how he mastered the mighty words of power,
tamed an ancient dragon, and crossed death's threshold to restore the balance.'''
Format = "HTML"

[Series]
Sequence = "1"
Title = "Earthsea"
```

The only required fields are `Title`, `Authors` and `ReleaseDate`.

`Series.Sequence` is a book's place in its series.
It is usually a single number, and a decimal such as `"1.5"` covers a novella that sits between two books.
An omnibus gives the span it covers instead, written as two numbers joined by a hyphen:

```toml
[Series]
Sequence = "1-3"
Title = "Earthsea"
```

Such a book is listed as `Earthsea books 1-3` rather than `Earthsea book 1`.

## Installation

Via [Homebrew](https://brew.sh):

```shell
brew tap CallumKerson/homebrew-tap/athenaeum
brew install athenaeum
```

To upgrade:

```shell
brew update && brew upgrade athenaeum
```

## DRM-Free M4B Audiobooks

Athenaeum only works on DRM-free `.m4b` audiobooks, and should only be used for private use of personally purchased
audiobooks.

To buy audiobooks that are compatible with Athenaeum, I would recommend [https://libro.fm](https://libro.fm), which can be used internationally.
According to their documentation [http://downpour.com](http://downpour.com) also provides DRM-free `.m4b` audiobooks, though it requires a US credit card.

## Note on test audiobooks

Any `*.m4b` files found in this repo are short public domain sound clips with
metadata that pretend to be audiobooks for testing purposes only.
