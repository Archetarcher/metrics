package memory

// Config is a struct for in memory storage, keeps configurations
type Config struct {
	Active          bool
	FileStoragePath string
	StoreInterval   int
	Restore         bool
}
