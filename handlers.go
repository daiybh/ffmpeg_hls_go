package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

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
	index = index - 1
	ffmpegObj := ffmpegMgr.GetLiveObj(index)
	if r.URL.Query().Has("starttime") && r.URL.Query().Has("endtime") {
		log.Printf("playHandler: starttime: %s, endtime: %s\n", r.URL.Query().Get("starttime"), r.URL.Query().Get("endtime"))
		ffmpegObj = ffmpegMgr.GetReplayObj(index)

		ffmpegObj.StartReplay(r.URL.Query().Get("starttime"), r.URL.Query().Get("endtime"))
	}
	if ffmpegObj != nil {
		log.Printf("find playHandler: %s\n", ffmpegObj.GetHLSURL())
		http.Redirect(w, r, "/"+ffmpegObj.GetHLSURL(), http.StatusFound)
		//http.Redirect(w, r, "/static/server.log", http.StatusFound)
	} else {
		http.Error(w, "stream not found", http.StatusNotFound)
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("statusHandler: %s %s\n", r.URL.Path, r.URL.Query())

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	objs := make([]map[string]string, 4)
	for i, obj := range ffmpegMgr.liveObjs {
		objs[i] = obj.Json()
	}
	for i, obj := range ffmpegMgr.replayObjs {
		objs[i+2] = obj.Json()
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
