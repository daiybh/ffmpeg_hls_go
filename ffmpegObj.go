package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
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

func (f *FFmpegObj) Start() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	cmd := exec.Command(f.ffmpegConfig.FfmpegPath,
		"-i",
		f.streamConfig.StreamURL,
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
	log.Printf("cmd:%s %s", cmd.Path, cmd.Args)
	log.Printf("FFmpeg started: %s %s\n", f.streamConfig.StreamURL, f.streamConfig.HLSURL)
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
