go build -ldflags "-X 'main.goVersion=$(go version)' -X 'main.gitHash=$(git show -s --format=%H)' -X 'main.buildTime=$(git show -s --format=%cd)'" -o ffmpeg_hls_go.exe ./cmd/ffmpeg_hls_go/
