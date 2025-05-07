package main

import (
	"net/http"
	"strings"
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

	for _, theTask := range tasks {
		if theTask.Id == id {
			c.IndentedJSON(http.StatusOK, theTask)
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

//En esta función vamos a enfoncarnos en actualizar la tarea en funcion del id que pongamos

func putTask(c *gin.Context) {
	id := c.Param("id")

	// Estructura para recibir los nuevos datos
	var newTask Task

	// Validar si el cuerpo de la solicitud es válido
	if err := c.BindJSON(&newTask); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid"})
		return
	}

	// Buscar la tarea por ID y actualizar los campos
	for i := range tasks {
		if tasks[i].Id == id {
			// Actualizar los campos de la tarea
			tasks[i].Title = newTask.Title
			tasks[i].Description = newTask.Description
			tasks[i].Status = newTask.Status

			// Si el estado es "Completed", asignar la fecha actual
			if newTask.Status == "Completed" {
				tasks[i].CompletedAt = time.Now()
			} else {
				// Asignar un valor vacío para CompletedAt si no está completada
				tasks[i].CompletedAt = time.Time{}
			}

			// Devolver la tarea actualizada
			c.JSON(http.StatusOK, tasks[i])
			return
		}
	}

	// Si no se encuentra la tarea con ese ID
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Task not found"})
}

//función para borrar

func deleteTasks(c *gin.Context) {
	id := c.Param("id")
	for i, task := range tasks {
		if task.Id == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "task deleted"})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
}

//funcion para obtener tareas por titulo

func getTasksByTitle(c *gin.Context) {
	title := c.Param("title")
	var matchedTasks []Task

	for _, task := range tasks {
		if strings.Contains(strings.ToLower(task.Title), strings.ToLower(title)) {
			matchedTasks = append(matchedTasks, task)
		}
	}
	c.IndentedJSON(http.StatusOK, matchedTasks)
}

// función para filtrar por status
func getTasksByStatus(c *gin.Context) {
	status := c.Param("status")
	var matchedTasks []Task

	for _, task := range tasks {
		if strings.Contains(strings.ToLower(task.Status), strings.ToLower(status)) {
			matchedTasks = append(matchedTasks, task)
		}
	}
	c.IndentedJSON(http.StatusOK, matchedTasks)
}

func main() {
	//vamos a ir creando las rutas

	//para inicializar las rutas crearé una variable llamada router
	router := gin.Default()

	//esto nos devuelve toda la data
	router.GET("/tasks", getTasks)

	//esto nos filtra por id
	router.GET("/tasks/:id", getTasksById)

	//esto nos permite añadir tareas
	router.POST("/tasks", postTasks)

	//esto nos permite actualizar tareas por ID
	router.PUT("/tasks/:id", putTask)

	//esto borra tareas por id
	router.DELETE("/tasks/:id", deleteTasks)

	//esto sirve para obtener las tasks por título
	router.GET("/tasks/title/:title", getTasksByTitle)

	//sirve para filtrar por status
	router.GET("/tasks/status/:status", getTasksByStatus)

	//la primera petición como tal está hecha, ahora voy a crear el servidor
	router.Run("localhost:8080")

}
