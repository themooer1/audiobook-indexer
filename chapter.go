package scanner

import (
	"bytes"
	"os"

	"github.com/dhowden/tag"
	"github.com/edsrzf/mmap-go"
	library "github.com/themooer1/audiobook-library"
)

/*
Used temporarily for sorting audiobook chapters.  Once sorted
we can convert this into an AudioBookChapter, discarding
discNum and trackNum and just keeping the overall index in the
audiobook
*/
type RelativeAudioBookChapter struct {
	title      string
	bookTitle  string
	bookAuthor string
	discNum    int
	trackNum   int
	filePath   string
}

func (r *RelativeAudioBookChapter) Title() string {
	return r.title
}

func (r *RelativeAudioBookChapter) BookTitle() string {
	return r.bookTitle
}

func (r *RelativeAudioBookChapter) BookAuthor() string {
	return r.bookAuthor
}

func (r *RelativeAudioBookChapter) DiscNum() int {
	return r.discNum
}

func (r *RelativeAudioBookChapter) TrackNum() int {
	return r.trackNum
}

func (r *RelativeAudioBookChapter) FilePath() string {
	return r.filePath
}

func fromFile(audioFilePath string) (RelativeAudioBookChapter, error) {

	rawAudioFile, err := os.Open(audioFilePath)
	if err != nil {
		return RelativeAudioBookChapter{}, err
	} else {
		defer rawAudioFile.Close()
	}

	// MMAP the rawAudioFile to get better performance
	// when the tag parser makes lots of small reads.
	var audioFileMappedBytes mmap.MMap
	audioFileMappedBytes, err = mmap.Map(rawAudioFile, mmap.RDONLY, 0)
	if err != nil {
		return RelativeAudioBookChapter{}, err
	} else {
		defer audioFileMappedBytes.Unmap()
	}

	audioFileMmappedBytesReader := bytes.NewReader(audioFileMappedBytes)

	// bufferedAudioFile := bufio.NewReader(rawAudioFile)
	// audioFileBytes := rawAudioFile.read

	metadata, err := tag.ReadFrom(audioFileMmappedBytesReader)
	if err != nil {
		return RelativeAudioBookChapter{}, err
	}

	title := metadata.Title()
	bookTitle := metadata.Album()
	bookAuthor := metadata.Artist()
	discNum, _ := metadata.Disc()
	trackNum, _ := metadata.Track()

	return RelativeAudioBookChapter{title, bookTitle, bookAuthor, discNum, trackNum, audioFilePath}, nil
}

func (r *RelativeAudioBookChapter) intoAudioBookChapter(index int) library.AudioBookChapter {
	return library.AudioBookChapter{
		Title: r.title,
		Index: index,
		Url:   r.filePath,
	}
}
