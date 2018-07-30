package sync

import (
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"

	"golang.org/x/text/unicode/norm"

	"github.com/golang/protobuf/ptypes"

	"github.com/havoc-io/mutagen/pkg/filesystem"
)

const (
	// scannerCopyBufferSize specifies the size of the internal buffer that a
	// scanner uses to copy file data.
	// TODO: Figure out if we should set this on a per-machine basis. This value
	// is taken from Go's io.Copy method, which defaults to allocating a 32k
	// buffer if none is provided.
	scannerCopyBufferSize = 32 * 1024

	// defaultInitialCacheCapacity specifies the default capacity for new
	// filesystem and ignore caches when the corresponding existing cache is nil
	// or empty. It is designed to save several rounds of cache capacity
	// doubling on insert without always allocating a huge cache. Its value is
	// somewhat arbitrary.
	defaultInitialCacheCapacity = 1024
)

type scanner struct {
	root                   string
	hasher                 hash.Hash
	cache                  *Cache
	ignorer                *ignorer
	ignoreCache            map[string]bool
	newCache               *Cache
	newIgnoreCache         map[string]bool
	buffer                 []byte
	deviceID               uint64
	recomposeUnicode       bool
	preservesExecutability bool
}

func (s *scanner) file(path string, info os.FileInfo) (*Entry, error) {
	// Extract metadata.
	mode := info.Mode()
	modificationTime := info.ModTime()
	size := uint64(info.Size())

	// Compute executability.
	executable := s.preservesExecutability && anyExecutableBitSet(mode)

	// Convert the timestamp to Protocol Buffers format.
	modificationTimeProto, err := ptypes.TimestampProto(modificationTime)
	if err != nil {
		return nil, errors.Wrap(err, "unable to convert modification time format")
	}

	// Try to find a cached digest. We only enforce that type, modification
	// time, and size haven't changed in order to re-use digests.
	var digest []byte
	cached, hit := s.cache.Entries[path]
	match := hit &&
		(os.FileMode(cached.Mode)&os.ModeType) == (mode&os.ModeType) &&
		cached.ModificationTime != nil &&
		modificationTimeProto.Seconds == cached.ModificationTime.Seconds &&
		modificationTimeProto.Nanos == cached.ModificationTime.Nanos &&
		cached.Size == size
	if match {
		digest = cached.Digest
	}

	// If we weren't able to pull a digest from the cache, compute one manually.
	if digest == nil {
		// Open the file and ensure its closure.
		file, err := os.Open(filepath.Join(s.root, path))
		if err != nil {
			return nil, errors.Wrap(err, "unable to open file")
		}
		defer file.Close()

		// Reset the hash state.
		s.hasher.Reset()

		// Copy data into the hash and very that we copied as much as expected.
		if copied, err := io.CopyBuffer(s.hasher, file, s.buffer); err != nil {
			return nil, errors.Wrap(err, "unable to hash file contents")
		} else if uint64(copied) != size {
			return nil, errors.New("hashed size mismatch")
		}

		// Compute the digest.
		digest = s.hasher.Sum(nil)
	}

	// Add a cache entry.
	s.newCache.Entries[path] = &CacheEntry{
		Mode:             uint32(mode),
		ModificationTime: modificationTimeProto,
		Size:             size,
		Digest:           digest,
	}

	// Success.
	return &Entry{
		Kind:       EntryKind_File,
		Executable: executable,
		Digest:     digest,
	}, nil
}

func (s *scanner) symlink(path string, enforcePortable bool) (*Entry, error) {
	// Read the link target.
	target, err := os.Readlink(filepath.Join(s.root, path))
	if err != nil {
		return nil, errors.Wrap(err, "unable to read symlink target")
	}

	// If requested, enforce that the link is portable, otherwise just ensure
	// that it's non-empty (this is required even in POSIX raw mode).
	if enforcePortable {
		target, err = normalizeSymlinkAndEnsurePortable(path, target)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("invalid symlink (%s)", path))
		}
	} else if target == "" {
		return nil, errors.New("symlink target is empty")
	}

	// Success.
	return &Entry{
		Kind:   EntryKind_Symlink,
		Target: target,
	}, nil
}

