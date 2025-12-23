package resources

import "strings"

// noCopy may be embedded into structs which must not be copied
type noCopy struct{}

// Lock is a no-op method to prevent copying of structs embedding noCopy.
func (*noCopy) Lock() {}

// Unlock is a no-op method to prevent copying of structs embedding noCopy.
func (*noCopy) Unlock() {}

func trimAssetPathRoot(path string) string {
	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 2 {
		return path
	}
	return parts[1]
}
