package handles

import (
	"encoding/json"
	"ffmpeg_hls_go/internal/logger"
	"ffmpeg_hls_go/internal/video"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func PlayHandler(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger("handlers.log", false)
	log.Printf("playHandler: %s %s", r.URL.Path, r.URL.Query())
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) != 3 {
		log.Printf("Invalid path: %s", r.URL.Path)
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	index := 0
	log.Printf("pathParts[2]: %s", pathParts[2])
	if strings.Contains(pathParts[2], ".m3u8") {
		m3u8Name := strings.Split(pathParts[2], ".")[0]

		index, _ = strconv.Atoi(m3u8Name)
	} else {
		index, _ = strconv.Atoi(pathParts[2])
	}
	log.Printf("index,%d", index)

	if index < 1 || index > 2 {
		log.Printf("Invalid live stream number: %s", pathParts[2])
		http.Error(w, "Invalid live stream number", http.StatusBadRequest)
		return
	}

	index = index - 1
	ffmpegObj := video.GetFFmpegMgr().GetLiveObj(index)
	if r.URL.Query().Has("starttime") && r.URL.Query().Has("endtime") {
		log.Printf("playHandler: starttime: %s, endtime: %s", r.URL.Query().Get("starttime"), r.URL.Query().Get("endtime"))
		ffmpegObj = video.GetFFmpegMgr().GetReplayObj(index)

		error := ffmpegObj.StartReplay(r.URL.Query().Get("starttime"), r.URL.Query().Get("endtime"))
		if error != nil {
			http.Error(w, error.Error(), http.StatusNotFound)
			return
		}
	}
	if ffmpegObj != nil {
		log.Printf("find playHandler: %s", ffmpegObj.GetHLSURL())
		http.Redirect(w, r, "/"+ffmpegObj.GetHLSURL(), http.StatusFound)
		//http.Redirect(w, r, "/static/server.log", http.StatusFound)
	} else {
		log.Printf("stream not found")
		http.Error(w, "stream not found", http.StatusNotFound)
	}
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("statusHandler: %s %s", r.URL.Path, r.URL.Query())

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	objs := make([]map[string]string, 4)
	for i, obj := range video.GetFFmpegMgr().LiveObjs {
		objs[i] = obj.Json()
	}
	for i, obj := range video.GetFFmpegMgr().ReplayObjs {
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
