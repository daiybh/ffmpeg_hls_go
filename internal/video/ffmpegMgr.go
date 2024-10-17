package video

import (
	"ffmpeg_hls_go/internal/configs"
	"ffmpeg_hls_go/internal/logger"
	"sync"
)

type FFmpegMgr struct {
	LiveObjs   [2]*FFmpegObj
	ReplayObjs [2]*FFmpegObj
}

var (
	ffmpegMgrInstance *FFmpegMgr
	once              sync.Once
)

func GetFFmpegMgr() *FFmpegMgr {
	once.Do(func() {
		ffmpegMgrInstance = &FFmpegMgr{}

	})
	return ffmpegMgrInstance
}
func (mgr *FFmpegMgr) Start(config *configs.Config) {
	log := logger.GetLoggerInstance()
	log.Println("FFmpegMgr start....")
	// Starting two live FFmpegObj
	mgr.LiveObjs[0] = NewFFmpegObj(true, 0)
	mgr.LiveObjs[1] = NewFFmpegObj(true, 1)

	// Starting two replay FFmpegObj

	mgr.ReplayObjs[0] = NewFFmpegObj(false, 0)
	mgr.ReplayObjs[1] = NewFFmpegObj(false, 1)

	// Start all FFmpeg processes
	for _, obj := range mgr.LiveObjs {
		obj.Start()
	}
}

func (mgr *FFmpegMgr) GetLiveObj(index int) *FFmpegObj {
	if index >= 0 && index < len(mgr.LiveObjs) {
		return mgr.LiveObjs[index]
	}
	return nil
}

func (mgr *FFmpegMgr) GetReplayObj(index int) *FFmpegObj {
	if index >= 0 && index < len(mgr.ReplayObjs) {
		return mgr.ReplayObjs[index]
	}
	return nil
}
