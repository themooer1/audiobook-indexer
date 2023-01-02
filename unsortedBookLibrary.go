package scanner

import library "github.com/themooer1/audiobook-library"

type UnsortedBookLibrary struct {
	books map[BookTitle]UnsortedBook
}

func (u *UnsortedBookLibrary) Initialize() {
	u.books = make(map[BookTitle]UnsortedBook)
}

func (u *UnsortedBookLibrary) AddChapter(chapter RelativeAudioBookChapter) {
	bookTitle := chapter.bookTitle

	var b UnsortedBook
	var ok bool
	if b, ok = u.books[bookTitle]; !ok {
		b.Initialize(bookTitle)
	}

	b.AddChapter(chapter)
	u.books[bookTitle] = b
}

func (u *UnsortedBookLibrary) AddAllToAudioBookLibrary(library *library.AudioBookLibrary, sorter Sorter[RelativeAudioBookChapter]) []error {
	var errors []error

	// Add books to library
	for _, unsortedBook := range u.books {
		// Sort each book
		book, err := unsortedBook.IntoAudioBook(sorter)

		// exit if any book fails
		if err != nil {
			errors = append(errors, err...)
		}

		// Add sorted book to library
		library.Add(*book)
	}

	return errors
}
