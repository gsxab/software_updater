package ui

type Config struct {
	DateFormat string `json:"date_format"`
	TimeFormat string `json:"time_format"`
}

var uiConfig *Config

func DefaultConfig() *Config {
	return &Config{
		DateFormat: "2006-01-02",
		TimeFormat: "2006-01-02 03:04:05",
	}
}
