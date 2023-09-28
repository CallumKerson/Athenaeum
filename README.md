# Athenaeum

[![GitHub release](https://img.shields.io/github/v/release/CallumKerson/Athenaeum?display_name=release&style=flat-square)](https://github.com/CallumKerson/Athenaeum/releases/latest)
![Build status](https://img.shields.io/github/actions/workflow/status/CallumKerson/Athenaeum/main.yaml?style=flat-square)
[![Go Reference](https://pkg.go.dev/badge/github.com/CallumKerson/Athenaeum.svg)](https://pkg.go.dev/github.com/CallumKerson/Athenaeum)
[![Go Report Card](https://goreportcard.com/badge/github.com/CallumKerson/Athenaeum?style=flat-square)](https://goreportcard.com/report/github.com/CallumKerson/Athenaeum)

An audiobook server that provides a podcast feed.

![Athenaeum](docs/athenaeum.jpg)

## Basics

This server will create a podcast feed from a collection of `.m4b` audiobooks and server that podcast feed on a selected
port. A config file (by default located at `~/.athenaeum/config.yaml`) can be used to customize the server.

For example, a minimal config file would look like:

```yaml
Host: https://athenaeum.testserver.net
Media:
  Root: ~/audiobooks
```

This tells the server where it is hosted, and where the root for the `.m4b` audiobooks is. To set up the host, a reverse
proxy is recommended. [Nginx](https://www.nginx.com) is standard, but for simple home use I would recommend [Caddy](https://caddyserver.com).

The above config file would produce a podcast feed at `https://athenaeum.testserver.net/podcast/feed.rss`, which can
then be added to your favourite podcast player. I use [Overcast](https://overcast.fm).

### Audiobook Media Layout

The layout of the audiobooks in the root media folder is quite flexible, but the following format is recommended:

```shell
$MEDIA_ROOT/Author/Audiobook/Audiobook.m4b
```

To detect a `.m4b` audiobook, a corresponding `.toml` file must exist in the same directory and with the same name as the
audiobook file. This file provides metadata for Athenaeum to use when constructing a podcast feed.
For example, if `A Wizard of Earthsea.m4b` exists in the media root, then `A Wizard of Earthsea.toml` must exist in the
same directory for Athenaeum to discover it and serve it in the podcast feed.

The following an example of the format for a `.toml` metadata file:

```toml
Title = "A Wizard of Earthsea"
Authors = ["Ursula K. Le Guin"]
ReleaseDate = 1968-11-01
Genres = ["Children's", "Fantasy"]
Narrators = ["Kobna Holdbrook-Smith"]

[Description]
Text = "<p>Ged, the greatest sorcerer in all Earthsea, was called Sparrowhawk in his reckless youth.</p><p>Hungry for power and knowledge, Sparrowhawk tampered with long-held secrets and loosed a terrible shadow upon the world. This is the tale of his testing, how he mastered the mighty words of power, tamed an ancient dragon, and crossed death's threshold to restore the balance.</p>"
Format = "HTML"

[Series]
Sequence = "1"
Title = "Earthsea"
```

The only required fields are `Title`, `Authors` and `ReleaseDate`.

## Installation

Via [Homebrew](https://brew.sh):

```shell
brew tap CallumKerson/homebrew-tap/athenaeum
brew install athenaeum
brew services start CallumKerson/homebrew-tap/athenaeum
```

To upgrade:

```shell
brew update && brew upgrade athenaeum
```

## DRM-Free M4B Audiobooks

Athenaeum only works on DRM-free `.m4b` audiobooks, and should only be used for private use of personally purchased audiobooks.

To buy audiobooks that are compatible with Athenaeum, I would recommend [https://libro.fm](https://libro.fm), which can
be used internationally. According to their documentation [http://downpour.com](http://downpour.com) also provides DRM-free
`.m4b` audiobooks, though it requires a US credit card.

## Note on test audiobooks

Any `*.m4b` files found in this repo are short public domain sound clips with metadata that pretend to be audiobooks for testing purposes only.
