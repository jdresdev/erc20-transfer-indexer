package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	DatabaseURL   string
	RPCURL        string
	BatchSize     uint64
	Confirmations uint64
	PollInterval  time.Duration
}

func Load() (*Config, error) {
	batchSize, err := getUint("BATCH_SIZE", 500)
	if err != nil {
		return nil, err
	}

	confirmations, err := getUint("CONFIRMATIONS", 6)
	if err != nil {
		return nil, err
	}

	pollInterval, err := getDuration("POLL_INTERVAL", "3s")
	if err != nil {
		return nil, err
	}

	return &Config{
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		RPCURL:        os.Getenv("RPC_URL"),
		BatchSize:     batchSize,
		Confirmations: confirmations,
		PollInterval:  pollInterval,
	}, nil
}

func getUint(key string, defaultVal uint64) (uint64, error) {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal, nil
	}

	parsed, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, err
	}

	return parsed, nil
}

func getDuration(key string, defaultVal string) (time.Duration, error) {
	val := os.Getenv(key)
	if val == "" {
		val = defaultVal
	}

	return time.ParseDuration(val)
}
