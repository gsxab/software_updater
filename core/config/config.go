package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Files    *FileConfig       `yaml:"files,omitempty"`
	Database *DatabaseConfig   `yaml:"database,omitempty"`
	Selenium *SeleniumConfig   `yaml:"selenium,omitempty"`
	Engine   *EngineConfig     `yaml:"engine,omitempty"`
	Extra    map[string]string `yaml:"extra,omitempty"`
}

type FileConfig struct {
	ScreenshotDir string `yaml:"screenshot_dir,omitempty"`
	HTMLDir       string `yaml:"html_dir,omitempty"`
	CURLSaveDir   string `yaml:"curl_save_dir,omitempty"`
}

type DatabaseConfig struct {
	Driver string `yaml:"driver,omitempty"`
	DSN    string `yaml:"dsn,omitempty"`
}

type SeleniumConfig struct {
	Params     []string `yaml:"params,omitempty"`
	DriverPath string   `yaml:"driver_path,omitempty"`
}

type EngineConfig struct {
	ForceCrawl  bool `yaml:"force_crawl,omitempty"`
	ForceUpdate bool `yaml:"force_update,omitempty"`
	DebugLog    bool `yaml:"debug_log,omitempty"`
	DebugCheck  bool `yaml:"debug_check,omitempty"`
}

var config = DefaultConfig()

func DefaultConfig() *Config {
	return &Config{
		Files: &FileConfig{
			ScreenshotDir: "./screenshot/",
			HTMLDir:       "./html/",
		},
		Database: &DatabaseConfig{
			Driver: "sqlite",
			DSN:    "./software.db",
		},
		Selenium: &SeleniumConfig{
			Params:     nil,
			DriverPath: "./chromedriver",
		},
		Engine: &EngineConfig{
			ForceCrawl:  false,
			ForceUpdate: false,
			DebugLog:    false,
			DebugCheck:  false,
		},
		Extra: make(map[string]string),
	}
}

func Load(path string) (*Config, error) {
	log.Printf("loading configuration at: %s", path)

	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Printf("readfile failed: %s", err)
		return nil, err
	}

	// default config
	err = yaml.Unmarshal(bytes, config)
	if err != nil {
		log.Printf("yaml unmarshal failed: %s", err)
		return nil, err
	}

	return config, nil
}

func Current() *Config {
	return config
}
