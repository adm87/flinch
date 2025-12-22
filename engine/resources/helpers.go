package resources

// noCopy may be embedded into structs which must not be copied
type noCopy struct{}

// Lock is a no-op method to prevent copying of structs embedding noCopy.
func (*noCopy) Lock() {}

// Unlock is a no-op method to prevent copying of structs embedding noCopy.
func (*noCopy) Unlock() {}

func trimAssetPathRoot(path string) string {
	if len(path) > 0 && (path[0] == '/' || path[0] == '\\') {
		return path[1:]
	}
	return path
}
