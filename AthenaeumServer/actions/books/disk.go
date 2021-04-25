package books

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/model"
	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/pkg/paths"
)

func organiseFile(book *model.Book, libraryRoot string) (*model.Book, bool, error) {
	currentRelativeLocation := book.File.FileLocation
	var expectedLocation string
	hasMoved := false
	if book.Series == nil {
		expectedLocation = fmt.Sprintf("/%s/%s%s",
			paths.CleanNameOnly(book.AuthorString()),
			paths.CleanNameOnly(book.Title),
			book.File.FileExtension,
		)
	} else {
		expectedLocation = fmt.Sprintf("/%s/%s/%s %s%s",
			paths.CleanNameOnly(book.AuthorString()),
			paths.CleanNameOnly(book.Series.Title),
			paths.CleanNameOnly(book.Series.EntryString()),
			paths.CleanNameOnly(book.Title),
			book.File.FileExtension,
		)
	}

	if currentRelativeLocation != expectedLocation {
		currentFullPath := filepath.Join(libraryRoot, currentRelativeLocation)
		targetFullPath := filepath.Join(libraryRoot, expectedLocation)
		err := os.MkdirAll(filepath.Dir(targetFullPath), os.ModePerm)
		if err != nil {
			return book, hasMoved, err
		}
		err = os.Rename(currentFullPath, targetFullPath)
		if err != nil {
			return book, hasMoved, err
		}
		book.File.FileLocation = expectedLocation
		hasMoved = true
	}
	return book, hasMoved, nil
}
