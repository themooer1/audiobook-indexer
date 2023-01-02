package scanner

import (
	"testing"
)

type MockChapterWithFilePath struct {
	ord      int
	filePath string
}

func (m *MockChapterWithFilePath) FilePath() string {
	return m.filePath
}

func TestSortByFilename(t *testing.T) {
	type args struct {
		items []MockChapterWithFilePath
	}
	tests := []struct {
		name string
		args args
		// SortByFilename never returns errors
	}{
		{
			"Sorted List",
			args{
				items: []MockChapterWithFilePath{
					{
						ord:      0,
						filePath: "filea.wav",
					},
					{
						ord:      1,
						filePath: "fileb.wav",
					},
					{
						ord:      2,
						filePath: "filec.wav",
					},
				},
			},
		},
		{
			"Unsorted List",
			args{
				items: []MockChapterWithFilePath{
					{
						ord:      1,
						filePath: "fileb.wav",
					},
					{
						ord:      0,
						filePath: "filea.wav",
					},
					{
						ord:      2,
						filePath: "filec.wav",
					},
				},
			},
		},
		{
			"Long Paths",
			args{
				items: []MockChapterWithFilePath{
					{
						ord:      2,
						filePath: "/x/users/Mario/filec.wav",
					},
					{
						ord:      1,
						filePath: "/y/users/Mario/fileb.wav",
					},
					{
						ord:      0,
						filePath: "/z/users/Mario/filea.wav",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sorted, errors := SortByFilename(tt.args.items)

			// SortByFilename should handle all inputs
			if len(errors) != 0 {
				t.Error("SortByFilename should never return errors, but it did: ", errors)
			}

			for index, chapter := range sorted {
				if index != chapter.ord {
					t.Fatalf("Expected %v in position %d, but it was in position %d.", chapter, index, chapter.ord)
				}
			}
		})
	}
}
