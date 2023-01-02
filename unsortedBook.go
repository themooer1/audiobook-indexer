// Temporary container to group chapters from the same book

package scanner

import (
	"errors"

	library "github.com/themooer1/audiobook-library"
)

type BookTitle = string

type UnsortedBook struct {
	BookTitle BookTitle
	Chapters  []RelativeAudioBookChapter
}

var ErrEmptyBook = errors.New("book has no chapters")

func (u *UnsortedBook) Initialize(bookTitle string) {
	u.BookTitle = bookTitle
	u.Chapters = make([]RelativeAudioBookChapter, 0, 10)
}

func (u *UnsortedBook) AddChapter(chapter RelativeAudioBookChapter) {
	u.Chapters = append(u.Chapters, chapter)
}

func (u *UnsortedBook) IntoAudioBook(sort Sorter[RelativeAudioBookChapter]) (*library.AudioBook, []error) {
	sortedRelChapters, errors := sort(u.Chapters)
	if errors != nil {
		return nil, errors
	}

	chapters := make([]library.AudioBookChapter, len(sortedRelChapters))
	for i, relChapter := range sortedRelChapters {
		chapters[i] = relChapter.intoAudioBookChapter(i)
	}

	if len(chapters) == 0 {
		return nil, []error{ErrEmptyBook}
	}

	bookTitle := u.Chapters[0].bookTitle
	bookAuthor := u.Chapters[0].bookAuthor
	bookDescription := "Description not available"

	return &library.AudioBook{
		Title:       bookTitle,
		Author:      bookAuthor,
		Description: bookDescription,
		Chapters:    chapters,
	}, nil
}
