package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gocraft/work"
	"github.com/rs/cors"

	"github.com/leowilbur/task_manager/docs/api"
	"github.com/leowilbur/task_manager/pkg/models"
)

// API is the REST API
type API struct {
	*gin.Engine
	Enqueuer *work.Enqueuer
	Context  *models.Context
}

// New creates a new API using the given dependencies
func New(context *models.Context) (*API, error) {
	gin.SetMode(gin.ReleaseMode)

	r := &API{
		Engine:  gin.New(),
		Context: context,
	}

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	corsMiddleware := cors.AllowAll()
	r.Use(func(c *gin.Context) {
		corsMiddleware.HandlerFunc(c.Writer, c.Request)
	})

	r.GET("/swagger.json", func(r *gin.Context) {
		r.Header("Content-Type", "application/json")
		r.String(http.StatusOK, api.GetSwaggerJSON())
	})

	r.GET("/tasks", r.GetTaskList)
	r.POST("/tasks", r.CreateTask)
	r.DELETE("/tasks/:task_id", r.DeleteTask)
	r.POST("/task/start", r.TaskStart)
	r.POST("/task/stop", r.TaskStop)
	r.POST("/task/import", r.TaskImport)
	r.POST("/task/export", r.TaskExport)

	return r, nil
}
