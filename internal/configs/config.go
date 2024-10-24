package configs

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	TokenServer struct {
		UserName    string `yaml:"username"`
		Password    string `yaml:"password"`
		TokenApiUrl string `yaml:"token_api_url"`
	} `yaml:"token_server"`
	Logging struct {
		LogPath    string `yaml:"log_path"`
		MaxSize    int    `yaml:"max_size"`
		MaxBackups int    `yaml:"max_backups"`
		MaxAge     int    `yaml:"max_age"`
		Compress   bool   `yaml:"compress"`
		Loglevel   int    `yaml:"Loglevel"`
	} `yaml:"logging"`

	FfmpegConfig FFmpegConfig `yaml:"FFmpegConfig"`
	Streams      struct {
		Live   []StreamConfig `yaml:"live"`
		Replay []StreamConfig `yaml:"replay"`
	} `yaml:"streams"`
}

type FFmpegConfig struct {
	FfmpegPath    string `yaml:"ffmpeg_path"`
	Hls_time      int    `yaml:"hls_time"`
	Hls_list_size int    `yaml:"hls_list_size"`
	Loglevel      string `yaml:"loglevel"`
}
type StreamConfig struct {
	StreamURL string `yaml:"stream_url"`
	HLSURL    string `yaml:"hls_url"`
}

func CreateDefaultConfig() *Config {
	config := &Config{}
	config.Server.Port = 8080
	config.Logging.LogPath = "log/"
	config.Logging.MaxSize = 10
	config.Logging.MaxBackups = 3
	config.Logging.MaxAge = 30
	config.Logging.Loglevel = 1
	config.Logging.Compress = true

	config.FfmpegConfig.FfmpegPath = "ffmpeg"
	config.FfmpegConfig.Hls_time = 1
	config.FfmpegConfig.Hls_list_size = 5
	config.FfmpegConfig.Loglevel = "error"

	// 以下是配置 HLS 流
	config.Streams.Live = []StreamConfig{
		{
			StreamURL: "rtsp://admin:ist20171016@192.168.1.55:12409/Streaming/Channels/102",
			HLSURL:    "static/live1/live.m3u8",
		},
		{
			StreamURL: "rtsp://admin:ist20171016@192.168.1.55:12409/Streaming/Channels/402",
			HLSURL:    "static/live2/live.m3u8",
		},
	}
	config.Streams.Replay = []StreamConfig{
		{
			StreamURL: "rtsp://admin:ist20171016@192.168.1.55:12409/Streaming/tracks/101?starttime=20240926t090000z&endtime=20240926t092000z",
			HLSURL:    "static/Replay1/replay.m3u8",
		},
		{
			StreamURL: "rtsp://admin:ist20171016@192.168.1.55:12409/Streaming/tracks/401?starttime=20240926t090000z&endtime=20240926t092000z",
			HLSURL:    "static/Replay2/replay.m3u8",
		},
	}

	config.TokenServer.UserName = "42536518181156892813"
	config.TokenServer.Password = "LN1c0eZsilRvgl2Mt5bZJIzeqtYkqN"
	config.TokenServer.TokenApiUrl = "https://t.jdc.taep.org.cn:17701/api/v1/token/getAccessToken"
	return config
}

// WriteConfigToFile writes the config to a file in YAML format
func WriteConfigToFile(config *Config, filename string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// LoadConfig reads and parses the YAML configuration file
func loadConfig(filename string) (*Config, error) {
	// Read the file content
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("config file not found. writing default config...")

			defaultConfig := CreateDefaultConfig()
			// Write the default config to the file
			if writeErr := WriteConfigToFile(defaultConfig, filename); writeErr != nil {
				return nil, fmt.Errorf("failed to write default config: %w", writeErr)
			}

			return defaultConfig, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	config := CreateDefaultConfig()
	// Parse the YAML content
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}

var (
	configInstance *Config
	configOnce     sync.Once
)

func GetConfigInstance() *Config {
	configOnce.Do(func() {
		configInstance, _ = loadConfig("config.yaml")
	})
	return configInstance
}
