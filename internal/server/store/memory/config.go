package memory

// Config is a struct for in memory storage, keeps configurations
type Config struct {
	FileStoragePath string
	Restore         bool
	Active          bool
	StoreInterval   int
}
