package assets

import (
	_ "embed"
)

//go:embed videos/pls-rember.mp4
var defaultVideo []byte

//go:embed videos/sember.mp4
var semberVideo []byte

//go:embed videos/evil.webm
var evilVideo []byte

//go:embed videos/cul-animaetion.webm
var coolAnimationVideo []byte

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
