package hardware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myhro/ovh-checker/api/errors"
	"github.com/myhro/ovh-checker/models/hardware"
)

type offersRequest struct {
	Country string `form:"country" binding:"required"`
	First   int    `form:"first" binding:"required"`
	Last    int    `form:"last" binding:"required"`
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
