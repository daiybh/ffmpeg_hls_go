package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type FFmpegObj struct {
	streamConfig *StreamConfig
	ffmpegConfig *FFmpegConfig
	cmd          *exec.Cmd
	mu           sync.Mutex
}

func NewFFmpegObj(streamConfig *StreamConfig, ffmpegConfig *FFmpegConfig) *FFmpegObj {
	return &FFmpegObj{
		streamConfig: streamConfig,
		ffmpegConfig: ffmpegConfig,
	}
}
func (f *FFmpegObj) Json() map[string]string {
	result := map[string]string{
		"stream_url":   f.streamConfig.StreamURL,
		"hls_url":      f.streamConfig.HLSURL,
		"cmd":          "",
		"processState": "",
	}
	if f.cmd != nil {
		result["cmd"] = f.cmd.String()
		result["processState"] = f.cmd.ProcessState.String()
	}
	return result
}
func (f *FFmpegObj) StartReplay(starttime string, endtime string) error {
	urls := strings.Split(f.streamConfig.StreamURL, "?")
	streamURL := fmt.Sprintf("%s?starttime=%s&endtime=%s", urls[0], starttime, endtime)
	f.streamConfig.StreamURL = streamURL
	f.Stop()
	return f.Start()
}
func (f *FFmpegObj) Start() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.streamConfig.StreamURL == "" || f.streamConfig.HLSURL == "" {
		return fmt.Errorf("stream_url or hls_url is empty")
	}
	destDir := filepath.Dir(f.streamConfig.HLSURL)
	os.MkdirAll(destDir, 0755)
	cmd := exec.Command(f.ffmpegConfig.FfmpegPath,
		"-i",
		f.streamConfig.StreamURL,
		"-c:v", "copy",
		"-an",
		"-start_number", "0",
		"-f", "hls",
		"-hls_time", strconv.Itoa(f.ffmpegConfig.Hls_time),
		"-hls_list_size", strconv.Itoa(f.ffmpegConfig.Hls_list_size),
		"-hls_flags", "delete_segments",
		f.streamConfig.HLSURL)
	f.cmd = cmd
	err := cmd.Start()
	if err != nil {
		wrappedErr := fmt.Errorf("failed to start FFmpeg: %w", err)
		log.Println(wrappedErr)
		return wrappedErr
	}
	log.Printf("%s", cmd.Args)
	//log.Printf("FFmpeg started: %s %s\n", f.streamConfig.StreamURL, f.streamConfig.HLSURL)
	log.Printf("FFMPEG pid:%d state:%s", cmd.Process.Pid, cmd.ProcessState.String())
	go cmd.Wait() // Run FFmpeg asynchronously
	return nil
}

func (f *FFmpegObj) Stop() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.cmd != nil {
		return f.cmd.Process.Kill()
	}
	return nil
}

func (f *FFmpegObj) GetHLSURL() string {
	return f.streamConfig.HLSURL
}
