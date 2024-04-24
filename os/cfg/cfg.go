package cfg

type Config struct {
	adapter Adapter
}

func NewWithAdapter(adapter Adapter) *Config {
	return &Config{
		adapter: adapter,
	}
}
