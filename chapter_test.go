package scanner

import (
	"reflect"
	"testing"

	library "github.com/themooer1/audiobook-library"
)

type expectedRelativeAudioBookChapter struct {
	bookTitle string
	discNum   int
	trackNum  int
	filePath  string
}

func assertRelativeAudioBookChapter(t *testing.T, got RelativeAudioBookChapter, expected *expectedRelativeAudioBookChapter) {
	if got.bookTitle != expected.bookTitle {
		t.Errorf("bookTitle = %v, want %v", got.bookTitle, expected.bookTitle)
	}
	if got.discNum != expected.discNum {
		t.Errorf("discNum = %v, want %v", got.discNum, expected.discNum)
	}
	if got.trackNum != expected.trackNum {
		t.Errorf("trackNum = %v, want %v", got.trackNum, expected.trackNum)
	}
	if got.filePath != expected.filePath {
		t.Errorf("filePath = %v, want %v", got.filePath, expected.filePath)
	}
}

func readAudioBookChapterAndTest(t *testing.T, expected *expectedRelativeAudioBookChapter) {
	relativeAudioBookChapter, err := fromFile(expected.filePath)
	if err != nil {
		t.Errorf("error reading file %v, %v", expected.filePath, err)
	}
	assertRelativeAudioBookChapter(t, relativeAudioBookChapter, expected)
}

func Test_fromFile(t *testing.T) {
	tests := []struct {
		name    string
		want    expectedRelativeAudioBookChapter
		wantErr bool
	}{
		{
			name: "RelativeAudioBookChapter fromFile frankenstein_00_shelley_64kb.mp3",
			want: expectedRelativeAudioBookChapter{
				bookTitle: "Frankenstein",
				discNum:   0,
				trackNum:  1,
				filePath:  "testdata/audiobooks/frankenstein/frankenstein_00_shelley_64kb.mp3",
			},
			wantErr: false,
		},
		{
			name: "RelativeAudioBookChapter fromFile frankenstein_01_shelley_64kb.mp3",
			want: expectedRelativeAudioBookChapter{
				bookTitle: "Frankenstein",
				discNum:   0,
				trackNum:  2,
				filePath:  "testdata/audiobooks/frankenstein/frankenstein_01_shelley_64kb.mp3",
			},
			wantErr: false,
		},
		{
			name: "RelativeAudioBookChapter fromFile frankenstein_02_shelley_64kb.mp3",
			want: expectedRelativeAudioBookChapter{
				bookTitle: "Frankenstein",
				discNum:   0,
				trackNum:  3,
				filePath:  "testdata/audiobooks/frankenstein/frankenstein_02_shelley_64kb.mp3",
			},
			wantErr: false,
		},
		{
			name: "RelativeAudioBookChapter fromFile crimeandpunishment_00_dostoyevsky_64kb.mp3",
			want: expectedRelativeAudioBookChapter{
				bookTitle: "Crime and Punishment (Version 3)",
				discNum:   0,
				trackNum:  1,
				filePath:  "testdata/audiobooks/crimeandpunishment/crimepunishment_00_dostoyevsky_64kb.mp3",
			},
			wantErr: false,
		},
		{
			name: "RelativeAudioBookChapter fromFile crimeandpunishment_01_dostoyevsky_64kb.mp3",
			want: expectedRelativeAudioBookChapter{
				bookTitle: "Crime and Punishment (Version 3)",
				discNum:   0,
				trackNum:  2,
				filePath:  "testdata/audiobooks/crimeandpunishment/crimepunishment_01_dostoyevsky_64kb.mp3",
			},
			wantErr: false,
		},
		{
			name: "RelativeAudioBookChapter fromFile crimeandpunishment_02_dostoyevsky_64kb.mp3",
			want: expectedRelativeAudioBookChapter{
				bookTitle: "Crime and Punishment (Version 3)",
				discNum:   0,
				trackNum:  3,
				filePath:  "testdata/audiobooks/crimeandpunishment/crimepunishment_02_dostoyevsky_64kb.mp3",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			readAudioBookChapterAndTest(t, &tt.want)
		})
	}
}

func TestRelativeAudioBookChapter_intoAudioBookChapter(t *testing.T) {
	type fields struct {
		bookTitle string
		discNum   int
		trackNum  int
		filePath  string
	}
	type args struct {
		index int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   library.AudioBookChapter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RelativeAudioBookChapter{
				bookTitle: tt.fields.bookTitle,
				discNum:   tt.fields.discNum,
				trackNum:  tt.fields.trackNum,
				filePath:  tt.fields.filePath,
			}
			if got := r.intoAudioBookChapter(tt.args.index); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RelativeAudioBookChapter.intoAudioBookChapter() = %v, want %v", got, tt.want)
			}
		})
	}
}
