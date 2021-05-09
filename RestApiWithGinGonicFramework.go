package main

import (
	"math/rand"
	"strconv"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

type User struct {
	ID   string `json : "id"`
	Name string `json : "name"`
	Age  int    `json : "age"`
}

var Users []User

// var Users2 = []User{}

func main() {
	server := gin.Default()
	userRoutes := server.Group("/users")
	{
		userRoutes.GET("/", getUser)
		userRoutes.POST("/", Createuser)
		userRoutes.PUT("/:id", editUser)
		userRoutes.DELETE("/:id", deleteUser)

	}
	server.Run(":8080")
}

func getUser(ctx *gin.Context) {
	/*ctx.JSON(200, gin.H{
		"message": "Hello Everyone",
	})
	*/
	ctx.JSON(200, Users)
}

func Createuser(ctx *gin.Context) {
	var reqBody User
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(422, gin.H{
			"err":     "True",
			"Message": "Invalid Message Body",
		})
		return
	}
	//	reqBody.ID = uuid.New().String()   // Can also use UUID in insted of using Rand()
	reqBody.ID = strconv.Itoa(rand.Intn(1000))
	Users = append(Users, reqBody)

	ctx.JSON(200, gin.H{
		"err":     "False",
		"Message": "Post Successful",
	})
}

func editUser(cxt *gin.Context) {
	id := cxt.Param("id")
	var reqBody User
	if err := cxt.ShouldBindJSON(&reqBody); err != nil {
		cxt.JSON(422, gin.H{
			"err":     "True",
			"Message": "Invalid Message Body",
		})
		return
	}
	for i, u := range Users {
		if u.ID == id {
			Users[i].Name = reqBody.Name
			Users[i].Age = reqBody.Age

			cxt.JSON(200, gin.H{
				"err":     "False",
				"Message": "Post is Edited",
			})
			return
		}
	}
	cxt.JSON(404, gin.H{
		"err":     "True",
		"Message": "Invalid ID",
	})
}
func deleteUser(cxt *gin.Context) {
	id := cxt.Param("id")
	for i, u := range Users {
		if u.ID == id {
			Users = append(Users[:i], Users[i+1:]...)

			cxt.JSON(200, gin.H{
				"err":     "False",
				"Message": "User Deleted",
			})
			return
		}
	}
	cxt.JSON(404, gin.H{
		"err":     "True",
		"Message": "Invalid ID",
	})

}
