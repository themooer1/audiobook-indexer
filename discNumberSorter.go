package scanner

import (
	"github.com/themooer1/gort/math"
)

type DiscNumber = int
type TrackNumber = int

type HasDiscNumber interface {
	DiscNum() DiscNumber
}

type HasTrackNumber interface {
	TrackNum() TrackNumber
}

type HasDiscAndTrackNumber interface {
	HasDiscNumber
	HasTrackNumber
}

type DiscNumberSorter[T any, TPtr interface {
	*T
	HasDiscAndTrackNumber
}] struct {
	chaptersByDiscByTrack map[DiscNumber]map[TrackNumber]T
	chapters              []T
	minDiscNumber         DiscNumber
	minTrackNumberByDisc  map[DiscNumber]TrackNumber
	maxDiscNumber         DiscNumber
	maxTrackNumberByDisc  map[DiscNumber]TrackNumber
}

func (s *DiscNumberSorter[T, TPtr]) Initalize(chapters []T) {
	s.chapters = chapters
	// Initial guess 10 tracks to a disc
	s.chaptersByDiscByTrack = make(map[DiscNumber]map[TrackNumber]T, len(s.chapters)/10)
	s.minTrackNumberByDisc = make(map[DiscNumber]TrackNumber)
	s.maxTrackNumberByDisc = make(map[DiscNumber]TrackNumber)

	for _, chapter := range s.chapters {
		discNum := TPtr(&chapter).DiscNum()
		trackNum := TPtr(&chapter).TrackNum()

		// Index chapter by disc and track number (creating maps if necessary)
		if _, ok := s.chaptersByDiscByTrack[discNum]; !ok {
			s.chaptersByDiscByTrack[discNum] = make(map[TrackNumber]T)
		}
		s.chaptersByDiscByTrack[discNum][trackNum] = chapter

		// Update min/max disc and track numbers
		s.minDiscNumber = math.MinInt(discNum, s.minDiscNumber)
		s.maxDiscNumber = math.MaxInt(discNum, s.maxDiscNumber)
		s.minTrackNumberByDisc[discNum] = math.MinInt(trackNum, s.minTrackNumberByDisc[discNum])
		s.maxTrackNumberByDisc[discNum] = math.MaxInt(trackNum, s.maxTrackNumberByDisc[discNum])
	}
}

func (s *DiscNumberSorter[T, TPtr]) sortedUnsafe() []T {
	sortedChapters := []T{}

	for discNum := s.minDiscNumber; discNum <= s.maxDiscNumber; discNum++ {
		if minTrackNumber, ok := s.minTrackNumberByDisc[discNum]; ok {
			if maxTrackNumber, ok := s.maxTrackNumberByDisc[discNum]; ok {
				for trackNum := minTrackNumber; trackNum <= maxTrackNumber; trackNum++ {
					if chapter, ok := s.chaptersByDiscByTrack[discNum][trackNum]; ok {
						sortedChapters = append(sortedChapters, chapter)
					}
				}
			}
		}
	}

	return sortedChapters
}

func (s *DiscNumberSorter[T, TPtr]) Sorted() ([]T, []error) {
	// if len(s.chaptersByDiscByTrack) < len(s.chapters) {
	// 	return nil, []error{ErrCannotSort}
	// } else {
	// 	return s.sortedUnsafe(), nil
	// }
	return s.sortedUnsafe(), nil
}

func SortByDiscNumber[T any, TPtr interface {
	*T
	HasDiscAndTrackNumber
}](items []T) ([]T, []error) {
	var s DiscNumberSorter[T, TPtr]

	s.Initalize(items)
	return s.Sorted()

}