func (s *scanner) directory(path string, info os.FileInfo, symlinkMode SymlinkMode) (*Entry, error) {
	// Verify that we haven't crossed a directory boundary (which might
	// potentially change executability preservation or Unicode decomposition
	// behavior).
	if id, err := filesystem.DeviceID(info); err != nil {
		return nil, errors.Wrap(err, "unable to extract directory device ID")
	} else if id != s.deviceID {
		return nil, errors.New("scan crossed filesystem boundary")
	}

	// Read directory contents.
	directoryContents, err := filesystem.DirectoryContents(filepath.Join(s.root, path))
	if err != nil {
		return nil, errors.Wrap(err, "unable to read directory contents")
	}

	// Compute entries.
	contents := make(map[string]*Entry, len(directoryContents))
	for _, c := range directoryContents {
		// Extract the content name.
		name := c.Name()

		// Determine whether or not this is a test file, and if so skip it. It's
		// possible to see other sessions' test files when session roots
		// overlap.
		if filesystem.IsExecutabilityProbeFileName(name) {
			continue
		} else if filesystem.IsUnicodeProbeFileName(name) {
			continue
		}

		// Recompose Unicode in the content name if necessary.
		if s.recomposeUnicode {
			name = norm.NFC.String(name)
		}

		// Compute the content path.
		contentPath := pathJoin(path, name)

		// Compute the kind for this content, skipping if unsupported.
		kind := EntryKind_File
		if mode := c.Mode(); mode&os.ModeDir != 0 {
			kind = EntryKind_Directory
		} else if mode&os.ModeSymlink != 0 {
			kind = EntryKind_Symlink
		} else if mode&os.ModeType != 0 {
			continue
		}

		// Determine whether or not this path is ignored and update the new
		// ignore cache.
		ignored, ok := s.ignoreCache[contentPath]
		if !ok {
			ignored = s.ignorer.ignored(contentPath, kind == EntryKind_Directory)
		}
		s.newIgnoreCache[contentPath] = ignored
		if ignored {
			continue
		}

		// Handle based on kind.
		var entry *Entry
		if kind == EntryKind_File {
			entry, err = s.file(contentPath, c)
		} else if kind == EntryKind_Symlink {
			if symlinkMode == SymlinkMode_SymlinkPortable {
				entry, err = s.symlink(contentPath, true)
			} else if symlinkMode == SymlinkMode_SymlinkIgnore {
				continue
			} else if symlinkMode == SymlinkMode_SymlinkPOSIXRaw {
				entry, err = s.symlink(contentPath, false)
			} else {
				panic("unsupported symlink mode")
			}
		} else if kind == EntryKind_Directory {
			entry, err = s.directory(contentPath, c, symlinkMode)
		} else {
			panic("unhandled entry kind")
		}

		// Watch for errors and add the entry.
		if err != nil {
			return nil, err
		}

		// Add the content.
		contents[name] = entry
	}

	// Success.
	return &Entry{
		Kind:     EntryKind_Directory,
		Contents: contents,
	}, nil
}

