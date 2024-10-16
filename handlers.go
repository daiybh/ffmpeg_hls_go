package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func liveHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) != 3 {
		log.Printf("Invalid path: %s\n", r.URL.Path)
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	index, err := strconv.Atoi(pathParts[2])
	if err != nil || index < 1 || index > 2 {
		log.Printf("Invalid live stream number: %s\n", pathParts[2])
		http.Error(w, "Invalid live stream number", http.StatusBadRequest)
		return
	}
	liveObj := ffmpegMgr.GetLiveObj(index)
	if liveObj != nil {
		http.Redirect(w, r, liveObj.GetHLSURL(), http.StatusFound)
	} else {
		http.Error(w, "Live stream not found", http.StatusNotFound)
	}
}
func playHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("playHandler: %s %s\n", r.URL.Path, r.URL.Query())
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) != 3 {
		log.Printf("Invalid path: %s\n", r.URL.Path)
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	m3u8Name := strings.Split(pathParts[2], ".")[0]
	index, err := strconv.Atoi(m3u8Name)
	if err != nil || index < 1 || index > 2 {
		log.Printf("Invalid live stream number: %s\n", pathParts[2])
		http.Error(w, "Invalid live stream number", http.StatusBadRequest)
		return
	}
	ffmpegObj := ffmpegMgr.GetLiveObj(index)
	if r.URL.Query().Has("starttime") && r.URL.Query().Has("endtime") {
		log.Printf("playHandler: starttime: %s, endtime: %s\n", r.URL.Query().Get("starttime"), r.URL.Query().Get("endtime"))
		ffmpegObj = ffmpegMgr.GetReplayObj(index)
	}
	if ffmpegObj != nil {
		log.Printf("playHandler: %s\n", ffmpegObj.GetHLSURL())
		//http.Redirect(w, r, ffmpegObj.GetHLSURL(), http.StatusFound)
		http.Redirect(w, r, "/static/server.log", http.StatusFound)
	} else {
		http.Error(w, "stream not found", http.StatusNotFound)
	}
}
func replayHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) != 3 {
		log.Printf("Invalid path: %s\n", r.URL.Path)
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	index, err := strconv.Atoi(pathParts[2])
	if err != nil || index < 1 || index > 2 {
		log.Printf("Invalid replay stream number: %s\n", pathParts[2])
		http.Error(w, "Invalid replay stream number", http.StatusBadRequest)
		return
	}
	replayObj := ffmpegMgr.GetReplayObj(index)
	if replayObj != nil {
		http.Redirect(w, r, replayObj.GetHLSURL(), http.StatusFound)
	} else {
		http.Error(w, "Replay stream not found", http.StatusNotFound)
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("statusHandler: %s %s\n", r.URL.Path, r.URL.Query())

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	objs := make([]map[string]string, len(ffmpegMgr.liveObjs))
	for i, obj := range ffmpegMgr.liveObjs {
		objs[i] = map[string]string{
			"hls_url":    obj.GetHLSURL(),
			"stream_url": obj.streamConfig.StreamURL,
			"cmd":        obj.cmd.ProcessState.String(),
		}
	}
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)

	err := encoder.Encode(objs)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
