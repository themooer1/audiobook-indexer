package scanner

import (
	"path"
	"sort"
)

type HasFilePath interface {
	FilePath() string
}

func SortByFilename[T any, TPtr interface {
	HasFilePath
	*T
}](items []T) ([]T, []error) {
	sort.Slice(
		items,
		func(i, j int) bool {
			return path.Base(TPtr(&items[i]).FilePath()) < path.Base(TPtr(&items[j]).FilePath())
		},
	)

	return items, []error{}
}
