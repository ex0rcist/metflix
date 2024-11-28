package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/ex0rcist/metflix/internal/entities"
	"github.com/spf13/pflag"
)

// Agent config.
type Config struct {
	Address        entities.Address  `env:"ADDRESS" json:"address"`
	PollInterval   int               `env:"POLL_INTERVAL" json:"poll_interval"`
	ReportInterval int               `env:"REPORT_INTERVAL" json:"report_interval"`
	RateLimit      int               `env:"RATE_LIMIT" json:"-"`
	Secret         entities.Secret   `env:"KEY" json:"key"`
	PublicKeyPath  entities.FilePath `env:"CRYPTO_KEY" json:"crypto_key"`
}

func NewConfig() (*Config, error) {
	var err error

	config := &Config{
		Address:        "0.0.0.0:8080",
		PollInterval:   2,
		ReportInterval: 10,
		RateLimit:      -1,
	}

	err = config.parse()
	if err != nil {
		return nil, err
	}

	return config, err
}

// Stringer
func (c Config) String() string {
	str := []string{
		fmt.Sprintf("address=%s", c.Address),
		fmt.Sprintf("poll-interval=%v", c.PollInterval),
		fmt.Sprintf("report-interval=%v", c.ReportInterval),
		fmt.Sprintf("rate-limit=%v", c.RateLimit),
	}

	if len(c.Secret) > 0 {
		str = append(str, fmt.Sprintf("secret=%v", c.Secret))
	}

	if len(c.PublicKeyPath) > 0 {
		str = append(str, fmt.Sprintf("public-key=%v", c.PublicKeyPath))
	}

	return "agent config: " + strings.Join(str, "; ")
}

func (c *Config) parse() error {
	err := c.tryLoadJSONConfig()
	if err != nil {
		return err
	}

	address := c.Address
	pflag.VarP(&address, "address", "a", "address:port for HTTP API requests")

	secret := c.Secret
	pflag.VarP(&secret, "secret", "k", "a key to sign outgoing data")

	publicKeyPath := c.PublicKeyPath
	pflag.VarP(&publicKeyPath, "crypto-key", "", "path to public key to encrypt agent -> server communications")

	configPath := entities.FilePath("") // register var for compatibility
	pflag.VarP(&configPath, "config", "c", "path to configuration file in JSON format")

	pflag.IntVarP(&c.PollInterval, "poll-interval", "p", c.PollInterval, "interval (s) for polling stats")
	pflag.IntVarP(&c.ReportInterval, "report-interval", "r", c.ReportInterval, "interval (s) for polling stats")
	pflag.IntVarP(&c.RateLimit, "rate-limit", "l", c.RateLimit, "number of max simultaneous requests to server")

	pflag.Parse()

	pflag.Visit(func(f *pflag.Flag) {
		switch f.Name {
		case "address":
			c.Address = address
		case "secret":
			c.Secret = secret
		case "crypto-key":
			c.PublicKeyPath = publicKeyPath
		}
	})

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
		return fmt.Errorf("agent.loadConfigFromFile - os.ReadFile: %w", err)
	}

	if err := json.Unmarshal(data, c); err != nil {
		return fmt.Errorf("agent.loadConfigFromFile - json.Unmarshal: %w", err)
	}

	return nil
}
