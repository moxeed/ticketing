package app

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

const (
	New TicketState = iota
	Open
	Resolved
	WaitingForInput
	SuccessfulClosed
	UnsuccessfulClosed
)

type TicketState int

type Ticket struct {
	gorm.Model
	UserId              *int
	UserName            string
	PhoneNumber         string
	Content             string
	Origin              string
	State               TicketState
	HandlerUserId       *int
	LastStateChangeDate time.Time
}

func NewTicket(userId *int, userName string, phoneNumber string, content string, origin string) Ticket {
	return Ticket{
		UserId:              userId,
		UserName:            userName,
		PhoneNumber:         phoneNumber,
		Content:             content,
		Origin:              origin,
		State:               New,
		LastStateChangeDate: time.Now(),
	}
}

func (ticket *Ticket) setState(state TicketState) {
	ticket.State = state
	ticket.LastStateChangeDate = time.Now()
}

func (ticket *Ticket) Resolve(handlerId int) {
	ticket.HandlerUserId = &handlerId
	ticket.setState(Resolved)
}

func (ticket *Ticket) Close(isSuccessful bool) {
	if isSuccessful {
		ticket.setState(SuccessfulClosed)
	} else {
		ticket.setState(UnsuccessfulClosed)
	}
}

func CreateTicket(model TicketCreateModel, db *gorm.DB) TicketModel {
	ticket := NewTicket(model.UserId,
		model.UserName,
		model.PhoneNumber,
		model.Content,
		model.Origin)

	db.Save(&ticket)

	return TicketModel{
		Id:                  ticket.ID,
		UserId:              ticket.UserId,
		UserName:            ticket.UserName,
		PhoneNumber:         ticket.PhoneNumber,
		Content:             ticket.Content,
		Origin:              ticket.Origin,
		HandlerUserId:       ticket.HandlerUserId,
		State:               ticket.State,
		CreatedAt:           ticket.CreatedAt,
		LastStateChangeDate: ticket.LastStateChangeDate,
	}
}

func CloseTicket(ticketId uint, isSuccessful bool, db *gorm.DB) error {
	ticket := Ticket{}

	db.First(&ticket, ticketId)
	if ticket.ID == 0 {
		return fmt.Errorf("ticket Not Found")
	}

	ticket.Close(isSuccessful)
	db.Save(&ticket)

	return nil
}

func parseFilter(filterModel TicketFilterModel) (int, int, string, string) {
	page := uint(1)
	if filterModel.Page != nil {
		page = *filterModel.Page
	}

	length := uint(10)
	if filterModel.Length != nil {
		length = *filterModel.Length
	}

	orderCol := "id"
	if filterModel.OrderCol != nil {
		orderCol = *filterModel.OrderCol
	}

	orderDirection := "desc"
	if filterModel.OrderAscending != nil && *filterModel.OrderAscending {
		orderDirection = "asc"
	}

	return int(page), int(length), orderCol, orderDirection
}

func GetTicket(ticketId uint, db *gorm.DB) TicketModel {
	ticket := Ticket{}
	db.First(&ticket, ticketId)
	return convertSingle(ticket)
}

func GetTickets(filterModel TicketFilterModel, db *gorm.DB) []TicketModel {
	page, length, orderCol, orderDirection := parseFilter(filterModel)

	tickets := make([]Ticket, 0)

	query := db

	if filterModel.FromDateTime != nil {
		query = query.Where("last_state_change_date > ?", *filterModel.FromDateTime)
	}
	if filterModel.ToDateTime != nil {
		query = query.Where("last_state_change_date < ?", *filterModel.ToDateTime)
	}
	if filterModel.State != nil {
		query = query.Where("state = ?", *filterModel.State)
	}
	if filterModel.Search != "" {
		likeTerm := fmt.Sprintf("%%%s%%", filterModel.Search)
		query = query.Where("phone_number Like ? OR user_name Like ?", likeTerm, likeTerm)
	}

	query.Order(fmt.Sprintf("%s %s", orderCol, orderDirection)).
		Offset((page - 1) * length).
		Limit(length).
		Find(&tickets)

	return convert(tickets)
}

func convert(tickets []Ticket) []TicketModel {
	models := make([]TicketModel, 0)

	for _, ticket := range tickets {
		models = append(models, TicketModel{
			Id:                  ticket.ID,
			UserId:              ticket.UserId,
			UserName:            ticket.UserName,
			PhoneNumber:         ticket.PhoneNumber,
			Content:             ticket.Content,
			Origin:              ticket.Origin,
			HandlerUserId:       ticket.HandlerUserId,
			State:               ticket.State,
			CreatedAt:           ticket.CreatedAt,
			LastStateChangeDate: ticket.LastStateChangeDate,
		})
	}

	return models
}

func convertSingle(ticket Ticket) TicketModel {
	return TicketModel{
		Id:                  ticket.ID,
		UserId:              ticket.UserId,
		UserName:            ticket.UserName,
		PhoneNumber:         ticket.PhoneNumber,
		Content:             ticket.Content,
		Origin:              ticket.Origin,
		HandlerUserId:       ticket.HandlerUserId,
		State:               ticket.State,
		CreatedAt:           ticket.CreatedAt,
		LastStateChangeDate: ticket.LastStateChangeDate,
	}
}
