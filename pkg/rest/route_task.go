package rest

import (
	"bytes"
	"encoding/csv"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leowilbur/task_manager/pkg/models"
	"github.com/robfig/cron/v3"
)

// GetTaskList lists all tasks
func (a *API) GetTaskList(r *gin.Context) {
	tasks := []models.Task{}
	for _, task := range a.Context.JobInfo {
		tasks = append(tasks, task)
	}

	r.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    tasks,
	})
}

// CreateTask inserts a new task
func (a *API) CreateTask(r *gin.Context) {
	var input models.Task
	if err := r.BindJSON(&input); err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON input: " + err.Error(),
		})
		return
	}

	entryID, err := a.Context.Cron.AddFunc(input.Cron, func() {
		log.Printf("Execute cron job:%s at %v \n", input.Name, time.Now())
	})

	if err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Fail to add new cron task: " + err.Error(),
		})
		return
	}

	task := models.Task{ID: int(entryID), Name: input.Name, Cron: input.Cron}
	a.Context.JobInfo[entryID] = task

	r.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    task,
	})
}

// DeleteTask remove the exist task
func (a *API) DeleteTask(r *gin.Context) {
	taskID, err := strconv.ParseInt(r.Param("task_id"), 10, 64)
	if err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input: " + err.Error(),
		})
		return
	}

	a.Context.Cron.Remove(cron.EntryID(taskID))
	delete(a.Context.JobInfo, cron.EntryID(taskID))

	r.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

// TaskStart start tasks
func (a *API) TaskStart(r *gin.Context) {
	a.Context.Cron.Start()

	r.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

// TaskStop stop tasks
func (a *API) TaskStop(r *gin.Context) {
	a.Context.Cron.Stop()

	r.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

// TaskImport import tasks from csv file
func (a *API) TaskImport(r *gin.Context) {
	rFile, err := r.FormFile("file")
	if err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid file upload, Error: " + err.Error(),
		})
		return
	}

	file, err := rFile.Open()
	if err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid file format, Error: " + err.Error(),
		})
		return
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 2

	csvData, err := reader.ReadAll()
	if err != nil {
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid csv file format, Error: " + err.Error(),
		})
		return
	}

	tasks := []models.Task{}
	for i := 1; i < len(csvData); i++ { // skip header
		task := models.Task{Name: csvData[i][0], Cron: csvData[i][1]}
		entryID, err := a.Context.Cron.AddFunc(task.Cron, func() {
			log.Printf("Execute cron job:%s at %v \n", task.Name, time.Now())
		})

		if err != nil {
			log.Printf("Fail to add new cron task: " + err.Error())
		}
		task.ID = int(entryID)
		a.Context.JobInfo[entryID] = task
		tasks = append(tasks, task)
	}

	r.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    tasks,
	})
}

// TaskExport export tasks to csv file
func (a *API) TaskExport(r *gin.Context) {
	b := &bytes.Buffer{}
	w := csv.NewWriter(b)

	if err := w.Write([]string{"Name", "Cron"}); err != nil {
		log.Fatalln("error writing record to csv:", err)
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "error writing record to csv:" + err.Error(),
		})
		return
	}

	for _, task := range a.Context.JobInfo {
		var record []string
		record = append(record, task.Name)
		record = append(record, task.Cron)
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
			r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "error writing record to csv:" + err.Error(),
			})
			return
		}
	}
	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
		r.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "error writing record to csv:" + err.Error(),
		})
	}

	r.Data(http.StatusOK, "text/csv", b.Bytes())
}
