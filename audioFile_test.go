package scanner

import (
	"io/fs"
	"os"
	"testing"
)

func Test_hasSupportedAudioFileExtension(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"boga.mp3", args{"boga.mp3"}, true},
		{"boga.flac", args{"boga.flac"}, true},
		{".mp3", args{".mp3"}, true},
		{"banana.jpg", args{"banana.jpg"}, false},
		{"maurice.png", args{"maurice.png"}, false},
		{"maurice.pngmp3", args{"maurice.pngmp3"}, false},
		{"maurice.mp3png", args{"maurice.mp3png"}, false},
		{"maurice.mp3.png", args{"maurice.mp3.png"}, false},
		{"maurice.png.mp3", args{"maurice.png.mp3"}, true},
		{"maurice.mp3.flac", args{"maurice.mp3.flac"}, true},
		{"maurice.flac.mp3", args{"maurice.flac.mp3"}, true},
		{"maurice.flac.mp3.png", args{"maurice.flac.mp3.png"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasSupportedAudioFileExtension(tt.args.fileName); got != tt.want {
				t.Errorf("hasAudioFileExtension() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getFileInfo(t *testing.T, fileName string) fs.FileInfo {
	fileInfo, err := fs.Stat(os.DirFS("./testdata/Test_isSupportedAudioFile/root"), fileName)

	if err != nil {
		t.Error(err)
	}

	return fileInfo
}

func Test_isSupportedAudioFile(t *testing.T) {
	type args struct {
		fileInfo fs.FileInfo
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"MP3 File", args{getFileInfo(t, "preface.mp3")}, true},
		{"XML File", args{getFileInfo(t, "rss.xml")}, false},
		{"PNG File", args{getFileInfo(t, "screenshot.png")}, false},
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	} else {
		println(cwd)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSupportedAudioFile(tt.args.fileInfo); got != tt.want {
				t.Errorf("isAudioFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
