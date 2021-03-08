package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Person struct {
	gorm.Model

	Name  string
	Email string `gorm:"typevarchar(100);unique_index"`
	Books []Book
}

type Book struct {
	gorm.Model

	Title      string
	Author     string
	CallNumber int `gorm:"unique_index"`
}

var (
	db  *gorm.DB
	err error
)

func main() {
	//creating a connection to the postgre database

	//loading environment
	dialect := os.Getenv("DIALECT")
	host := os.Getenv("HOST")
	dbPort := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbName := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")

	//establishing connection
	datasource := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbName, password, dbPort)

	//openning connection
	db, err = gorm.Open(dialect, datasource)
	if err != nil {
		log.Fatalln("log err", err)
	} else {
		fmt.Println("Connection establishment successfull")
	}

	//closing connection
	defer db.Close()

	//make migration to the database if the database have not been created
	db.AutoMigrate(&Person{})
	db.AutoMigrate(&Book{})

	router := mux.NewRouter()

	router.HandleFunc("/people", getpeople).Methods("GET")
	router.HandleFunc("/person/{id}", getperson).Methods("GET")
	router.HandleFunc("/delete/person/{id}", deleteperson).Methods("DELETE")
	router.HandleFunc("/create/person", createperson).Methods("POST")

	router.HandleFunc("/books", getbooks).Methods("GET")
	router.HandleFunc("/book/{id}", getbook).Methods("GET")
	router.HandleFunc("/delete/book/{id}", deletebook).Methods("DELETE")
	router.HandleFunc("/create/book", createbook).Methods("POST")

	log.Fatal(http.ListenAndServe(":8082", router))
}

func getpeople(w http.ResponseWriter, r *http.Request) {
	var people []Person
	db.Find(&people)

	json.NewEncoder(w).Encode(&people)
}

func getperson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var person Person
	var books Book
	db.First(&person, params["id"])
	db.Model(&person).Related(&books)

	json.NewEncoder(w).Encode(&person)
}

func createperson(w http.ResponseWriter, r *http.Request) {
	var person Person
	json.NewDecoder(r.Body).Decode(&person)

	createdPerson := db.Create(&person)
	err = createdPerson.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&person)
	}
}

func deleteperson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var person Person
	db.First(&person, params["id"])
	db.Delete(&person)
	json.NewEncoder(w).Encode(&person)
}

func getbooks(w http.ResponseWriter, r *http.Request) {
	var books []Book

	db.Find(&books)

	json.NewEncoder(w).Encode(&books)
}

func getbook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var books Book
	db.First(&books, params)
	json.NewEncoder(w).Encode(&books)

}

func deletebook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var book Book
	db.First(&book, params)
	db.Delete(&book)
	json.NewEncoder(w).Encode(&book)
}

func createbook(w http.ResponseWriter, r *http.Request) {
	var book Book
	json.NewDecoder(r.Body).Decode(&book)

	createdPerson := db.Create(&book)
	err = createdPerson.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&book)
	}

}
