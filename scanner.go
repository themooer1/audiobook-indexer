package scanner

import (
	"fmt"
	"io/fs"
	"sync"

	library "github.com/themooer1/audiobook-library"
)

func fileScanner(filesToScan <-chan string, chaptersOut chan<- RelativeAudioBookChapter, errorHandler func(path string, err error), wg *sync.WaitGroup) {
	for file := range filesToScan {
		chapter, err := fromFile(file)

		if err != nil {
			errorHandler(file, err)
		} else {
			chaptersOut <- chapter
		}
	}

	wg.Done()
}

func startFileScanners(scanners int, filesToScan <-chan string, chaptersOut chan<- RelativeAudioBookChapter, errorHandler func(path string, err error)) {
	var wg sync.WaitGroup
	wg.Add(scanners)

	for i := 0; i < scanners; i++ {
		go fileScanner(filesToScan, chaptersOut, errorHandler, &wg)
	}

	wg.Wait()
	close(chaptersOut)
}

func importIntoUnsortedLibrary(library *UnsortedBookLibrary, chapters <-chan RelativeAudioBookChapter) {
	for c := range chapters {
		library.AddChapter(c)
	}
}

// Scan scans the given directory for audio files, uses the sorter to organize them into audiobooks
// and adds them to the given library
func Scan(rootDir string, library *library.AudioBookLibrary, sorter Sorter[RelativeAudioBookChapter]) []error {
	audioFilesToScan := make(chan string, 100)
	chapters := make(chan RelativeAudioBookChapter)
	var unsortedLibrary UnsortedBookLibrary
	unsortedLibrary.Initialize()

	filter := func(path string, d fs.DirEntry) (bool, error) {
		info, err := d.Info()
		if err != nil {
			return false, err
		}

		return isSupportedAudioFile(info), nil
	}

	onError := func(path string, d fs.DirEntry, err error) error {
		fmt.Errorf("Failed to enumerate: %s, %s", path, err)

		// Scanning errors should be non-fatal
		return nil
	}

	onFile := func(path string, d fs.DirEntry) {
		audioFilesToScan <- path
	}

	onScanError := func(path string, err error) {
		fmt.Errorf("Failed to scan: %s, %s", path, err)
	}

	go func() {
		Walk(rootDir, filter, onFile, onError)
		close(audioFilesToScan)
	}()
	go startFileScanners(7, audioFilesToScan, chapters, onScanError)
	importIntoUnsortedLibrary(&unsortedLibrary, chapters)

	return unsortedLibrary.AddAllToAudioBookLibrary(library, sorter)
}

// ScanToNewLibrary scans the given directory for audio files, uses the sorter to organize them into audiobooks
// and returns a new library containing the audiobooks
func ScanToNewLibrary(rootDir string, sorter Sorter[RelativeAudioBookChapter]) (*library.AudioBookLibrary, []error) {
	lib := library.AudioBookLibrary{}
	lib.Initialize()

	errors := Scan(rootDir, &lib, sorter)

	return &lib, errors
}
