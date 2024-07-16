package pgx

type Config struct {
	Active          bool
	FileStoragePath string
	StoreInterval   int
	Restore         bool
	DatabaseDsn     string
}
