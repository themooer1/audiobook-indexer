package scanner

import (
	"testing"
)

type MockChapterWithDiscNumber struct {
	ord      int
	discNum  DiscNumber
	trackNum TrackNumber
}

func (m *MockChapterWithDiscNumber) DiscNum() DiscNumber {
	return m.discNum
}

func (m *MockChapterWithDiscNumber) TrackNum() TrackNumber {
	return m.trackNum
}

func TestSortByDiscNum(t *testing.T) {
	type args struct {
		items []MockChapterWithDiscNumber
	}
	tests := []struct {
		name string
		args args
		// SortByFilename never returns errors
	}{
		{
			"One Track",
			args{
				items: []MockChapterWithDiscNumber{
					{
						ord:      0,
						discNum:  1,
						trackNum: 1,
					},
				},
			},
		},
		{
			"Sort by Disc Only",
			args{
				items: []MockChapterWithDiscNumber{
					{
						ord:      2,
						discNum:  2,
						trackNum: 0,
					},
					{
						ord:      0,
						discNum:  0,
						trackNum: 0,
					},
					{
						ord:      1,
						discNum:  1,
						trackNum: 0,
					},
				},
			},
		},
		{
			"Sort by Track Only",
			args{
				items: []MockChapterWithDiscNumber{
					{
						ord:      2,
						discNum:  0,
						trackNum: 2,
					},
					{
						ord:      0,
						discNum:  0,
						trackNum: 0,
					},
					{
						ord:      1,
						discNum:  0,
						trackNum: 1,
					},
				},
			},
		},
		{
			"Sort by Disc and Track",
			args{
				items: []MockChapterWithDiscNumber{
					{
						ord:      0,
						discNum:  0,
						trackNum: 0,
					},
					{
						ord:      1,
						discNum:  0,
						trackNum: 1,
					},
					{
						ord:      2,
						discNum:  0,
						trackNum: 2,
					},
					{
						ord:      3,
						discNum:  1,
						trackNum: 0,
					},
					{
						ord:      4,
						discNum:  1,
						trackNum: 1,
					},
					{
						ord:      5,
						discNum:  1,
						trackNum: 2,
					},
					{
						ord:      6,
						discNum:  2,
						trackNum: 0,
					},
					{
						ord:      7,
						discNum:  2,
						trackNum: 1,
					},
				},
			},
		},
		{
			"Sort with Missing Tracks",
			args{
				items: []MockChapterWithDiscNumber{
					{
						ord:      0,
						discNum:  0,
						trackNum: 0,
					},
					{
						ord:      1,
						discNum:  0,
						trackNum: 1,
					},
					{
						ord:      2,
						discNum:  0,
						trackNum: 2,
					},
					{
						ord:      3,
						discNum:  1,
						trackNum: 2,
					},
					{
						ord:      4,
						discNum:  2,
						trackNum: 1,
					},
				},
			},
		},
		{
			"Sort with Missing Discs",
			args{
				items: []MockChapterWithDiscNumber{
					{
						ord:      0,
						discNum:  0,
						trackNum: 0,
					},
					{
						ord:      1,
						discNum:  0,
						trackNum: 1,
					},
					{
						ord:      2,
						discNum:  3,
						trackNum: 1,
					},
					{
						ord:      3,
						discNum:  3,
						trackNum: 2,
					},
					{
						ord:      4,
						discNum:  3,
						trackNum: 3,
					},
				},
			},
		},
		{
			"Sort with Negative Tracks",
			args{
				items: []MockChapterWithDiscNumber{
					{
						ord:      0,
						discNum:  0,
						trackNum: -3,
					},
					{
						ord:      1,
						discNum:  0,
						trackNum: -2,
					},
					{
						ord:      2,
						discNum:  0,
						trackNum: -1,
					},
					{
						ord:      3,
						discNum:  1,
						trackNum: 0,
					},
					{
						ord:      4,
						discNum:  1,
						trackNum: 1,
					},
					{
						ord:      5,
						discNum:  1,
						trackNum: 2,
					},
					{
						ord:      6,
						discNum:  2,
						trackNum: -4,
					},
					{
						ord:      7,
						discNum:  2,
						trackNum: -2,
					},
				},
			},
		},
		{
			"Sort with Negative Discs",
			args{
				items: []MockChapterWithDiscNumber{
					{
						ord:      0,
						discNum:  -3,
						trackNum: 0,
					},
					{
						ord:      1,
						discNum:  -3,
						trackNum: 1,
					},
					{
						ord:      2,
						discNum:  -3,
						trackNum: 2,
					},
					{
						ord:      3,
						discNum:  1,
						trackNum: 0,
					},
					{
						ord:      4,
						discNum:  1,
						trackNum: 1,
					},
					{
						ord:      5,
						discNum:  1,
						trackNum: 2,
					},
					{
						ord:      6,
						discNum:  2,
						trackNum: 0,
					},
					{
						ord:      7,
						discNum:  2,
						trackNum: 1,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sorted, errors := SortByDiscNumber(tt.args.items)

			// SortByFilename should handle all inputs
			if len(errors) != 0 {
				t.Error("SortByFilename should never return errors, but it did: ", errors)
			}

			for index, chapter := range sorted {
				if index != chapter.ord {
					t.Error(sorted)
					t.Fatalf("Expected %v in position %d, but it was in position %d.", chapter, chapter.ord, index)
				}
			}
		})
	}
}
