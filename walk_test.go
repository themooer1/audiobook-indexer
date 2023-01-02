package scanner

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

type FilePath = string

type walkFile struct {
	path FilePath
}

type walkFolder struct {
	path FilePath
}

type walkError struct {
	path FilePath
}

type walkTestResult struct {
	root    FilePath
	files   map[FilePath]walkFile
	folders map[FilePath]walkFolder
	errors  map[FilePath]walkError
}

func (r *walkTestResult) Initialize(root FilePath) {
	r.root = root
	r.files = make(map[FilePath]walkFile)
	r.folders = make(map[FilePath]walkFolder)
	r.errors = make(map[FilePath]walkError)
}

func (r *walkTestResult) addFile(path FilePath) {
	fullpath := filepath.Join(r.root, path)
	r.files[fullpath] = walkFile{fullpath}
}

func (r *walkTestResult) addFolder(path FilePath) {
	fullpath := filepath.Join(r.root, path)
	r.folders[fullpath] = walkFolder{fullpath}

	// Git doesn't track empty folders, so create them
	os.MkdirAll(fullpath, 0755)
}

func (r *walkTestResult) addError(path FilePath) {
	fullpath := filepath.Join(r.root, path)
	r.errors[fullpath] = walkError{fullpath}
}

func (r *walkTestResult) getFileHandler(t *testing.T) WalkFileHandler {
	return func(path string, dirEntry fs.DirEntry) {
		if dirEntry.IsDir() {
			if _, ok := r.folders[path]; ok {
				delete(r.folders, path)
			} else {
				t.Errorf("Unexpected folder: %v", path)
			}
		} else {
			if _, ok := r.files[path]; ok {
				delete(r.files, path)
			} else {
				t.Errorf("Unexpected file: %v", path)
			}
		}
	}
}

func (r *walkTestResult) getErrorHandler(t *testing.T) WalkErrorHandler {
	return func(path string, dirEntry fs.DirEntry, err error) error {
		if _, ok := r.errors[path]; ok {
			delete(r.errors, path)
		} else {
			t.Errorf("Unexpected error: %v", path)
		}
		return nil
	}
}

func (r *walkTestResult) isSuccessful() bool {
	return len(r.files) == 0 && len(r.folders) == 0 && len(r.errors) == 0
}

func printFilter(path string, dirEntry fs.DirEntry) (bool, error) {
	fmt.Println("Filter: ", path, ", ", dirEntry)

	return true, nil
}

func printErrorHandler(path string, dirEntry fs.DirEntry, err error) error {
	fmt.Println("Error: ", path, ", ", dirEntry)

	return nil
}

func printFileHandler(path string, dirEntry fs.DirEntry) {
	fmt.Println("File: ", path, ", ", dirEntry)
}

func TestWalk(t *testing.T) {
	// Empty Test
	emptyTestResult := walkTestResult{}
	emptyTestResult.Initialize("testdata/TestWalk/empty/root")
	emptyTestResult.addFolder("")

	// Empty Folder Test
	emptyFolderTestResult := walkTestResult{}
	emptyFolderTestResult.Initialize("testdata/TestWalk/emptyFolder/root")
	emptyFolderTestResult.addFolder(".")
	emptyFolderTestResult.addFolder("folder")
	emptyFolderTestResult.addFolder("folder/subfolder")

	// Flat Test
	flatTestResult := walkTestResult{}
	flatTestResult.Initialize("testdata/TestWalk/flat/root")
	flatTestResult.addFolder(".")
	flatTestResult.addFile("file1")
	flatTestResult.addFile("file2")

	// Folder in Folder Test
	folderInFolderTestResult := walkTestResult{}
	folderInFolderTestResult.Initialize("testdata/TestWalk/folderInFolder/root")
	folderInFolderTestResult.addFolder(".")
	folderInFolderTestResult.addFolder("folder")
	folderInFolderTestResult.addFolder("folder/subfolder")
	folderInFolderTestResult.addFile("folder/subfolder/file")

	// Tree Test
	treeTestResult := walkTestResult{}
	treeTestResult.Initialize("testdata/TestWalk/tree/root")
	treeTestResult.addFolder(".")
	treeTestResult.addFile("rootfile")
	treeTestResult.addFolder("folder")
	treeTestResult.addFile("folder/folderfile")
	treeTestResult.addFolder("folder/subfolder")
	treeTestResult.addFile("folder/subfolder/subfolderfile1")
	treeTestResult.addFile("folder/subfolder/subfolderfile2")

	type args struct {
		rootDir      string
		filter       WalkFilter
		fileHandler  WalkFileHandler
		errorHandler WalkErrorHandler
	}
	tests := []struct {
		name string
		args args
	}{
		// {
		// 	"PrintTr2ee",
		// 	args{
		// 		"./testdata/TestWalk/tree/root",
		// 		printFilter,
		// 		printFileHandler,
		// 		printErrorHandler,
		// 	},
		// },
		{
			"Empty Test",
			args{
				"./testdata/TestWalk/empty/root",
				printFilter,
				emptyTestResult.getFileHandler(t),
				emptyTestResult.getErrorHandler(t),
			},
		},
		{
			"Empty Folder Test",
			args{
				"./testdata/TestWalk/emptyFolder/root",
				printFilter,
				emptyFolderTestResult.getFileHandler(t),
				emptyFolderTestResult.getErrorHandler(t),
			},
		},
		{
			"Flat Test",
			args{
				"./testdata/TestWalk/flat/root",
				printFilter,
				flatTestResult.getFileHandler(t),
				flatTestResult.getErrorHandler(t),
			},
		},
		{
			"Folder in Folder Test",
			args{
				"./testdata/TestWalk/folderInFolder/root",
				printFilter,
				folderInFolderTestResult.getFileHandler(t),
				folderInFolderTestResult.getErrorHandler(t),
			},
		},
		{
			"Tree Test",
			args{
				"./testdata/TestWalk/tree/root",
				printFilter,
				treeTestResult.getFileHandler(t),
				treeTestResult.getErrorHandler(t),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Walk(tt.args.rootDir, tt.args.filter, tt.args.fileHandler, tt.args.errorHandler)
		})
	}
}
