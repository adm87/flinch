package resources

import (
	"fmt"
	"io"
	"io/fs"
	"sync"
	"sync/atomic"
)

var (
	batchID = atomic.Uint64{}

	assetLocks = &sync.Pool{
		New: func() any {
			return &AssetLock{}
		},
	}
)

// ============================== Assets ==============================

// Asset is a unique identifier for a resource within an AssetManifest.
type Asset uint64

// AssetManifest maps Asset identifiers to their corresponding file paths.
type AssetManifest map[Asset]string

// AssetLock is a guard used to prevent multiple threads from operating on the same resource simultaneously.
//
// Each AssetLock must be released exactly once by calling Release(). Calling Release() multiple times
// or using an AssetLock after release will panic.
//
// AssetLocks are pooled and reused internally - do not retain references after calling Release().
type AssetLock struct {
	noCopy

	batchID uint64 // Associated batch ID that currently holds the lock
	asset   Asset  // The asset being locked

	rs      *ResourceSystem // The ResourceSystem managing this lock
	assetMu *sync.Mutex     // Mutex for the specific asset
}

// Release unlocks the AssetLock, allowing other threads to acquire the resource.
//
// This method must be called exactly once. Calling Release() multiple times will panic.
// After calling Release(), the AssetLock should not be used further.
func (al *AssetLock) Release() {
	if al.assetMu == nil {
		panic("AssetLock released multiple times. DO NOT release an AssetLock more than once.")
	}

	al.assetMu.Unlock()
	al.assetMu = nil

	al.rs.lockReleased(al)
}

// ============================== Resource System ==============================

// ResourceSystemOptions defines configuration options for a ResourceSystem.
type ResourceSystemOptions struct {
	// TrimRoot indicates whether to trim the root directory from resource paths.
	//
	// Some filesystems may require the root directory to be trimmed from resource paths
	// to correctly locate resources.
	TrimRoot bool
}

// ResourceSystem represents a collection of resources, providing utilities for loading and managing them.
//
// ResourceSystem is safe for concurrent use by multiple goroutines. It manages per-asset locking to
// ensure only one goroutine can load or operate on a specific asset at a time, while allowing concurrent
// access to different assets.
//
// Note: The ResourceSystem does not manage the lifecycle of the loaded assets themselves. It is the caller's
// responsibility to handle asset caching, unloading, and memory management as needed.
type ResourceSystem struct {
	options    ResourceSystemOptions
	manifest   AssetManifest
	filesystem fs.FS
	name       string

	locks   map[uint64]*AssetLock
	assetMu map[Asset]*sync.Mutex
	mu      sync.RWMutex
}

// NewResourceSystem creates a new ResourceSystem with the given name, manifest, and options.
func NewResourceSystem(name string, manifest AssetManifest, options ResourceSystemOptions) *ResourceSystem {
	return &ResourceSystem{
		locks:    make(map[uint64]*AssetLock),
		assetMu:  make(map[Asset]*sync.Mutex),
		name:     name,
		manifest: manifest,
		options:  options,
	}
}

// SetFileSystem sets the filesystem to be used by the ResourceSystem for loading assets.
//
// The filesystem must implement the fs.FS interface. If no filesystem is set, attempts to
// read assets will result in an error.
func (rs *ResourceSystem) SetFileSystem(fs fs.FS) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.filesystem = fs
}

// AcquireAssetLock attempts to acquire a lock for the specified asset within the resource system.
//
// If the asset is already locked by another batch, this method will block until the lock is available.
//
// The batchID should be unique per loading operation (typically a goroutine ID or operation ID).
// A batch can only hold one lock at a time - attempting to acquire multiple locks simultaneously
// will panic.
//
// The returned AssetLock must be released by calling Release() exactly once when done.
func (rs *ResourceSystem) LockAsset(batchID uint64, asset Asset) *AssetLock {
	rs.mu.Lock()
	if _, exists := rs.manifest[asset]; !exists {
		rs.mu.Unlock()
		panic(fmt.Sprintf("asset 0x%x does not exist in resource system %s", asset, rs.name))
	}

	if _, exists := rs.locks[batchID]; exists {
		rs.mu.Unlock()
		panic(fmt.Sprintf("resource batch %d attempted to acquire multiple locks simultaneously", batchID))
	}

	assetMutex, exists := rs.assetMu[asset]
	if !exists {
		assetMutex = &sync.Mutex{}
		rs.assetMu[asset] = assetMutex
	}

	lock := assetLocks.Get().(*AssetLock)
	lock.batchID = batchID
	lock.asset = asset
	lock.assetMu = assetMutex
	lock.rs = rs

	rs.locks[batchID] = lock
	rs.mu.Unlock()

	assetMutex.Lock()

	return lock
}

// ReadBytes reads the raw byte data of the specified asset from the resource system.
//
// If the asset does not exist, an error is returned.
func (rs *ResourceSystem) ReadBytes(asset Asset) ([]byte, error) {
	rs.mu.RLock()
	path := rs.manifest[asset]
	fs := rs.filesystem
	rs.mu.RUnlock()

	if fs == nil {
		return nil, fmt.Errorf("resource system %s has no associated filesystem", rs.name)
	}

	if rs.options.TrimRoot {
		path = trimAssetPathRoot(path)
	}

	file, err := fs.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open asset 0x%x (%s) in resource system %s: %w", asset, path, rs.name, err)
	}
	defer file.Close()

	return io.ReadAll(file)
}

// CreateBatch creates a new LoadingOperation batch with the specified loading tasks.
//
// The returned LoadingOperation can be executed to perform the loading tasks within the resource system.
func (rs *ResourceSystem) CreateBatch(tasks ...LoadingTask) *LoadingOperation {
	return &LoadingOperation{
		batchID: batchID.Add(1),
		rs:      rs,
		tasks:   tasks,
	}
}

// lockReleased is an internal method called when an AssetLock is released.
//
// This method removes the lock from the active locks map and returns it to the pool.
// Asset mutexes are intentionally kept in the assetMu map to prevent race conditions
// where threads may hold references to mutexes after unlocking rs.mu.
func (rs *ResourceSystem) lockReleased(lock *AssetLock) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	delete(rs.locks, lock.batchID)
	assetLocks.Put(lock)
}
