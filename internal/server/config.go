package server

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/spf13/pflag"
)

// Backend config
type Config struct {
	Address         entities.Address  `env:"ADDRESS" json:"address"`
	GRPCAddress     entities.Address  `env:"GRPC_ADDRESS" json:"grpc_address"`
	StoreInterval   int               `env:"STORE_INTERVAL" json:"store_interval"`
	StorePath       string            `env:"FILE_STORAGE_PATH" json:"store_file"`
	RestoreOnStart  bool              `env:"RESTORE" json:"restore"`
	DatabaseDSN     string            `env:"DATABASE_DSN" json:"database_dsn"`
	Secret          entities.Secret   `env:"KEY" json:"key"`
	ProfilerAddress entities.Address  `env:"PROFILER_ADDRESS" json:"profiler_address"`
	PrivateKeyPath  entities.FilePath `env:"CRYPTO_KEY" json:"crypto_key"`
	TrustedSubnet   *net.IPNet        `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
	ConfigFilePath  entities.FilePath `env:"CONFIG"`
}

func NewConfig() (*Config, error) {
	var err error

	config := &Config{
		Address:         "0.0.0.0:8080",
		GRPCAddress:     "0.0.0.0:50051",
		StoreInterval:   300,
		RestoreOnStart:  true,
		ProfilerAddress: "0.0.0.0:8081",
	}

	err = config.parse()
	if err != nil {
		return nil, err
	}

	return config, err
}

func (c *Config) parse() error {
	var err error

	err = c.tryLoadJSONConfig()
	if err != nil {
		return err
	}

	err = c.parseFlags(os.Args[0], os.Args[1:])
	if err != nil {
		return err
	}

	err = c.parseEnv()
	if err != nil {
		return err
	}

	return err
}

func (c *Config) parseFlags(progname string, args []string) error {
	flags := pflag.NewFlagSet(progname, pflag.ContinueOnError)

	address := c.Address
	flags.VarP(&address, "address", "a", "address:port for HTTP API requests")

	secret := c.Secret
	flags.VarP(&secret, "secret", "k", "a key to sign outgoing data")

	privateKeyPath := c.PrivateKeyPath
	flags.VarP(&privateKeyPath, "crypto-key", "", "path to public key to encrypt agent -> server communications")

	configPath := entities.FilePath("") // register var for compatibility
	flags.VarP(&configPath, "config", "c", "path to configuration file in JSON format")

	defaultSubnet := net.IPNet{}
	trustedSubnet := flags.IPNetP("trusted-subnet", "t", defaultSubnet, "trusted subnet in CIDR notation")

	// define flags
	flags.IntVarP(&c.StoreInterval, "store-interval", "i", c.StoreInterval, "interval (s) for dumping metrics to the disk, zero value means saving after each request")
	flags.StringVarP(&c.StorePath, "store-file", "f", c.StorePath, "path to file to store metrics")
	flags.BoolVarP(&c.RestoreOnStart, "restore", "r", c.RestoreOnStart, "whether to restore state on startup")
	flags.StringVarP(&c.DatabaseDSN, "database", "d", c.DatabaseDSN, "PostgreSQL database DSN")

	pErr := flags.Parse(args)
	if pErr != nil {
		return pErr
	}

	// fill values
	flags.Visit(func(f *pflag.Flag) {
		switch f.Name {
		case "address":
			c.Address = address
		case "secret":
			c.Secret = secret
		case "crypto-key":
			c.PrivateKeyPath = privateKeyPath
		case "trusted-subnet":
			c.TrustedSubnet = trustedSubnet
		}
	})

	return nil
}

func (c *Config) parseEnv() error {
	if err := env.Parse(c); err != nil {
		return err
	}

	return nil
}

func (c *Config) tryLoadJSONConfig() error {
	configArg := os.Getenv("CONFIG")

	// args is higher prior
	for i, arg := range os.Args {
		if (arg == "-c" || arg == "--config") && i+1 < len(os.Args) {
			configArg = os.Args[i+1]
			break
		}
	}

	if len(configArg) > 0 {
		err := c.loadConfigFromFile(entities.FilePath(configArg))
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) loadConfigFromFile(src entities.FilePath) error {
	data, err := os.ReadFile(src.String())
	if err != nil {
		return fmt.Errorf("server.loadConfigFromFile - os.ReadFile: %w", err)
	}

	if err := json.Unmarshal(data, c); err != nil {
		return fmt.Errorf("server.loadConfigFromFile - json.Unmarshal: %w", err)
	}

	return nil
}
