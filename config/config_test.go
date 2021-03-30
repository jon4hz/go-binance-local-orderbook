package config

import (
	"fmt"
	"os"
	"testing"
)

func TestGetBeforeConfigIsLoaded(t *testing.T) {
	defer func() { recover() }()
	Get()
	t.Fatal("Should've panicked because the configuration hasn't been loaded yet")
}

func TestLoadFileThatDoesNotExistAndWithoutEnvVars(t *testing.T) {
	defer func() { recover() }()
	_ = Load("file-that-does-not-exist.yaml")
	t.Error("Should've panicked, because the file specified doesn't exist and env is not set")
}

func TestLoadDefaultConfigurationFile(t *testing.T) {
	defer func() { recover() }()
	_ = Load(DefaultConfigurationFilePath)
	t.Error("Should've panicked, because there's no configuration files at the default path nor the default fallback path nor is the env set")

}

func TestConfigWithEnvVariables(t *testing.T) {
	os.Setenv("NAME", "binance")
	os.Setenv("MARKET", "BTCUSDT")

	os.Setenv("POSTGRES_DB", "orderbook")
	os.Setenv("POSTGRES_USER", "username")
	os.Setenv("POSTGRES_PASSWORD", "password")
	os.Setenv("POSTGRES_SERVER", "localhost")
	os.Setenv("POSTGRES_PORT", "5432")

	err := Load("")
	if err != nil {
		t.Error(fmt.Sprintf("Couldn't read variables from environment: %s", err))
	}
	cfg := Get()
	if cfg == nil {
		t.Error("Config can't be nil")
	}
}

func TestWithMissingVarExchangeName(t *testing.T) {
	defer func() { recover() }()
	os.Setenv("MARKET", "BTCUSDT")

	os.Setenv("POSTGRES_DB", "orderbook")
	os.Setenv("POSTGRES_USER", "username")
	os.Setenv("POSTGRES_PASSWORD", "password")
	os.Setenv("POSTGRES_SERVER", "localhost")
	os.Setenv("POSTGRES_PORT", "5432")

	_, varSet := os.LookupEnv("NAME")
	if varSet {
		os.Unsetenv("NAME")
	}

	err := Load("")
	if err != nil {
		t.Error(fmt.Sprintf("Error while loading the config from env: %s", err))
	}
	t.Error("Should have panicked, because exchange.name is not set")
}

func TestWithMissingVarExchangeMarket(t *testing.T) {
	defer func() { recover() }()
	os.Setenv("NAME", "binance")

	os.Setenv("POSTGRES_DB", "orderbook")
	os.Setenv("POSTGRES_USER", "username")
	os.Setenv("POSTGRES_PASSWORD", "password")
	os.Setenv("POSTGRES_SERVER", "localhost")
	os.Setenv("POSTGRES_PORT", "5432")

	_, varSet := os.LookupEnv("MARKET")
	if varSet {
		os.Unsetenv("MARKET")
	}

	err := Load("")
	if err != nil {
		t.Error(fmt.Sprintf("Error while loading the config from env: %s", err))
	}
	t.Error("Should have panicked, because exchange.name is not set")
}

func TestWithMissingVarExchangeAPIKey(t *testing.T) {
	defer func() { recover() }()
	os.Setenv("NAME", "binance")

	os.Setenv("POSTGRES_DB", "orderbook")
	os.Setenv("POSTGRES_USER", "username")
	os.Setenv("POSTGRES_PASSWORD", "password")
	os.Setenv("POSTGRES_SERVER", "localhost")
	os.Setenv("POSTGRES_PORT", "5432")

	_, varSet := os.LookupEnv("API_KEY")
	if varSet {
		os.Unsetenv("API_KEY")
	}

	err := Load("")
	if err != nil {
		t.Error(fmt.Sprintf("Error while loading the config from env: %s", err))
	}
	t.Error("Should have panicked, because exchange.api_key is not set")
}

func TestWithMissingVarExchangeAPISecret(t *testing.T) {
	defer func() { recover() }()
	os.Setenv("NAME", "binance")

	os.Setenv("POSTGRES_DB", "orderbook")
	os.Setenv("POSTGRES_USER", "username")
	os.Setenv("POSTGRES_PASSWORD", "password")
	os.Setenv("POSTGRES_SERVER", "localhost")
	os.Setenv("POSTGRES_PORT", "5432")

	_, varSet := os.LookupEnv("API_SECRET")
	if varSet {
		os.Unsetenv("API_SECRET")
		os.Unsetenv(("ASDFasdfa"))
	}

	err := Load("")
	if err != nil {
		t.Error(fmt.Sprintf("Error while loading the config from env: %s", err))
	}
	t.Error("Should have panicked, because exchange.api_secret is not set")
}

func TestWithMissingVarDatabaseDB(t *testing.T) {
	defer func() { recover() }()
	os.Setenv("NAME", "binance")
	os.Setenv("MARKET", "BTCUSDT")

	os.Setenv("POSTGRES_USER", "username")
	os.Setenv("POSTGRES_PASSWORD", "password")
	os.Setenv("POSTGRES_SERVER", "localhost")
	os.Setenv("POSTGRES_PORT", "5432")

	_, varSet := os.LookupEnv("POSTGRES_DB")
	if varSet {
		os.Unsetenv("POSTGRES_DB")
	}

	err := Load("")
	if err != nil {
		t.Error(fmt.Sprintf("Error while loading the config from env: %s", err))
	}
	t.Error("Should have panicked, because database.db is not set")
}

