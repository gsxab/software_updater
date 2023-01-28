package config

type Config struct {
	Secret     *string `json:"secret,omitempty"`
	Addr       string  `json:"addr,omitempty"`
	DateFormat string  `json:"date_format"`
	TimeFormat string  `json:"time_format"`
}

var WebUIConfig *Config

func DefaultConfig() *Config {
	return &Config{
		Secret:     nil,
		Addr:       "127.0.0.1:8080",
		DateFormat: "2006-01-02",
		TimeFormat: "2006-01-02 03:04:05",
	}
}
