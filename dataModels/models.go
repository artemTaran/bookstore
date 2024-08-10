package dataModels

import "time"

type Book struct {
	ID          uint      `gorm:"primaryKey" json:"id,omitempty"` // `omitempty` уберет поле из JSON, если оно пустое
	Title       string    `gorm:"type:varchar(500);not null" json:"title"`
	ISBN        string    `gorm:"type:varchar(20);not null; unique" json:"isbn"`
	Authors     []Author  `gorm:"many2many:book_authors" json:"authors"`
	Description string    `gorm:"type:varchar(2000)" json:"description,omitempty"`
	Language    string    `gorm:"type:varchar(20); not null" json:"language"`
	Year        int16     `gorm:"type:smallint; not null" json:"year"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
}

type Author struct {
	ID        uint      `gorm:"primaryKey" json:"id,omitempty"` // `omitempty` уберет поле из JSON, если оно пустое
	FirstName string    `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName  string    `gorm:"type:varchar(100);not null" json:"last_name"`
	Books     []Book    `gorm:"many2many:book_authors" json:"books,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
}
