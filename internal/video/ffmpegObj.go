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
	"time"
)

type FFmpegObj struct {
	IsLive       bool
	Index        int16
	cmd          *exec.Cmd
	mu           sync.Mutex
	streamURL    string
	hlsURL       string
	ffmpegConfig *configs.FFmpegConfig
	StartedCount uint64
}

func NewFFmpegObj(_isLive bool, _index int16) *FFmpegObj {
	config := configs.GetConfigInstance()

	return &FFmpegObj{
		IsLive:       _isLive,
		Index:        _index,
		streamURL:    config.Streams.Live[_index].StreamURL,
		hlsURL:       config.Streams.Live[_index].HLSURL,
		ffmpegConfig: &config.FfmpegConfig,
		StartedCount: 0,
	}
}
func (f *FFmpegObj) Json() map[string]string {
	result := map[string]string{
		"stream_url":   f.streamURL,
		"hls_url":      f.hlsURL,
		"cmd":          "",
		"processState": "",
		"StartedCount": strconv.FormatUint(f.StartedCount, 10),
	}
	if f.cmd != nil {
		result["cmd"] = f.cmd.String()
		result["processState"] = f.cmd.ProcessState.String()
	}
	return result
}
func (f *FFmpegObj) StartReplay(starttime string, endtime string) error {
	if f.streamURL == "" {
		return fmt.Errorf("stream_url is empty")
	}
	urls := strings.Split(f.streamURL, "?")
	streamURL := fmt.Sprintf("%s?starttime=%s&endtime=%s", urls[0], starttime, endtime)
	f.streamURL = streamURL
	f.Stop()
	return f.Start()
}

// 启动 ffmpeg 进程
func (f *FFmpegObj) startFFmpeg() (*exec.Cmd, error) {
	// 这里指定你的 ffmpeg 命令和参数
	f.mu.Lock()
	defer f.mu.Unlock()
	f.StartedCount += 1
	log := logger.GetLogger("ffmpegobj.log", false)
	if f.streamURL == "" || f.hlsURL == "" {
		return nil, fmt.Errorf("stream_url or hls_url is empty")
	}
	destDir := filepath.Dir(f.hlsURL)
	os.MkdirAll(destDir, 0755)
	cmd := exec.Command(f.ffmpegConfig.FfmpegPath,
		"-hide_banner",
		"-loglevel", f.ffmpegConfig.Loglevel,
		"-i",
		f.streamURL,
		"-c:v", "copy",
		"-an",
		"-start_number", "0",
		"-f", "hls",
		"-hls_time", strconv.Itoa(f.ffmpegConfig.Hls_time),
		"-hls_list_size", strconv.Itoa(f.ffmpegConfig.Hls_list_size),
		"-hls_flags", "delete_segments",
		//"-metadata title", f.hlsURL,
		f.hlsURL)
	f.cmd = cmd

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	err := cmd.Start()
	if err != nil {
		wrappedErr := fmt.Errorf("failed to start FFmpeg: %w", err)
		log.Println(wrappedErr)
		return cmd, wrappedErr
	}
	ffmpegLogger := logger.GetFFmpegLogger()
	go io.Copy(ffmpegLogger.Writer(), stdout)
	go io.Copy(ffmpegLogger.Writer(), stderr)
	log.Infof("===%d====", f.StartedCount)
	log.Printf("%s", cmd.Args)
	log.Printf("FFmpeg started: %s %s", f.streamURL, f.hlsURL)
	log.Printf("FFMPEG pid:%d state:%s", cmd.Process.Pid, cmd.ProcessState.String())
	return cmd, nil
}

// 监控 ffmpeg 进程
func (f *FFmpegObj) Start() error {
	cmd, _ := f.startFFmpeg()
	if f.IsLive {
		go f.monitorFFmpeg(cmd)
	}
	return nil
}

func (f *FFmpegObj) monitorFFmpeg(cmd *exec.Cmd) {
	for {
		// 检查进程是否退出
		err := cmd.Process.Signal(nil)
		if err != nil {
			// 进程已经退出
			cmd, _ = f.startFFmpeg() // 重新启动 ffmpeg
		}

		// 每隔 5 秒检查一次进程状态
		time.Sleep(1 * time.Second)
	}
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
