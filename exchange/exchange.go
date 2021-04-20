package exchange

type Config struct {
	Name   string `mapstructure:"NAME"`
	Market string `mapstructure:"MARKET"`
}

var (
	LastUpdateID int64
	SmallU       int64
	BigU         int64
	Prev_u       int64
	Notified     bool
)
