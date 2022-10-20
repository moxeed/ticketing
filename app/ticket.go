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

const (
	NoResponseFromStudent CloseReason = iota
	InvalidTicket
	Expired
	InvalidContact
	Other
)

type TicketState int
type CloseReason int

type Ticket struct {
	gorm.Model
	UserId              *int
	ClientId            string
	UserName            string
	PhoneNumber         string
	Content             string
	Origin              string
	State               TicketState
	CloseReason         *CloseReason
	CloseDescription    *string
	HandlerUserId       *int
	LastStateChangeDate time.Time
}

func NewTicket(userId *int, clientId string, userName string, phoneNumber string, content string, origin string) Ticket {
	return Ticket{
		UserId:              userId,
		ClientId:            clientId,
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

func (ticket *Ticket) Close(isSuccessful bool, reason CloseReason, description string) {
	ticket.CloseReason = &reason
	ticket.CloseDescription = &description
	if isSuccessful {
		ticket.setState(SuccessfulClosed)
	} else {
		ticket.setState(UnsuccessfulClosed)
	}
}

func CreateTicket(model TicketCreateModel, db *gorm.DB) TicketModel {
	ticket := NewTicket(model.UserId,
		model.ClientId,
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

func CloseTicket(ticketId uint, isSuccessful bool, closeReason CloseReason, closeDescription string, db *gorm.DB) error {
	ticket := Ticket{}

	db.First(&ticket, ticketId)
	if ticket.ID == 0 {
		return fmt.Errorf("ticket Not Found")
	}

	ticket.Close(isSuccessful, closeReason, closeDescription)
	db.Save(&ticket)

	return nil
}

func GetTicket(ticketId uint, db *gorm.DB) TicketModel {
	ticket := Ticket{}
	db.First(&ticket, ticketId)
	return convert(ticket)
}

func convert(ticket Ticket) TicketModel {
	return TicketModel{
		Id:                  ticket.ID,
		UserId:              ticket.UserId,
		ClientId:            ticket.ClientId,
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
