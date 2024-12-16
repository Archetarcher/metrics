package config

import (
	"flag"
)

const (
	flagRunAddrName                = "a"
	flagGRPCRunAddrName            = "ga"
	flagLogLevelName               = "l"
	flagKeyName                    = "k"
	flagFileStoragePathName        = "f"
	flagStoreIntervalName          = "i"
	flagRestoreName                = "r"
	flagEnableGRPCName             = "eg"
	flagDatabaseDsnName            = "d"
	flagDatabaseMigrationsPathName = "m"
	flagPrivateKeyPathName         = "crypto-key"
	flagConfigPathName             = "c config"
	flagTrustedSubnetName          = "t"
)

func (c *AppConfig) initFlags() {
	c.mux.Lock()
	defer c.mux.Unlock()

	flag.StringVar(&c.RunAddr, flagRunAddrName, ":8080", "address and port to run server")
	flag.StringVar(&c.LogLevel, flagLogLevelName, "info", "log level")
	flag.StringVar(&c.Key, flagKeyName, "", "key")

	flag.StringVar(&c.FileStoragePath, flagFileStoragePathName, "/tmp/metrics-pgx.json", "file storage path")
	flag.IntVar(&c.StoreInterval, flagStoreIntervalName, 300, "seconds to save data to file")
	flag.BoolVar(&c.Restore, flagRestoreName, true, "load data from file")

	flag.StringVar(&c.DatabaseDsn, flagDatabaseDsnName, "", "dsn")
	flag.StringVar(&c.MigrationsPath, flagDatabaseMigrationsPathName, "internal/server/migrations", "migrations")

	flag.StringVar(&c.PrivateKeyPath, flagPrivateKeyPathName, "private.pem", "crypto-key")
	flag.StringVar(&c.TrustedSubnet, flagTrustedSubnetName, "", "allowed ip address")

	flag.StringVar(&c.ConfigPath, flagConfigPathName, "server-config.json", "config file")

	flag.BoolVar(&c.EnableGRPC, flagEnableGRPCName, true, "run grpc server or not")
	flag.StringVar(&c.GRPCRunAddr, flagGRPCRunAddrName, ":3200", "address and port to run grpc server")

}

func (c *AppConfig) parseFlags() {
	flag.Parse()
}
