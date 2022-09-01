package cmd

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"ticketing/app"
	"ticketing/common"
)

func InvalidRequest(g *gin.Context, valid interface{}) {
	g.JSON(400, valid)
}

// CreateTicket godoc
// @Tags         ticket
// @Description  Create New Ticket
// @Accept       json
// @Produce      json
// @Param        ticket body app.TicketCreateModel true "ticket data"
// @Success      200  {object}  app.TicketModel
// @Failure      400  {object}  common.Error
// @Failure      404  {object}  common.Error
// @Router       /ticket [post]
func CreateTicket(g *gin.Context) {
	model := app.TicketCreateModel{}
	if err := g.BindJSON(&model); err != nil {
		InvalidRequest(g, model)
		return
	}

	db := common.OpenDb()
	result := app.CreateTicket(model, db)

	g.JSON(200, result)
}

// CloseTicket godoc
// @Tags         ticket
// @Description  Close Ticket
// @Accept       json
// @Produce      json
// @Param        ticketId path uint true "ticket id"
// @Param        successful query bool true "is successful"
// @Success      200  {object}  app.TicketModel
// @Failure      400  {object}  common.Error
// @Failure      404  {object}  common.Error
// @Router       /admin/ticket/{ticketId} [post]
func CloseTicket(g *gin.Context) {
	var ticketId uint64
	var successful bool
	var err error

	ticketId, err = strconv.ParseUint(g.Param("ticketId"), 10, 0)
	if err != nil {
		InvalidRequest(g, "/:ticketId(unsigned number)")
		return
	}
	successful, err = strconv.ParseBool(g.Query("successful"))
	if err != nil {
		InvalidRequest(g, "?successful=true|false")
		return
	}

	db := common.OpenDb()
	result := app.CloseTicket(uint(ticketId), successful, db)

	g.JSON(200, result)
}

// GetTicket godoc
// @Tags         ticket
// @Description  Get Ticket By Id
// @Accept       json
// @Produce      json
// @Param        ticketId path uint true "ticket id"
// @Success      200  {object}  app.TicketModel
// @Failure      400  {object}  common.Error
// @Failure      404  {object}  common.Error
// @Router       /ticket/{ticketId} [get]
func GetTicket(g *gin.Context) {
	ticketId, err := strconv.ParseUint(g.Param("ticketId"), 10, 0)
	if err != nil {
		InvalidRequest(g, "/:ticketId(unsigned number)")
		return
	}

	db := common.OpenDb()
	result := app.GetTicket(uint(ticketId), db)

	g.JSON(200, result)
}

// GetTickets godoc
// @Tags         ticket
// @Description  Get Tickets By Filter Model
// @Accept       json
// @Produce      json
// @Param        filterModel body app.TicketFilterModel true "filters"
// @Success      200  {array}  app.TicketModel
// @Failure      400  {object}  common.Error
// @Router       /tickets [post]
func GetTickets(g *gin.Context) {
	model := app.TicketFilterModel{}
	if err := g.BindJSON(&model); err != nil {
		InvalidRequest(g, model)
		return
	}

	db := common.OpenDb()
	result := app.GetTickets(model, db)

	g.JSON(200, result)
}
