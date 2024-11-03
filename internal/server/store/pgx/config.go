package pgx

// Config is a struct for pgx storage, keeps configurations
type Config struct {
	Active          bool
	FileStoragePath string
	StoreInterval   int
	Restore         bool
	DatabaseDsn     string
	MigrationsPath  string
}
