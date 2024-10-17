package video

import (
	"ffmpeg_hls_go/internal/configs"
	"ffmpeg_hls_go/internal/logger"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type FFmpegObj struct {
	IsLive       bool
	Index        int16
	cmd          *exec.Cmd
	mu           sync.Mutex
	streamURL    string
	hlsURL       string
	ffmpegConfig *configs.FFmpegConfig
}

func NewFFmpegObj(_isLive bool, _index int16) *FFmpegObj {
	config := configs.GetConfigInstance()

	return &FFmpegObj{
		IsLive:       _isLive,
		Index:        _index,
		streamURL:    config.Streams.Live[_index].StreamURL,
		hlsURL:       config.Streams.Live[_index].HLSURL,
		ffmpegConfig: &config.FfmpegConfig,
	}
}
func (f *FFmpegObj) Json() map[string]string {
	result := map[string]string{
		"stream_url":   f.streamURL,
		"hls_url":      f.hlsURL,
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
	urls := strings.Split(f.streamURL, "?")
	streamURL := fmt.Sprintf("%s?starttime=%s&endtime=%s", urls[0], starttime, endtime)
	f.streamURL = streamURL
	f.Stop()
	return f.Start()
}
func (f *FFmpegObj) Start() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	log := logger.GetLoggerInstance()
	if f.streamURL == "" || f.hlsURL == "" {
		return fmt.Errorf("stream_url or hls_url is empty")
	}
	destDir := filepath.Dir(f.hlsURL)
	os.MkdirAll(destDir, 0755)
	cmd := exec.Command(f.ffmpegConfig.FfmpegPath,
		"-hide_banner",
		"-i",
		f.streamURL,
		"-c:v", "copy",
		"-an",
		"-start_number", "0",
		"-f", "hls",
		"-hls_time", strconv.Itoa(f.ffmpegConfig.Hls_time),
		"-hls_list_size", strconv.Itoa(f.ffmpegConfig.Hls_list_size),
		"-hls_flags", "delete_segments",
		f.hlsURL)
	f.cmd = cmd

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	err := cmd.Start()
	if err != nil {
		wrappedErr := fmt.Errorf("failed to start FFmpeg: %w", err)
		log.Println(wrappedErr)
		return wrappedErr
	}
	ffmpegLogger := logger.GetFFmpegLogger()
	go io.Copy(ffmpegLogger.Writer(), stdout)
	go io.Copy(ffmpegLogger.Writer(), stderr)

	log.Printf("%s", cmd.Args)
	//log.Printf("FFmpeg started: %s %s\n", f.streamConfig.StreamURL, f.streamConfig.HLSURL)
	log.Printf("FFMPEG pid:%d state:%s", cmd.Process.Pid, cmd.ProcessState.String())
	go func() {
		for {
			if err := cmd.Wait(); err != nil {
				log.Errorf("FFmpeg pid:%d exited with error %v", cmd.Process.Pid, err)
				if f.IsLive {
					f.Start()
				} else {
					break
				}
			}
		}
		log.Errorf("ffmpeg %d exited", cmd.Process.Pid)
	}()
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
	return f.hlsURL
}