// TODO: Note that the provided cache is assumed to be valid (i.e. that it
// doesn't have any nil entries), so callers should run EnsureValid on anything
// they pull from disk
func Scan(root string, hasher hash.Hash, cache *Cache, ignores []string, ignoreCache map[string]bool, symlinkMode SymlinkMode) (*Entry, bool, bool, *Cache, map[string]bool, error) {
	// A nil cache is technically valid, but if the provided cache is nil,
	// replace it with an empty one, that way we don't have to use the
	// GetEntries accessor everywhere.
	if cache == nil {
		cache = &Cache{}
	}

	// Create the ignorer.
	ignorer, err := newIgnorer(ignores)
	if err != nil {
		return nil, false, false, nil, nil, errors.Wrap(err, "unable to create ignorer")
	}

	// Verify that the symlink mode is valid for this platform.
	if symlinkMode == SymlinkMode_SymlinkPOSIXRaw && runtime.GOOS == "windows" {
		return nil, false, false, nil, nil, errors.New("raw POSIX symlinks not supported on Windows")
	}

	// Create a new cache to populate. Estimate its capacity based on the
	// existing cache length. If the existing cache is empty, create one with
	// the default capacity.
	initialCacheCapacity := defaultInitialCacheCapacity
	if cacheLength := len(cache.Entries); cacheLength != 0 {
		initialCacheCapacity = cacheLength
	}
	newCache := &Cache{
		Entries: make(map[string]*CacheEntry, initialCacheCapacity),
	}

	// Create a new ignore cache to populate. Estimate its capacity based on the
	// existing ignore cache length. If the existing cache is empty, create one
	// with the default capacity.
	initialIgnoreCacheCapacity := defaultInitialCacheCapacity
	if ignoreCacheLength := len(ignoreCache); ignoreCacheLength != 0 {
		initialIgnoreCacheCapacity = ignoreCacheLength
	}
	newIgnoreCache := make(map[string]bool, initialIgnoreCacheCapacity)

	// Create a scanner.
	s := &scanner{
		root:           root,
		hasher:         hasher,
		cache:          cache,
		ignorer:        ignorer,
		ignoreCache:    ignoreCache,
		newCache:       newCache,
		newIgnoreCache: newIgnoreCache,
		buffer:         make([]byte, scannerCopyBufferSize),
	}

	// Create the snapshot.
	if info, err := os.Lstat(root); err != nil {
		if os.IsNotExist(err) {
			return nil, false, false, newCache, newIgnoreCache, nil
		} else {
			return nil, false, false, nil, nil, errors.Wrap(err, "unable to probe scan root")
		}
	} else if mode := info.Mode(); mode&os.ModeDir != 0 {
		// Grab and set the device ID for the root directory.
		if id, err := filesystem.DeviceID(info); err != nil {
			return nil, false, false, nil, nil, errors.Wrap(err, "unable to probe root device ID")
		} else {
			s.deviceID = id
		}

		// Probe and set Unicode decomposition behavior for the root directory.
		if decomposes, err := filesystem.DecomposesUnicode(root); err != nil {
			return nil, false, false, nil, nil, errors.Wrap(err, "unable to probe root Unicode decomposition behavior")
		} else {
			s.recomposeUnicode = decomposes
		}

		// Probe and set executability preservation behavior for the root directory.
		if preserves, err := filesystem.PreservesExecutability(root); err != nil {
			return nil, false, false, nil, nil, errors.Wrap(err, "unable to probe root executability preservation behavior")
		} else {
			s.preservesExecutability = preserves
		}

		// Perform a recursive scan.
		if rootEntry, err := s.directory("", info, symlinkMode); err != nil {
			return nil, false, false, nil, nil, err
		} else {
			return rootEntry, s.preservesExecutability, s.recomposeUnicode, newCache, newIgnoreCache, nil
		}
	} else if mode&os.ModeType != 0 {
		// We disallow symlinks as synchronization roots because there's no easy
		// way to propagate changes to them.
		return nil, false, false, nil, nil, errors.New("invalid scan root type")
	} else {
		// Probe and set executability preservation behavior for the parent of the root directory.
		if preserves, err := filesystem.PreservesExecutability(filepath.Dir(root)); err != nil {
			return nil, false, false, nil, nil, errors.Wrap(err, "unable to probe root parent executability preservation behavior")
		} else {
			s.preservesExecutability = preserves
		}

		// Perform a scan of the root file.
		if rootEntry, err := s.file("", info); err != nil {
			return nil, false, false, nil, nil, err
		} else {
			return rootEntry, s.preservesExecutability, false, newCache, newIgnoreCache, nil
		}
	}
}
