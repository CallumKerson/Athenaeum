package scribble

import (
	"encoding/json"

	"github.com/CallumKerrEdwards/Athenaeum/AthenaeumServer/model"
)

func (c *ScribbleStore) Add(book model.Book) error {
	c.logger.Debugf("Writing book to collection %s with key %s", bookCollection, book.Id)
	return c.Write(bookCollection, book.Id, book)
}

func (c *ScribbleStore) Remove(book model.Book) error {
	c.logger.Debugf("Removing book with key %s from collection", book.Id, bookCollection)
	return c.Delete(bookCollection, book.Id)
}

func (c *ScribbleStore) Get(id string) (*model.Book, error) {
	book := model.Book{}
	err := c.Read(bookCollection, id, &book)
	return &book, err
}

func (c *ScribbleStore) GetAll() ([]model.Book, error) {
	records, err := c.ReadAll(bookCollection)
	if err != nil {
		c.logger.Errorf("Cannot get all items from collection %s", bookCollection)
		return make([]model.Book, 0), err
	}
	foundBooks := []model.Book{}
	for _, record := range records {
		bookFound := model.Book{}
		if err := json.Unmarshal([]byte(record), &bookFound); err != nil {
			c.logger.Errorf("Cannot unmarshal %s to Book", string(record))
			return make([]model.Book, 0), err
		}
		foundBooks = append(foundBooks, bookFound)
	}
	return foundBooks, nil
}
