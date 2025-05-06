package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Task struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"CreatedAt"`
	CompletedAt time.Time `json:"CompletedAt"`
}

var tasks = []Task{}

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
	var newTask Task

	//aquí usamos la referencia de la variable que hemos creado
	c.BindJSON(&newTask)

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
