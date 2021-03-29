package exchange

type Config struct {
	Name      string `mapstructure:"NAME"`
	Market    string `mapstructure:"MARKET"`
	APIKey    string `mapstructure:"API_KEY"`
	APISecret string `mapstructure:"API_SECRET"`
}
