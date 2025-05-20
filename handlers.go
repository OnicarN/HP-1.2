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

func getTasks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, tasks)
}

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

func postTasks(c *gin.Context) {
	newTaskId := uuid.New().String()
	var newTask Task

	if err := c.BindJSON(&newTask); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "you can not request"})
		return
	}

	switch {
	case strings.TrimSpace(newTask.Title) == "":
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "you have to write the title!",
		})
		return

	case strings.TrimSpace(newTask.Description) == "":
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "you have to write a description!",
		})
		return

	case strings.TrimSpace(newTask.Status) == "":
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "you have to choose a not null status!",
		})
		return
	}

	statusLowerCase := strings.ToLower(newTask.Status)
	valid := false

	for _, status := range statusOptions {
		if statusLowerCase == status {
			valid = true
			break
		}
	}

	if !valid {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "This status is not valid",
		})
	}

	newTask.Id = newTaskId
	newTask.CreatedAt = time.Now()

	if strings.ToLower(newTask.Status) == "completed" {
		newTask.CompletedAt = time.Now()
	} else {
		newTask.CompletedAt = time.Time{}
	}

	tasks = append(tasks, newTask)

	c.IndentedJSON(http.StatusCreated, tasks)
}

func putTask(c *gin.Context) {
	id := c.Param("id")

	var newTask Task

	if err := c.BindJSON(&newTask); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid"})
		return
	}

	for i := range tasks {
		if tasks[i].Id == id {

			tasks[i].Title = newTask.Title
			tasks[i].Description = newTask.Description
			tasks[i].Status = newTask.Status

			if strings.ToLower(newTask.Status) == "completed" {
				tasks[i].CompletedAt = time.Now()
			} else {

				tasks[i].CompletedAt = time.Time{}
			}

			c.JSON(http.StatusOK, tasks[i])
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Task not found"})
}

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

func getTasksByTitle(c *gin.Context) {
	title := c.Query("title")
	var matchedTasks []Task

	for _, task := range tasks {
		if strings.Contains(strings.ToLower(task.Title), strings.ToLower(title)) {
			matchedTasks = append(matchedTasks, task)
		}
	}
	c.IndentedJSON(http.StatusOK, matchedTasks)
}

func getTasksByStatus(c *gin.Context) {
	status := c.Query("status")
	var matchedTasks []Task

	for _, task := range tasks {
		if strings.Contains(strings.ToLower(task.Status), strings.ToLower(status)) {
			matchedTasks = append(matchedTasks, task)
		}
	}
	c.IndentedJSON(http.StatusOK, matchedTasks)
}
