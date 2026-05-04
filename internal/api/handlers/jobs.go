package handlers

import (
	"fmt"
	"net/http"

	"github.com/Lakelimbo/kaeru/internal/jobs"
	"github.com/labstack/echo/v5"
)

type JobHandler struct {
	messaging *jobs.JobMessaging
}

func NewJobHandler(jobs *jobs.Job) *JobHandler {
	return &JobHandler{
		messaging: jobs.Messaging,
	}
}

func (h *JobHandler) Routes(e *echo.Group) {
	e.GET("/jobs/:id/stream", h.Stream)
}

// Stream is a route for server-sent-events (SSE) about a specific
// job.
//
//	@Summary		Stream
//	@Description	Realtime job status stream
//	@Tags			jobs
//	@Param			id	path	string	true	"Job ID"
//	@Produces		text/event-stream
//	@Success		200	{string}	string
//	@Failure		500	{string}	string
//	@Router			/api/v1/jobs/{id}/stream [get]
func (h *JobHandler) Stream(c *echo.Context) error {
	jobID := c.Param("id")

	res := c.Response()
	req := c.Request()

	res.Header().Set(echo.HeaderContentType, "text/event-stream")
	res.Header().Set(echo.HeaderCacheControl, "no-cache")
	res.Header().Set(echo.HeaderConnection, "keep-alive")

	rc := http.NewResponseController(res)

	ch := h.messaging.Sub(jobID)
	defer h.messaging.Unsub(jobID, ch)

	ctx := req.Context()
	for {
		select {
		case msg := <-ch:
			fmt.Fprintf(res, "data%s\n\n", msg)
			rc.Flush()

		case <-ctx.Done():
			return nil
		}
	}
}