func TestWithMissingVarDatabaseUsername(t *testing.T) {
	defer func() { recover() }()
	os.Setenv("NAME", "binance")
	os.Setenv("MARKET", "BTCUSDT")

	os.Setenv("POSTGRES_DB", "orderbook")
	os.Setenv("POSTGRES_PASSWORD", "password")
	os.Setenv("POSTGRES_SERVER", "localhost")
	os.Setenv("POSTGRES_PORT", "5432")

	_, varSet := os.LookupEnv("POSTGRES_USER")
	if varSet {
		os.Unsetenv("POSTGRES_USER")
	}

	err := Load("")
	if err != nil {
		t.Error(fmt.Sprintf("Error while loading the config from env: %s", err))
	}
	t.Error("Should have panicked, because database.username is not set")
}

func TestWithMissingVarDatabasePassword(t *testing.T) {
	defer func() { recover() }()
	os.Setenv("NAME", "binance")
	os.Setenv("MARKET", "BTCUSDT")

	os.Setenv("POSTGRES_DB", "orderbook")
	os.Setenv("POSTGRES_USER", "username")
	os.Setenv("POSTGRES_SERVER", "localhost")
	os.Setenv("POSTGRES_PORT", "5432")

	_, varSet := os.LookupEnv("POSTGRES_PASSWORD")
	if varSet {
		os.Unsetenv("POSTGRES_PASSWORD")
	}

	err := Load("")
	if err != nil {
		t.Error(fmt.Sprintf("Error while loading the config from env: %s", err))
	}
	t.Error("Should have panicked, because database.password is not set")
}

func TestWithMissingVarDatabaseServer(t *testing.T) {
	defer func() { recover() }()
	os.Setenv("NAME", "binance")
	os.Setenv("MARKET", "BTCUSDT")

	os.Setenv("POSTGRES_DB", "orderbook")
	os.Setenv("POSTGRES_USER", "username")
	os.Setenv("POSTGRES_PASSWORD", "password")
	os.Setenv("POSTGRES_SERVER", "localhost")
	os.Setenv("POSTGRES_PORT", "5432")

	_, varSet := os.LookupEnv("POSTGRES_SERVER")
	if varSet {
		os.Unsetenv("POSTGRES_SERVER")
	}

	err := Load("")
	if err != nil {
		t.Error(fmt.Sprintf("Error while loading the config from env: %s", err))
	}
	t.Error("Should have panicked, because database.server is not set")
}

func TestWithMissingVarDatabasePort(t *testing.T) {
	defer func() { recover() }()
	os.Setenv("NAME", "binance")

	os.Setenv("POSTGRES_DB", "orderbook")
	os.Setenv("POSTGRES_USER", "username")
	os.Setenv("POSTGRES_PASSWORD", "password")
	os.Setenv("POSTGRES_SERVER", "localhost")
	os.Setenv("POSTGRES_PORT", "5432")

	_, varSet := os.LookupEnv("POSTGRES_PORT")
	if varSet {
		os.Unsetenv("POSTGRES_PORT")
	}

	err := Load("")
	if err != nil {
		t.Error(fmt.Sprintf("Error while loading the config from env: %s", err))
	}
	t.Log("Didn't panic, because database.port has default value")
}

func TestConfigWithEnvVariablesInvalidExchange(t *testing.T) {
	defer func() { recover() }()
	os.Setenv("NAME", "kucoin")
	os.Setenv("MARKET", "BTCUSDT")

	os.Setenv("POSTGRES_DB", "orderbook")
	os.Setenv("POSTGRES_USER", "username")
	os.Setenv("POSTGRES_PASSWORD", "password")
	os.Setenv("POSTGRES_SERVER", "localhost")
	os.Setenv("POSTGRES_PORT", "5432")

	err := Load("")
	if err != nil {
		t.Error(fmt.Sprintf("Couldn't read variables from environment: %s", err))
	}
	cfg := Get()
	if cfg == nil {
		t.Error("Config can't be nil")
	}
	t.Error("Should have panicked, because exchange.name is invalid")
}

func TestConfigWithoutDatabaseConfig(t *testing.T) {
	defer func() { recover() }()
	// unset all possible env vars
	vars := [7]string{"NAME", "MARKET", "POSTGRES_DB", "POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_SERVER", "POSTGRES_PORT"}
	var varSet bool
	for i := 0; i < len(vars); i++ {
		_, varSet = os.LookupEnv(vars[i])
		if varSet {
			os.Unsetenv(vars[i])
		}
	}
	os.Setenv("NAME", "binance")
	os.Setenv("MARKET", "BTCUSDT")
	os.Setenv("API_KEY", "asdf")
	os.Setenv("API_SECRET", "asdf1234")

	_ = Load("")
	t.Error("Should have panicked, because database configuration is missing")
}

func TestConfigReadFromFile(t *testing.T) {

	// unset all possible env vars
	vars := [7]string{"NAME", "MARKET", "POSTGRES_DB", "POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_SERVER", "POSTGRES_PORT"}
	var varSet bool
	for i := 0; i < len(vars); i++ {
		_, varSet = os.LookupEnv(vars[i])
		if varSet {
			os.Unsetenv(vars[i])
		}
	}

	err := Load("valid_test_config.yml")
	if err != nil {
		t.Error(fmt.Sprintf("Error while loading the config from file: %s", err))
	}
}

func TestConfigReadFromInvalidFile(t *testing.T) {

	// unset all possible env vars
	vars := [7]string{"NAME", "MARKET", "POSTGRES_DB", "POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_SERVER", "POSTGRES_PORT"}
	var varSet bool
	for i := 0; i < len(vars); i++ {
		_, varSet = os.LookupEnv(vars[i])
		if varSet {
			os.Unsetenv(vars[i])
		}
	}

	err := Load("invalid_test_config.txt")
	if err == nil {
		t.Error("Error can't be nil since config file is invalid")
	}
}
