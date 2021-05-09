package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Saurav@123"
	dbname   = "postgres"
)

type Persons struct {
	Person_id   int    `json:"person_id"`
	Person_name string `json:"person_name"`
	Person_city string `json :"person_city"`
}

func main() {
	server := gin.Default()
	server.GET("/users", getUsers)
	server.POST("/users", createPerson)
	server.GET("/users/:pid", getUser)
	server.PUT("/users/:pid", updateUser)
	server.DELETE("/users/:pid", deleteUser)

	server.Run(":8080")
}

func OpenConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	//defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println(psqlInfo)
	fmt.Println("Successfully connected!")
	return db
}

func getUsers(c *gin.Context) {
	db := OpenConnection()
	c.Writer.Header().Set("Content-Type", "application/json")

	var persons []Persons
	result, err := db.Query("SELECT person_id, person_name, person_city from persons")
	if err != nil {
		panic(err.Error())
	}
	//defer result.Close()
	for result.Next() {
		var person Persons
		err := result.Scan(&person.Person_id, &person.Person_name, &person.Person_city)
		if err != nil {
			panic(err.Error())
		}
		persons = append(persons, person)
	}
	json.NewEncoder(c.Writer).Encode(persons)
	//	fmt.Fprintf(c.Writer, "All Data Displayed")
	//defer db.Close()
}

func createPerson(c *gin.Context) {
	db := OpenConnection()
	fmt.Println("CHECKING")
	c.Writer.Header().Set("Content-Type", "application/json")
	stmt, err := db.Prepare("INSERT INTO persons VALUES($1,$2,$3)")
	fmt.Println("CHECKING 2")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("CHECKING 3")
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)

	person_id := keyVal["person_id"]
	person_name := keyVal["person_name"]
	person_city := keyVal["person_city"]

	_, err = stmt.Exec(person_id, person_name, person_city)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(c.Writer, "New post was created")
}

func getUser(c *gin.Context) {
	db := OpenConnection()
	c.Writer.Header().Set("Content-Type", "application/json")
	p := c.Param("pid")
	fmt.Println("VAl of p", p)

	fmt.Println("CHECKING 4")
	result, err := db.Query("SELECT person_id, person_name, person_city FROM persons WHERE person_id = $1;", p)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("CHECKING 5")
	//defer result.Close()
	var person Persons
	for result.Next() {
		err := result.Scan(&person.Person_id, &person.Person_name, &person.Person_city)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(c.Writer).Encode(person)
}

func updateUser(c *gin.Context) {
	db := OpenConnection()
	c.Writer.Header().Set("Content-Type", "application/json")
	p := c.Param("pid") // Param returns the the value of the key i.e what we pass in the URL
	//params := mux.Vars(r)   // In mux Vars() will return map[key:value]
	stmt, err := db.Prepare("UPDATE persons SET person_name = $1, person_city = $2 WHERE person_id = $3")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	Person_name := keyVal["person_name"]
	Person_city := keyVal["person_city"]

	_, err = stmt.Exec(Person_name, Person_city, p)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(c.Writer, "Post with ID = %s was updated", p)
}

func deleteUser(c *gin.Context) {
	db := OpenConnection()
	c.Writer.Header().Set("Content-Type", "application/json")
	p := c.Param("pid")
	//params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM persons WHERE person_id = $1")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(p)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(c.Writer, "Post with ID = %s was deleted", p)
}
