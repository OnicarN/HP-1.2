package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Task struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"CreatedAt"`
	CompletedAt time.Time `json:"CompletedAt"`
}

var tasks = []Task{}
var statusOptions = []string{"new", "ongoing", "completed"}

// creamos la primera funciónd de nuestro proyecto, la cual nos va a devolver todas las tareas
func getTasks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, tasks)
}

// con esta función vamos a poder filtrar las tareas por ID
func getTasksById(c *gin.Context) {
	id := c.Param("id")

	for _, t := range tasks {
		if t.Id == id {
			c.IndentedJSON(http.StatusOK, t)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
}

//ahora vamos a ir creando el post

func postTasks(c *gin.Context) {
	newTaskId := uuid.New().String()
	var newTask Task

	if err := c.BindJSON(&newTask); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "you can not request"})
		return
	}

	if newTask.Title == "" || newTask.Description == "" || newTask.Status == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Title, Description and Status should not be null",
		})
	}

	valid := false

	for _, status := range statusOptions {
		if newTask.Status == status {
			valid = true
			break
		}
	}

	if valid == false {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "This status is not valid",
		})
	}

	newTask.Id = newTaskId
	newTask.CreatedAt = time.Now()

	if newTask.Status == "completed" {
		newTask.CompletedAt = time.Now()
	} else {
		newTask.CompletedAt = time.Time{}
	}

	tasks = append(tasks, newTask)

	//vamos a devolverlo de la siguiente forma
	c.IndentedJSON(http.StatusCreated, tasks)
}

func main() {
	//vamos a ir creando las rutas

	//para inicializar las rutas crearé una variable llamada router
	router := gin.Default()

	//esto nos devuelve toda la data
	router.GET("/tasks", getTasks)

	//esto nos filtra por id
	router.GET("tasks/:id", getTasksById)

	//esto nos permite añadir tareas
	router.POST("/tasks", postTasks)

	//la primera petición como tal está hecha, ahora voy a crear el servidor
	router.Run("localhost:8080")
}
