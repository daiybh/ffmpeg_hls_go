package main

import (
	"ffmpeg_hls_go/internal/video"
	"ffmpeg_hls_go/internal/video/handles"
	"ffmpeg_hls_go/pkg/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/natefinch/lumberjack"
)

// Configure the logger to rotate automatically
func setupLogger(config *utils.Config) io.Writer {
	loggerFile := &lumberjack.Logger{
		Filename:   config.Logging.LogFile,    // Log file name
		MaxSize:    config.Logging.MaxSize,    // Max size in MB before rotating
		MaxBackups: config.Logging.MaxBackups, // Max number of old log files to keep
		MaxAge:     config.Logging.MaxAge,     // Max age in days to keep old log files
		Compress:   true,                      // Compress old log files
	}
	multiWriter := io.MultiWriter(os.Stdout, loggerFile)
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("Logger initialized")
	return loggerFile
}

var ffmpegMgr *video.FFmpegMgr

func main() {
	config, err := utils.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	logger := setupLogger(config)
	defer logger.(*lumberjack.Logger).Close()

	log.Printf("Starting HLS server...")
	// Start the FFmpeg Manager
	ffmpegMgr = video.GetFFmpegMgr()
	ffmpegMgr.Start(config)

	// Setup HTTP routes

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/ch/", handles.PlayHandler)
	http.HandleFunc("/status/", handles.StatusHandler)

	// Start the HTTP server
	address := fmt.Sprintf(":%d", config.Server.Port)
	log.Println("Server is listening on port " + address + "...")
	log.Fatal(http.ListenAndServe(address, nil))
}
