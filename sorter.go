package scanner

import (
	"errors"
)

// type ChapterSorterStatus int32

var ErrCannotSort = errors.New("cannot sort chapters")

// Returns a sorted chapter list
// type ChapterSorter func(chapters []RelativeAudioBookChapter) ([]RelativeAudioBookChapter, error)

type Sorter[T any] func(items []T) ([]T, []error)

func Compose[T any](sorters ...Sorter[T]) Sorter[T] {
	return func(items []T) ([]T, []error) {
		allErrors := []error{}

		for _, sort := range sorters {
			result, errors := sort(items)
			if len(errors) == 0 {
				return result, nil
			}

			allErrors = append(allErrors, errors...)
		}

		return nil, allErrors
	}
}

// func Compose(sorters ...ChapterSorter) ChapterSorter {
// 	return func(chapters []RelativeAudioBookChapter) ([]RelativeAudioBookChapter, error) {
// 		for _, sort := range sorters {
// 			result, err := sort(chapters)
// 			if err != nil {
// 				return result, nil
// 			}

// 		}

// 	}
// }
