# audiobook-scanner
Fast tag scanner for audiobook libraries.

Inspects the IDv3, FLAC, MP4, and any other tags supported by the excellent library, https://github.com/dhowden/tag.

Usage:
```golang
lib := library.AudioBookLibrary{}
lib.Initialize()

sorter := scanner.SortByDiscNumber[scanner.RelativeAudioBookChapter]

errors := scanner.Scan(audioRoot, &lib, sorter)
```
