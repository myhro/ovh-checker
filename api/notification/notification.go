package notification

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myhro/ovh-checker/api/errors"
	"github.com/myhro/ovh-checker/storage"
	"github.com/nleof/goyesql"
)

const (
	existentNotification = "a notification for this country and server already exists"
	successfullyCreated  = "notification created successfully"
)

// Handler holds objects to be reused between requests, like a database connection
type Handler struct {
	DB      storage.DB
	Queries goyesql.Queries
}

// NewHandler creates a new Handler
func NewHandler(db storage.DB) (*Handler, error) {
	queries, err := goyesql.ParseFile("sql/notification.sql")
	if err != nil {
		return nil, err
	}

	handler := Handler{
		DB:      db,
		Queries: queries,
	}

	return &handler, nil
}

type notificationRequest struct {
	Country   string `form:"country" binding:"required"`
	Server    string `form:"server" binding:"required"`
	Recurrent bool   `form:"recurrent"`
}

// Notification creates a notification request for the current user
func (h *Handler) Notification(c *gin.Context) {
	id := c.GetInt("auth_id")

	req := notificationRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		msg := errors.ValidationMessage(err)
		errors.BadRequestWithMessage(c, msg)
		return
	}

	_, err = h.DB.Exec(h.Queries["add-notification"], id, req.Server, req.Country, req.Recurrent)
	if err != nil && storage.ErrUniqueViolation(err) {
		errors.BadRequestWithMessage(c, existentNotification)
		return
	} else if err != nil {
		log.Print(err)
		errors.InternalServerError(c)
		return
	}

	body := gin.H{
		"message": successfullyCreated,
	}

	c.JSON(http.StatusOK, body)
}
