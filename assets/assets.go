package assets

import (
	"embed"
	"fmt"
	"math/rand/v2"
)

//go:embed videos/pls-rember.mp4
var defaultVideo []byte

//go:embed videos/sember.mp4
var semberVideo []byte

//go:embed videos/evil.webm
var evilVideo []byte

//go:embed videos/cul-animaetion.webm
var coolAnimationVideo []byte

//go:embed videos/9/*
var cirnoDayVideoFiles embed.FS

var cirnoDayVideos = loadVideos(cirnoDayVideoFiles, "videos/9")

type Video struct {
	Name string
	Data []byte
}

func GetDefaultVideo() Video {
	return Video{Name: "pls-rember.mp4", Data: defaultVideo}
}

func SemberVideo() Video {
	return Video{Name: "sember.mp4", Data: semberVideo}
}

func EvilVideo() Video {
	return Video{Name: "evil.webm", Data: evilVideo}
}

func CoolAnimation() Video {
	return Video{Name: "cul-animaetion.webm", Data: coolAnimationVideo}
}

func CirnoDayVideo() Video {
	return cirnoDayVideos[rand.IntN(len(cirnoDayVideos))]
}

func loadVideos(files embed.FS, directory string) []Video {
	entries, err := files.ReadDir(directory)
	if err != nil {
		panic(fmt.Sprintf("read embedded videos in %s: %v", directory, err))
	}

	videos := make([]Video, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := directory + "/" + entry.Name()
		data, err := files.ReadFile(path)
		if err != nil {
			panic(fmt.Sprintf("read embedded video %s: %v", path, err))
		}

		videos = append(videos, Video{Name: entry.Name(), Data: data})
	}

	if len(videos) == 0 {
		panic(fmt.Sprintf("no embedded videos in %s", directory))
	}

	return videos
}
