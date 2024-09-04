package dataModels

import (
	"github.com/pkg/errors"
	"time"
)

type ValidatorData interface {
	Validate() error
}

type Book struct {
	ID          uint      `gorm:"primaryKey" json:"id,omitempty"` // `omitempty` уберет поле из JSON, если оно пустое
	Title       string    `gorm:"type:varchar(500);not null" json:"title"`
	ISBN        string    `gorm:"type:varchar; not null; unique" json:"isbn"`
	Authors     []Author  `gorm:"many2many:book_authors" json:"authors"`
	Description string    `gorm:"type:varchar(2000)" json:"description,omitempty"`
	Language    string    `gorm:"type:varchar; not null" json:"language"`
	Year        int16     `gorm:"type:smallint; not null" json:"year"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
}

func (b Book) Validate() error {
	errMsg := "Please fill in these variables: "
	if b.Title == "" {
		errMsg += " title"
	}
	if b.ISBN == "" {
		errMsg += " isbn"
	}
	if b.Language == "" {
		errMsg += " language"
	}
	if b.Year == 0 {
		errMsg += " year"
	}
	if errMsg != "Please fill in these variables: " {
		return errors.New(errMsg)
	}
	return nil
}

type BookResponse struct {
	ID          uint             `json:"id"`
	Title       string           `json:"title"`
	ISBN        string           `json:"isbn"`
	Authors     []AuthorResponse `json:"authors"`
	Description string           `json:"description,omitempty"`
	Language    string           `json:"language"`
	Year        int16            `json:"year"`
}

type Author struct {
	ID        uint      `gorm:"primaryKey" json:"id,omitempty"` // `omitempty` уберет поле из JSON, если оно пустое
	FirstName string    `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName  string    `gorm:"type:varchar(100);not null" json:"last_name"`
	Books     []Book    `gorm:"many2many:book_authors" json:"books,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
}

func (a Author) Validate() error {
	errMsg := "Please fill in these variables: "
	if a.FirstName == "" {
		errMsg += " first_name"
	}
	if a.LastName == "" {
		errMsg += " last_name"
	}
	if errMsg != "Please fill in these variables: " {
		return errors.New(errMsg)
	}
	return nil
}

type AuthorResponse struct {
	ID        uint   `gorm:"primaryKey" json:"id,omitempty"` // `omitempty` уберет поле из JSON, если оно пустое
	FirstName string `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName  string `gorm:"type:varchar(100);not null" json:"last_name"`
}
