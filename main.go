package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type Task struct {
	Id            string     `json:"id"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"CreatedAt"`
	CompletededAt *time.Time `json:"CompletedAt"` 
}

var tasks = []Task{}
var statusOptions = []string{"new", "ongoing", "completed"}

var db *sql.DB

func getTasks(c *gin.Context) {
	tasks, err := getTheTasks()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving tasks"})
		return
	}

	c.IndentedJSON(http.StatusOK, tasks)
}

func getTheTasks() ([]Task, error) {

	rows, err := db.Query("SELECT * FROM task")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var tsk Task
		err := rows.Scan(&tsk.Id, &tsk.Title, &tsk.Description, &tsk.Status, &tsk.CreatedAt, &tsk.CompletededAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, tsk)
	}
	return tasks, nil

}

func getTaskById(id string) (Task, error) {
	var tsk Task

	row := db.QueryRow("SELECT id, title, description, status, createdAt, completededAt FROM task WHERE id = ?", id)
	err := row.Scan(&tsk.Id, &tsk.Title, &tsk.Description, &tsk.Status, &tsk.CreatedAt, &tsk.CompletededAt) // ← Corregido
	if err != nil {
		if err == sql.ErrNoRows {
			return tsk, fmt.Errorf("getTaskById %q: no such task", id)
		}
		return tsk, fmt.Errorf("getTaskById %q: %v", id, err)
	}

	return tsk, nil
}

func getTaskByIdHandler(c *gin.Context) {
	id := c.Param("id")
	task, err := getTaskById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, task)
}

//Hasta ahí hemos visto el tema de los getters

// ahora vamos a darle al tema del post

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
		return
	}

	valid := false
	for _, status := range statusOptions {
		if newTask.Status == status {
			valid = true
			break
		}
	}

	if !valid {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "This status is not valid",
		})
		return
	}

	newTask.Id = newTaskId
	newTask.CreatedAt = time.Now()

	if newTask.Status == "completed" {
		now := time.Now()
		newTask.CompletededAt = &now
	} else {
		newTask.CompletededAt = nil
	}

	if success := insertTaskInDb(newTask); success {
		c.JSON(http.StatusCreated, newTask)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error inserting task"})
	}
}

func insertTaskInDb(task Task) bool {
	stmt, err := db.Prepare("INSERT INTO task (id, title, description, status, createdAt, completededAt) VALUES (?, ?, ?, ?, ?, ?)")

	if err != nil {
		fmt.Println("Error preparing statement: ", err)
		return false
	}

	_, execErr := stmt.Exec(task.Id, task.Title, task.Description, task.Status, task.CreatedAt, task.CompletededAt)
	if execErr != nil {
		fmt.Println("Error executing statement:", execErr)
		return false
	}
	return true

}

// ahora vamos a escribir el código para actaulizar las tareas

func updateTaskInDb(id string, task Task) error {
	stmt, err := db.Prepare("UPDATE task SET title = ?, description = ?, status = ?, completededAt = ? WHERE id = ?")

	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(task.Title, task.Description, task.Status, task.CompletededAt, id)
	return err
}

func putTask(c *gin.Context) {
	id := c.Param("id")
	var newTask Task

	if err := c.BindJSON(&newTask); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid"})
		return
	}

	if strings.ToLower(newTask.Status) == "completed" {
		now := time.Now()
		newTask.CompletededAt = &now
	} else {
		newTask.CompletededAt = nil
	}

	err := updateTaskInDb(id, newTask)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to update task"})
		return
	}

	c.IndentedJSON(http.StatusOK, newTask)
}

//vale ahora vamos a centrarnos en borrar lo que viene siendo las tasks de la base de datos,

func deleteAllTasksById(taskId string) (int64, error) {
	stmt, err := db.Prepare("DELETE FROM task WHERE Id = ?")

	if err != nil {
		log.Print("An error occured during the delete tasks")
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(taskId)

	if err != nil {
		log.Print("Exec error")
		return 0, err
	}

	return res.RowsAffected()
}

func deleteTasks(c *gin.Context) {
	id := c.Param("id")

	rowsAffected, err := deleteAllTasksById(id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error deleting task"})
		return
	}

	if rowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Task not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Deleted %d task(s)", rowsAffected)})
}

func getTasksByTitle(title string) ([]Task, error) {
	query := "SELECT id, title, description, status, createdAt, completededAt FROM task WHERE title LIKE ?"
	rows, err := db.Query(query, "%"+title+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matchedTasks []Task
	for rows.Next() {
		var tsk Task
		err := rows.Scan(&tsk.Id, &tsk.Title, &tsk.Description, &tsk.Status, &tsk.CreatedAt, &tsk.CompletededAt)
		if err != nil {
			return nil, err
		}
		matchedTasks = append(matchedTasks, tsk)
	}
	return matchedTasks, nil
}

func getTaskByTitleHandler(c *gin.Context) {
	title := c.Param("title")
	task, err := getTasksByTitle(title)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, task)
}

/**
func getTasksByTitle(c *gin.Context) {
	title := c.Param("title")
	var matchedTasks []Task

	for _, task := range tasks {
		if strings.Contains(strings.ToLower(task.Title), strings.ToLower(title)) {
			matchedTasks = append(matchedTasks, task)
		}
	}
	c.IndentedJSON(http.StatusOK, matchedTasks)
}*/

//En esta última función vamos a buscar por status

func getTasksByStatusFromDb(status string) ([]Task, error) {
	query := "SELECT id, title, description, status, createdAt, completededAt FROM task WHERE status = ?"
	rows, err := db.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matchedTasks []Task
	for rows.Next() {
		var tsk Task
		err := rows.Scan(&tsk.Id, &tsk.Title, &tsk.Description, &tsk.Status, &tsk.CreatedAt, &tsk.CompletededAt)
		if err != nil {
			return nil, err
		}
		matchedTasks = append(matchedTasks, tsk)
	}
	return matchedTasks, nil
}

func getTasksByStatus(c *gin.Context) {
	status := c.Param("status")

	tasks, err := getTasksByStatusFromDb(status)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving tasks by status"})
		return
	}

	c.IndentedJSON(http.StatusOK, tasks)
}

func main() {
	cfg := mysql.NewConfig()
	cfg.User = "root"
	cfg.Passwd = "DaniPizzaeloy1"
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = "tasks1.2"
	cfg.ParseTime = true

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	router := gin.Default()
	router.GET("/tasks", getTasks)
	router.GET("/tasks/:id", getTaskByIdHandler)
	router.POST("/tasks", postTasks)
	router.PUT("/tasks/:id", putTask)
	router.DELETE("/tasks/:id", deleteTasks)
	router.GET("/tasks/title/:title", getTaskByTitleHandler)
	router.GET("/tasks/status/:status", getTasksByStatus)

	router.Run("localhost:8080")
}
