/*
	Este proyecto fue realizado por Hécor Daniel Polanco Menaç
	Solo he realizado la parte obligatoria
	Para buscar en las rutas como bien se sabe hay que poner esto
	si es para el "get tasks" típico por ejemplo ponemos : http://localhost:8080/api/tasks
	Eso es todo, espero que el proyecto sea de su agrado
	Un saludo.

*/

package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	api := router.Group("/api")

	api.GET("/tasks", getTasks)

	api.GET("/tasks/:id", getTasksById)

	api.POST("/tasks", postTasks)

	api.PUT("/tasks/:id", putTask)

	api.DELETE("/tasks/:id", deleteTasks)

	api.GET("/tasks/title", getTasksByTitle)

	api.GET("/tasks/status", getTasksByStatus)

	router.Run("localhost:8080")

}
