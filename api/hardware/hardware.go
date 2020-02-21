package hardware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myhro/ovh-checker/api/errors"
	"github.com/myhro/ovh-checker/database"
	"github.com/myhro/ovh-checker/models/hardware"
	"github.com/nleof/goyesql"
)

// Handler holds objects to be reused between requests, like a database connection
type Handler struct {
	DB      database.DB
	Queries goyesql.Queries
}

type offersRequest struct {
	Country string `form:"country" binding:"required"`
	First   int    `form:"first" binding:"required"`
	Last    int    `form:"last" binding:"required"`
}

// NewHandler creates a new Handler
func NewHandler() (*Handler, error) {
	db, err := database.New()
	if err != nil {
		return nil, err
	}

	queries, err := goyesql.ParseFile("sql/hardware.sql")
	if err != nil {
		return nil, err
	}

	handler := Handler{
		DB:      db,
		Queries: queries,
	}

	return &handler, nil
}

// Offers returns a list of latest offers for the requested hardware
func (h *Handler) Offers(c *gin.Context) {
	offersReq := offersRequest{}
	err := c.ShouldBind(&offersReq)
	if err != nil {
		msg := errors.ValidationMessage(err)
		errors.BadRequestWithMessage(c, msg)
		return
	}

	list := []hardware.LatestOffers{}
	err = h.DB.Select(&list, h.Queries["latest-offers"], offersReq.Country, offersReq.First, offersReq.Last)
	if err != nil {
		log.Print(err)
		errors.InternalServerError(c)
		return
	}

	c.JSON(http.StatusOK, list)
}
