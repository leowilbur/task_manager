package rest_test

import (
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/robfig/cron/v3"
	resty "gopkg.in/resty.v1"

	"github.com/leowilbur/task_manager/pkg/models"
	"github.com/leowilbur/task_manager/pkg/rest"
)

func TestRest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "./")
}

var _ = Describe("Tasks API Description", func() {
	Context("Task manager api should be work fine", func() {
		var (
			api       *rest.API
			apiServer *httptest.Server
			client    *resty.Client
			err       error
		)
		BeforeEach(func() {
			workerContext := &models.Context{
				Cron:    cron.New(),
				JobInfo: make(map[cron.EntryID]models.Task),
			}

			api, err = rest.New(workerContext)
			Expect(err).To(BeNil())

			apiServer = httptest.NewServer(api)

			client = resty.New().SetHostURL(apiServer.URL)

			Expect(client).NotTo(BeNil())
		})

		AfterEach(func() {
			apiServer.Close()
		})

		It("should get task list success", func() {
			resp, err := client.NewRequest().
				Get("/tasks")
			Expect(err).To(BeNil())
			Expect(resp.IsError()).To(BeFalse())
		})

		It("should create task success", func() {
			resp, err := client.NewRequest().
				SetBody(map[string]interface{}{
					"name": "Run every minutes",
					"cron": "* * * * *",
				}).
				Post("/tasks")
			Expect(err).To(BeNil())
			Expect(resp.IsError()).To(BeFalse())
		})

		It("should start task success", func() {
			resp, err := client.NewRequest().
				SetBody(nil).
				Post("/task/start")
			Expect(err).To(BeNil())
			Expect(resp.IsError()).To(BeFalse())
		})

		It("should stop task success", func() {
			resp, err := client.NewRequest().
				SetBody(nil).
				Post("/task/stop")
			Expect(err).To(BeNil())
			Expect(resp.IsError()).To(BeFalse())
		})

		It("should delete task success", func() {
			resp, err := client.NewRequest().
				Delete("/tasks/1")
			Expect(err).To(BeNil())
			Expect(resp.IsError()).To(BeFalse())
		})
	})
})
