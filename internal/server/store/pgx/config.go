package pgx

// Config is a struct for pgx storage, keeps configurations
type Config struct {
	FileStoragePath string
	DatabaseDsn     string
	MigrationsPath  string
	Restore         bool
	Active          bool
	StoreInterval   int
}
