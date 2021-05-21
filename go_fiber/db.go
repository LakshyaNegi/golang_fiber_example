package main

import (
	"log"

	"github.com/jinzhu/gorm"
)

var (
	DBConn *gorm.DB
)

type Shoe struct {
	gorm.Model
	Name  string  `json:"name,omitempty"`
	Size  int     `json:"size,omitempty"`
	Price float64 `json:"price,omitempty"`
	Email string  `json:"email,omitempty"`
}

type Users struct {
	gorm.Model
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

// var shoes = []*Shoe{
// 	{Id: 1, Name: "Air max 97", Size: 10, Price: 16499, Email: "admin"},
// 	{Id: 2, Name: "Air force 1", Size: 9, Price: 8999, Email: "admin"},
// 	{Id: 3, Name: "Cortez", Size: 8, Price: 7999, Email: "admin2"},
// }

func initDB() {
	//var err error
	DBConn, err := gorm.Open("sqlite3", "shoes.db")
	if err != nil {
		panic("Failed to connect database")
	}
	DBConn2, err := gorm.Open("sqlite3", "users.db")
	if err != nil {
		panic("Failed to connect database")
	}

	log.Printf("Database connected")

	DBConn.AutoMigrate(&Shoe{})
	DBConn2.AutoMigrate(&Users{})
	log.Printf("Database migrated")
}
