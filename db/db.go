package db

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"shop/config"
	"shop/dataModels"
	"shop/logger"
	"strings"
)

var db *gorm.DB

func InitDb() {
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal(err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", cfg.Host, cfg.UserName, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error(err)
	}
	if err = db.AutoMigrate(&dataModels.Book{}, &dataModels.Author{}); err != nil {
		logger.Error(err)
	}
}

func AddBook(newBook dataModels.Book) (err error) {
	for i, author := range newBook.Authors {
		var existingAuthor dataModels.Author
		if err = db.Where("first_name = ? AND last_name = ?", author.FirstName, author.LastName).First(&existingAuthor).Error; err == nil {
			newBook.Authors[i] = existingAuthor
		}
	}

	if err = db.Create(&newBook).Error; err != nil {
		if strings.Contains(err.Error(), "23505") {
			return fmt.Errorf("uniqueness Error: ISBN already exists")
		} else {
			return err
		}
	}
	return err
}

func GetBooks(quantity int) (books []dataModels.Book, err error) {
	if err = db.Preload("Authors").Limit(quantity).Find(&books).Error; err != nil {
		return nil, err
	}
	return books, err
}

func GetBookById(id int) (book dataModels.Book, err error) {
	err = db.Preload("Authors").Find(&book, id).Error
	return book, err
}

func SearchByTitle(title string) (books []dataModels.Book, err error) {
	err = db.Preload("Authors").Where("title ILIKE ?", "%"+title+"%").Find(&books).Error
	return books, err
}

func GetAuthors(quantity int) (authors []dataModels.Author, err error) {
	if err = db.Preload("Books").Limit(quantity).Find(&authors).Error; err != nil {
		return nil, err
	}
	return authors, err
}

func GetAuthorById(id int) (author dataModels.Author, err error) {
	err = db.Preload("Books").Find(&author, id).Error
	return author, err
}

func AddAuthor(newAuthor dataModels.Author) (msg string, err error) {
	err = db.Where("first_name = ? AND last_name = ?", newAuthor.FirstName, newAuthor.LastName).First(&dataModels.Author{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = db.Create(&newAuthor).Error; err != nil {
				return msg, err
			}
			msg = "author added successfully"
			return msg, nil
		} else {
			return msg, err
		}
	}
	msg = "author already exists"

	return msg, nil
}

func SearchByFLName(flName string) (authors []dataModels.Author, err error) {
	err = db.Preload("Books").Where("CONCAT(first_name, ' ',  last_name) ILIKE ?", "%"+flName+"%").
		Find(&authors).Error
	return authors, err
}
