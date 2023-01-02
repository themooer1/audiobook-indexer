package scanner

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/go-test/deep"
	library "github.com/themooer1/audiobook-library"
)

func Test_fileScanner(t *testing.T) {
	type args struct {
		filesToScan      []string
		expectedChapters map[RelativeAudioBookChapter]struct{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"No Files",
			args{
				filesToScan: []string{
					// No files to scan
				},
				expectedChapters: map[RelativeAudioBookChapter]struct{}{
					// No chapters
				},
			},
		},
		{
			"Single File",
			args{
				filesToScan: []string{"./testdata/audiobooks/frankenstein/frankenstein_00_shelley_64kb.mp3"},
				expectedChapters: map[RelativeAudioBookChapter]struct{}{
					{
						title:      "00 - Letters",
						bookTitle:  "Frankenstein",
						bookAuthor: "Mary W. Shelley",
						discNum:    0,
						trackNum:   1,
						filePath:   "./testdata/audiobooks/frankenstein/frankenstein_00_shelley_64kb.mp3",
					}: {},
				},
			},
		},
		{
			"Two Files",
			args{
				filesToScan: []string{
					"./testdata/audiobooks/frankenstein/frankenstein_00_shelley_64kb.mp3",
					"./testdata/audiobooks/frankenstein/frankenstein_01_shelley_64kb.mp3",
				},
				expectedChapters: map[RelativeAudioBookChapter]struct{}{
					{
						title:      "00 - Letters",
						bookTitle:  "Frankenstein",
						bookAuthor: "Mary W. Shelley",
						discNum:    0,
						trackNum:   1,
						filePath:   "./testdata/audiobooks/frankenstein/frankenstein_00_shelley_64kb.mp3",
					}: {},
					{
						title:      "01 - Chapter 1",
						bookTitle:  "Frankenstein",
						bookAuthor: "Mary W. Shelley",
						discNum:    0,
						trackNum:   2,
						filePath:   "./testdata/audiobooks/frankenstein/frankenstein_01_shelley_64kb.mp3",
					}: {},
				},
			},
		},
		{
			"Three Files Two Books",
			args{
				filesToScan: []string{
					"./testdata/audiobooks/frankenstein/frankenstein_00_shelley_64kb.mp3",
					"./testdata/audiobooks/frankenstein/frankenstein_01_shelley_64kb.mp3",
					"./testdata/audiobooks/crimeandpunishment/crimepunishment_00_dostoyevsky_64kb.mp3",
				},
				expectedChapters: map[RelativeAudioBookChapter]struct{}{
					{
						title:      "00 - Letters",
						bookTitle:  "Frankenstein",
						bookAuthor: "Mary W. Shelley",
						discNum:    0,
						trackNum:   1,
						filePath:   "./testdata/audiobooks/frankenstein/frankenstein_00_shelley_64kb.mp3",
					}: {},
					{
						title:      "01 - Chapter 1",
						bookTitle:  "Frankenstein",
						bookAuthor: "Mary W. Shelley",
						discNum:    0,
						trackNum:   2,
						filePath:   "./testdata/audiobooks/frankenstein/frankenstein_01_shelley_64kb.mp3",
					}: {},
					{
						title:      "00 - Preface",
						bookTitle:  "Crime and Punishment (Version 3)",
						bookAuthor: "Fyodor Dostoyevsky",
						discNum:    0,
						trackNum:   1,
						filePath:   "./testdata/audiobooks/crimeandpunishment/crimepunishment_00_dostoyevsky_64kb.mp3",
					}: {},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			wg.Add(1)

			filesIn := make(chan string, len(tt.args.filesToScan))
			chaptersOut := make(chan RelativeAudioBookChapter, len(tt.args.filesToScan))
			errorHandler := func(path string, err error) {
				t.Errorf("Error scanning %s, %v", path, err)
			}

			for _, file := range tt.args.filesToScan {
				filesIn <- file
			}
			close(filesIn)

			go fileScanner(filesIn, chaptersOut, errorHandler, &wg)

			go func() {
				wg.Wait()
				close(chaptersOut)
			}()

			// Check that we got the expected chapters
			for c := range chaptersOut {
				if _, ok := tt.args.expectedChapters[c]; !ok {
					t.Errorf("Expected %v, got %v", tt.args.expectedChapters, c)
				}

				delete(tt.args.expectedChapters, c)
			}

			if len(tt.args.expectedChapters) > 0 {
				t.Errorf("Some chapters weren't found: %v", tt.args.expectedChapters)
			}
		})
	}
}

func TestScan(t *testing.T) {
	type args struct {
		rootDir string
		sorter  Sorter[RelativeAudioBookChapter]
	}
	tests := []struct {
		name string
		args args
		want library.AudioBookLibrary
	}{
		{
			"Frankenstein (Flat Directory)",
			args{
				rootDir: "./testdata/audiobooks/frankenstein",
				sorter:  SortByDiscNumber[RelativeAudioBookChapter],
			},
			library.AudioBookLibrary{
				AudioBooksByName: map[string]library.AudioBook{
					"Frankenstein": {
						Title:       "Frankenstein",
						Author:      "Mary W. Shelley",
						Description: "Description not available",
						Chapters: []library.AudioBookChapter{
							{
								Title: "00 - Letters",
								Index: 0,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_00_shelley_64kb.mp3",
							},

							{
								Title: "01 - Chapter 1",
								Index: 1,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_01_shelley_64kb.mp3",
							},

							{
								Title: "02 - Chapter 2",
								Index: 2,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_02_shelley_64kb.mp3",
							},

							{
								Title: "03 - Chapters 3-4",
								Index: 3,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_03_shelley_64kb.mp3",
							},

							{
								Title: "04 - Chapter 5",
								Index: 4,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_04_shelley_64kb.mp3",
							},

							{
								Title: "05 - Chapter 6",
								Index: 5,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_05_shelley_64kb.mp3",
							},

							{
								Title: "06 - Chapter 7",
								Index: 6,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_06_shelley_64kb.mp3",
							},

							{
								Title: "07 - Chapter 8",
								Index: 7,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_07_shelley_64kb.mp3",
							},

							{
								Title: "08 - Chapter 9",
								Index: 8,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_08_shelley_64kb.mp3",
							},

							{
								Title: "09 - Chapter 10",
								Index: 9,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_09_shelley_64kb.mp3",
							},

							{
								Title: "10 - Chapter 11",
								Index: 10,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_10_shelley_64kb.mp3",
							},

							{
								Title: "11 - Chapter 12",
								Index: 11,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_11_shelley_64kb.mp3",
							},

							{
								Title: "12 - Chapter 13",
								Index: 12,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_12_shelley_64kb.mp3",
							},

							{
								Title: "13 - Chapters 14-15",
								Index: 13,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_13_shelley_64kb.mp3",
							},

							{
								Title: "14 - Chapter 16",
								Index: 14,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_14_shelley_64kb.mp3",
							},

							{
								Title: "15 - Chapter 17",
								Index: 15,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_15_shelley_64kb.mp3",
							},

							{
								Title: "16 - Chapter 18",
								Index: 16,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_16_shelley_64kb.mp3",
							},

							{
								Title: "17 - Chapter 19",
								Index: 17,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_17_shelley_64kb.mp3",
							},

							{
								Title: "18 - Chapter 20",
								Index: 18,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_18_shelley_64kb.mp3",
							},

							{
								Title: "19 - Chapters 21-22",
								Index: 19,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_19_shelley_64kb.mp3",
							},

							{
								Title: "20 - Chapter 23",
								Index: 20,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_20_shelley_64kb.mp3",
							},

							{
								Title: "21 - Chapter 24",
								Index: 21,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_21_shelley_64kb.mp3",
							},
						},
					},
				},
			},
		},
		{
			"Crime and Punishment (Flat Directory)",
			args{
				rootDir: "./testdata/audiobooks/crimeandpunishment",
				sorter:  SortByDiscNumber[RelativeAudioBookChapter],
			},
			library.AudioBookLibrary{
				AudioBooksByName: map[string]library.AudioBook{
					"Crime and Punishment (Version 3)": {
						Title:       "Crime and Punishment (Version 3)",
						Author:      "Fyodor Dostoyevsky",
						Description: "Description not available",
						Chapters: []library.AudioBookChapter{
							{
								Title: "00 - Preface",
								Index: 0,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_00_dostoyevsky_64kb.mp3",
							},

							{
								Title: "01 - Part 1 Chapter 1",
								Index: 1,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_01_dostoyevsky_64kb.mp3",
							},

							{
								Title: "02 - Part 1 Chapter 2",
								Index: 2,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_02_dostoyevsky_64kb.mp3",
							},

							{
								Title: "03 - Part 1 Chapter 3",
								Index: 3,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_03_dostoyevsky_64kb.mp3",
							},

							{
								Title: "04 - Part 1 Chapter 4",
								Index: 4,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_04_dostoyevsky_64kb.mp3",
							},

							{
								Title: "05 - Part 1 Chapter 5",
								Index: 5,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_05_dostoyevsky_64kb.mp3",
							},

							{
								Title: "06 - Part 1 Chapter 6",
								Index: 6,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_06_dostoyevsky_64kb.mp3",
							},

							{
								Title: "07 - Part 1 Chapter 7",
								Index: 7,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_07_dostoyevsky_64kb.mp3",
							},

							{
								Title: "08 - Part 2 Chapter 1",
								Index: 8,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_08_dostoyevsky_64kb.mp3",
							},

							{
								Title: "09 - Part 2 Chapter 2",
								Index: 9,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_09_dostoyevsky_64kb.mp3",
							},

							{
								Title: "10 - Part 2 Chapter 3",
								Index: 10,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_10_dostoyevsky_64kb.mp3",
							},

							{
								Title: "11 - Part 2 Chapter 4",
								Index: 11,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_11_dostoyevsky_64kb.mp3",
							},

							{
								Title: "12 - Part 2 Chapter 5",
								Index: 12,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_12_dostoyevsky_64kb.mp3",
							},

							{
								Title: "13 - Part 2 Chapter 6",
								Index: 13,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_13_dostoyevsky_64kb.mp3",
							},

							{
								Title: "14 - Part 2 Chapter 7",
								Index: 14,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_14_dostoyevsky_64kb.mp3",
							},

							{
								Title: "15 - Part 3 Chapter 1",
								Index: 15,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_15_dostoyevsky_64kb.mp3",
							},

							{
								Title: "16 - Part 3 Chapter 2",
								Index: 16,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_16_dostoyevsky_64kb.mp3",
							},

							{
								Title: "17 - Part 3 Chapter 3",
								Index: 17,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_17_dostoyevsky_64kb.mp3",
							},

							{
								Title: "18 - Part 3 Chapter 4",
								Index: 18,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_18_dostoyevsky_64kb.mp3",
							},

							{
								Title: "19 - Part 3 Chapter 5",
								Index: 19,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_19_dostoyevsky_64kb.mp3",
							},

							{
								Title: "20 - Part 3 Chapter 6",
								Index: 20,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_20_dostoyevsky_64kb.mp3",
							},

							{
								Title: "21 - Part 4 Chapter 1",
								Index: 21,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_21_dostoyevsky_64kb.mp3",
							},

							{
								Title: "22 - Part 4 Chapter 2",
								Index: 22,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_22_dostoyevsky_64kb.mp3",
							},

							{
								Title: "23 - Part 4 Chapter 3",
								Index: 23,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_23_dostoyevsky_64kb.mp3",
							},

							{
								Title: "24 - Part 4 Chapter 4",
								Index: 24,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_24_dostoyevsky_64kb.mp3",
							},

							{
								Title: "25 - Part 4 Chapter 5",
								Index: 25,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_25_dostoyevsky_64kb.mp3",
							},

							{
								Title: "26 - Part 4 Chapter 6",
								Index: 26,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_26_dostoyevsky_64kb.mp3",
							},

							{
								Title: "27 - Part 5 Chapter 1",
								Index: 27,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_27_dostoyevsky_64kb.mp3",
							},

							{
								Title: "28 - Part 5 Chapter 2",
								Index: 28,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_28_dostoyevsky_64kb.mp3",
							},

							{
								Title: "29 - Part 5 Chapter 3",
								Index: 29,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_29_dostoyevsky_64kb.mp3",
							},

							{
								Title: "30 - Part 5 Chapter 4",
								Index: 30,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_30_dostoyevsky_64kb.mp3",
							},

							{
								Title: "31 - Part 5 Chapter 5",
								Index: 31,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_31_dostoyevsky_64kb.mp3",
							},

							{
								Title: "32 - Part 6 Chapter 1",
								Index: 32,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_32_dostoyevsky_64kb.mp3",
							},

							{
								Title: "33 - Part 6 Chapter 2",
								Index: 33,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_33_dostoyevsky_64kb.mp3",
							},

							{
								Title: "34 - Part 6 Chapter 3",
								Index: 34,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_34_dostoyevsky_64kb.mp3",
							},

							{
								Title: "35 - Part 6 Chapter 4",
								Index: 35,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_35_dostoyevsky_64kb.mp3",
							},

							{
								Title: "36 - Part 6 Chapter 5",
								Index: 36,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_36_dostoyevsky_64kb.mp3",
							},

							{
								Title: "37 - Part 6 Chapter 6",
								Index: 37,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_37_dostoyevsky_64kb.mp3",
							},

							{
								Title: "38 - Part 6 Chapter 7",
								Index: 38,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_38_dostoyevsky_64kb.mp3",
							},

							{
								Title: "39 - Part 6 Chapter 8",
								Index: 39,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_39_dostoyevsky_64kb.mp3",
							},

							{
								Title: "40 - Epilogue",
								Index: 40,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_40_dostoyevsky_64kb.mp3",
							},
						},
					},
				},
			},
		},
		{
			"Two Books (Full Tree)",
			args{
				rootDir: "./testdata/audiobooks",
				sorter:  SortByDiscNumber[RelativeAudioBookChapter],
			},
			library.AudioBookLibrary{
				AudioBooksByName: map[string]library.AudioBook{
					"Crime and Punishment (Version 3)": {
						Title:       "Crime and Punishment (Version 3)",
						Author:      "Fyodor Dostoyevsky",
						Description: "Description not available",
						Chapters: []library.AudioBookChapter{
							{
								Title: "00 - Preface",
								Index: 0,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_00_dostoyevsky_64kb.mp3",
							},

							{
								Title: "01 - Part 1 Chapter 1",
								Index: 1,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_01_dostoyevsky_64kb.mp3",
							},

							{
								Title: "02 - Part 1 Chapter 2",
								Index: 2,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_02_dostoyevsky_64kb.mp3",
							},

							{
								Title: "03 - Part 1 Chapter 3",
								Index: 3,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_03_dostoyevsky_64kb.mp3",
							},

							{
								Title: "04 - Part 1 Chapter 4",
								Index: 4,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_04_dostoyevsky_64kb.mp3",
							},

							{
								Title: "05 - Part 1 Chapter 5",
								Index: 5,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_05_dostoyevsky_64kb.mp3",
							},

							{
								Title: "06 - Part 1 Chapter 6",
								Index: 6,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_06_dostoyevsky_64kb.mp3",
							},

							{
								Title: "07 - Part 1 Chapter 7",
								Index: 7,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_07_dostoyevsky_64kb.mp3",
							},

							{
								Title: "08 - Part 2 Chapter 1",
								Index: 8,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_08_dostoyevsky_64kb.mp3",
							},

							{
								Title: "09 - Part 2 Chapter 2",
								Index: 9,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_09_dostoyevsky_64kb.mp3",
							},

							{
								Title: "10 - Part 2 Chapter 3",
								Index: 10,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_10_dostoyevsky_64kb.mp3",
							},

							{
								Title: "11 - Part 2 Chapter 4",
								Index: 11,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_11_dostoyevsky_64kb.mp3",
							},

							{
								Title: "12 - Part 2 Chapter 5",
								Index: 12,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_12_dostoyevsky_64kb.mp3",
							},

							{
								Title: "13 - Part 2 Chapter 6",
								Index: 13,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_13_dostoyevsky_64kb.mp3",
							},

							{
								Title: "14 - Part 2 Chapter 7",
								Index: 14,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_14_dostoyevsky_64kb.mp3",
							},

							{
								Title: "15 - Part 3 Chapter 1",
								Index: 15,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_15_dostoyevsky_64kb.mp3",
							},

							{
								Title: "16 - Part 3 Chapter 2",
								Index: 16,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_16_dostoyevsky_64kb.mp3",
							},

							{
								Title: "17 - Part 3 Chapter 3",
								Index: 17,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_17_dostoyevsky_64kb.mp3",
							},

							{
								Title: "18 - Part 3 Chapter 4",
								Index: 18,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_18_dostoyevsky_64kb.mp3",
							},

							{
								Title: "19 - Part 3 Chapter 5",
								Index: 19,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_19_dostoyevsky_64kb.mp3",
							},

							{
								Title: "20 - Part 3 Chapter 6",
								Index: 20,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_20_dostoyevsky_64kb.mp3",
							},

							{
								Title: "21 - Part 4 Chapter 1",
								Index: 21,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_21_dostoyevsky_64kb.mp3",
							},

							{
								Title: "22 - Part 4 Chapter 2",
								Index: 22,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_22_dostoyevsky_64kb.mp3",
							},

							{
								Title: "23 - Part 4 Chapter 3",
								Index: 23,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_23_dostoyevsky_64kb.mp3",
							},

							{
								Title: "24 - Part 4 Chapter 4",
								Index: 24,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_24_dostoyevsky_64kb.mp3",
							},

							{
								Title: "25 - Part 4 Chapter 5",
								Index: 25,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_25_dostoyevsky_64kb.mp3",
							},

							{
								Title: "26 - Part 4 Chapter 6",
								Index: 26,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_26_dostoyevsky_64kb.mp3",
							},

							{
								Title: "27 - Part 5 Chapter 1",
								Index: 27,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_27_dostoyevsky_64kb.mp3",
							},

							{
								Title: "28 - Part 5 Chapter 2",
								Index: 28,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_28_dostoyevsky_64kb.mp3",
							},

							{
								Title: "29 - Part 5 Chapter 3",
								Index: 29,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_29_dostoyevsky_64kb.mp3",
							},

							{
								Title: "30 - Part 5 Chapter 4",
								Index: 30,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_30_dostoyevsky_64kb.mp3",
							},

							{
								Title: "31 - Part 5 Chapter 5",
								Index: 31,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_31_dostoyevsky_64kb.mp3",
							},

							{
								Title: "32 - Part 6 Chapter 1",
								Index: 32,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_32_dostoyevsky_64kb.mp3",
							},

							{
								Title: "33 - Part 6 Chapter 2",
								Index: 33,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_33_dostoyevsky_64kb.mp3",
							},

							{
								Title: "34 - Part 6 Chapter 3",
								Index: 34,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_34_dostoyevsky_64kb.mp3",
							},

							{
								Title: "35 - Part 6 Chapter 4",
								Index: 35,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_35_dostoyevsky_64kb.mp3",
							},

							{
								Title: "36 - Part 6 Chapter 5",
								Index: 36,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_36_dostoyevsky_64kb.mp3",
							},

							{
								Title: "37 - Part 6 Chapter 6",
								Index: 37,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_37_dostoyevsky_64kb.mp3",
							},

							{
								Title: "38 - Part 6 Chapter 7",
								Index: 38,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_38_dostoyevsky_64kb.mp3",
							},

							{
								Title: "39 - Part 6 Chapter 8",
								Index: 39,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_39_dostoyevsky_64kb.mp3",
							},

							{
								Title: "40 - Epilogue",
								Index: 40,
								Url:   "testdata/audiobooks/crimeandpunishment/crimepunishment_40_dostoyevsky_64kb.mp3",
							},
						},
					},
					"Frankenstein": {
						Title:       "Frankenstein",
						Author:      "Mary W. Shelley",
						Description: "Description not available",
						Chapters: []library.AudioBookChapter{
							{
								Title: "00 - Letters",
								Index: 0,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_00_shelley_64kb.mp3",
							},

							{
								Title: "01 - Chapter 1",
								Index: 1,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_01_shelley_64kb.mp3",
							},

							{
								Title: "02 - Chapter 2",
								Index: 2,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_02_shelley_64kb.mp3",
							},

							{
								Title: "03 - Chapters 3-4",
								Index: 3,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_03_shelley_64kb.mp3",
							},

							{
								Title: "04 - Chapter 5",
								Index: 4,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_04_shelley_64kb.mp3",
							},

							{
								Title: "05 - Chapter 6",
								Index: 5,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_05_shelley_64kb.mp3",
							},

							{
								Title: "06 - Chapter 7",
								Index: 6,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_06_shelley_64kb.mp3",
							},

							{
								Title: "07 - Chapter 8",
								Index: 7,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_07_shelley_64kb.mp3",
							},

							{
								Title: "08 - Chapter 9",
								Index: 8,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_08_shelley_64kb.mp3",
							},

							{
								Title: "09 - Chapter 10",
								Index: 9,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_09_shelley_64kb.mp3",
							},

							{
								Title: "10 - Chapter 11",
								Index: 10,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_10_shelley_64kb.mp3",
							},

							{
								Title: "11 - Chapter 12",
								Index: 11,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_11_shelley_64kb.mp3",
							},

							{
								Title: "12 - Chapter 13",
								Index: 12,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_12_shelley_64kb.mp3",
							},

							{
								Title: "13 - Chapters 14-15",
								Index: 13,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_13_shelley_64kb.mp3",
							},

							{
								Title: "14 - Chapter 16",
								Index: 14,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_14_shelley_64kb.mp3",
							},

							{
								Title: "15 - Chapter 17",
								Index: 15,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_15_shelley_64kb.mp3",
							},

							{
								Title: "16 - Chapter 18",
								Index: 16,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_16_shelley_64kb.mp3",
							},

							{
								Title: "17 - Chapter 19",
								Index: 17,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_17_shelley_64kb.mp3",
							},

							{
								Title: "18 - Chapter 20",
								Index: 18,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_18_shelley_64kb.mp3",
							},

							{
								Title: "19 - Chapters 21-22",
								Index: 19,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_19_shelley_64kb.mp3",
							},

							{
								Title: "20 - Chapter 23",
								Index: 20,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_20_shelley_64kb.mp3",
							},

							{
								Title: "21 - Chapter 24",
								Index: 21,
								Url:   "testdata/audiobooks/frankenstein/frankenstein_21_shelley_64kb.mp3",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lib := library.AudioBookLibrary{}
			lib.Initialize()

			errors := Scan(tt.args.rootDir, &lib, tt.args.sorter)

			if len(errors) > 0 {
				t.Errorf("Scan() returned errors: %v", errors)
			}

			if len(lib.AudioBooksByName) != len(tt.want.AudioBooksByName) {
				t.Errorf("Expected to find %d books, but got %d books.", len(lib.AudioBooksByName), len(tt.want.AudioBooksByName))
			}

			if diff := deep.Equal(tt.want.AudioBooksByName, lib.AudioBooksByName); diff != nil {
				for _, d := range diff {
					t.Error(d)
					fmt.Fprintf(os.Stderr, "\n")
				}
			}
		})
	}
}
