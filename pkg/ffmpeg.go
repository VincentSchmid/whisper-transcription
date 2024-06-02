package helpers

import (
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

func ConvertVideoToAudio(videoPath string, audioFilePath string, audioBitrate int) {
	ffmpeg_go.Input(videoPath).
		Output(audioFilePath, ffmpeg_go.KwArgs{"ac": 1, "audio_bitrate": audioBitrate}).
		OverWriteOutput().
		ErrorToStdOut().
		Run()
}
