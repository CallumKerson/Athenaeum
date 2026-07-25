package scan

import (
	"encoding/json"
	"os"
	"time"

	"github.com/CallumKerson/Athenaeum/internal/fsutil"
)

// Entry records what we learned from opening an .m4b file. Only the duration is
// expensive to obtain: it requires parsing the file's moov box. Size and modTime
// are stored purely to decide whether the duration is still valid.
type Entry struct {
	Size     int64         `json:"size"`
	ModTime  time.Time     `json:"modTime"`
	Duration time.Duration `json:"duration"`
}

// Cache maps a media-root-relative .m4b path to its last known Entry.
//
// Deliberately narrow: the .toml metadata is never cached, because editing a
// description does not change the .m4b's mtime. Caching whole audiobooks would
// silently ignore metadata edits.
type Cache struct {
	entries map[string]Entry
	Hits    int
	Misses  int
}

func NewCache() *Cache {
	return &Cache{entries: map[string]Entry{}}
}

// LoadCache reads a cache file. A missing or unreadable cache is not an error:
// it just means every .m4b gets re-parsed this run.
func LoadCache(path string) (*Cache, error) {
	cache := NewCache()
	contents, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cache, nil
		}
		return cache, err
	}
	if err := json.Unmarshal(contents, &cache.entries); err != nil {
		return NewCache(), err
	}
	// A cache file of "null" unmarshals to a nil map without error, which would
	// panic on the first Store.
	if cache.entries == nil {
		cache.entries = map[string]Entry{}
	}
	return cache, nil
}

// Lookup returns the cached duration if the file is unchanged since it was
// recorded. Both size and modification time must match.
func (c *Cache) Lookup(relPath string, size int64, modTime time.Time) (time.Duration, bool) {
	entry, found := c.entries[relPath]
	if !found || entry.Size != size || !entry.ModTime.Equal(modTime) {
		c.Misses++
		return 0, false
	}
	c.Hits++
	return entry.Duration, true
}

func (c *Cache) Store(relPath string, size int64, modTime time.Time, duration time.Duration) {
	c.entries[relPath] = Entry{Size: size, ModTime: modTime, Duration: duration}
}

// Prune drops entries for files that were not seen during this scan, so the
// cache cannot grow without bound as books are removed from the library.
func (c *Cache) Prune(seen map[string]bool) {
	for relPath := range c.entries {
		if !seen[relPath] {
			delete(c.entries, relPath)
		}
	}
}

func (c *Cache) Save(path string) error {
	contents, err := json.Marshal(c.entries)
	if err != nil {
		return err
	}
	return fsutil.WriteAtomic(path, contents)
}
