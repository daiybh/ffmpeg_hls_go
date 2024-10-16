package main

import "log"

type FFmpegMgr struct {
	liveObjs   [2]*FFmpegObj
	replayObjs [2]*FFmpegObj
}

func NewFFmpegMgr() *FFmpegMgr {
	return &FFmpegMgr{}
}

func (mgr *FFmpegMgr) Start(config *Config) {
	log.Println("FFmpegMgr start....")
	// Starting two live FFmpegObj
	mgr.liveObjs[0] = NewFFmpegObj(&config.Streams.Live[0], &config.FfmpegConfig)
	mgr.liveObjs[1] = NewFFmpegObj(&config.Streams.Live[1], &config.FfmpegConfig)

	// Starting two replay FFmpegObj
	streamConfig := &StreamConfig{}

	mgr.replayObjs[0] = NewFFmpegObj(streamConfig, &config.FfmpegConfig)
	mgr.replayObjs[1] = NewFFmpegObj(streamConfig, &config.FfmpegConfig)

	// Start all FFmpeg processes
	for _, obj := range mgr.liveObjs {
		obj.Start()
	}
	for _, obj := range mgr.replayObjs {
		obj.Start()
	}
}

func (mgr *FFmpegMgr) GetLiveObj(index int) *FFmpegObj {
	if index >= 0 && index < len(mgr.liveObjs) {
		return mgr.liveObjs[index]
	}
	return nil
}

func (mgr *FFmpegMgr) GetReplayObj(index int) *FFmpegObj {
	if index >= 0 && index < len(mgr.replayObjs) {
		return mgr.replayObjs[index]
	}
	return nil
}
