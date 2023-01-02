package scanner

import (
	"io/fs"
	"strings"
)

var audioFileExtensions = [...]string{".mp3", ".flac"}

func hasSupportedAudioFileExtension(fileName string) bool {
	for _, e := range audioFileExtensions {
		if strings.HasSuffix(fileName, e) {
			return true
		}
	}

	return false
}

func isSupportedAudioFile(fileInfo fs.FileInfo) bool {
	return !fileInfo.IsDir() && hasSupportedAudioFileExtension(fileInfo.Name())
}
