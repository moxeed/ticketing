package app

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"strings"
	"testing"
	"time"
)

const dsn = "host=localhost user=postgres password=1qaz@WSX dbname=bartarha_ticket port=5432"

func clearTickets(db *gorm.DB) {
	value := db.Unscoped().Where("1 = 1").Delete(&Ticket{})
	data := value.Error
	fmt.Print(data)
}

func insertTicket(db *gorm.DB) {
	userId := 1
	db.Create(&Ticket{
		UserId:              &userId,
		UserName:            "sdjksdj",
		PhoneNumber:         "091092392",
		Content:             "sklskdlsdkls",
		Origin:              "sdsdsd",
		HandlerUserId:       nil,
		State:               New,
		LastStateChangeDate: time.Now(),
	})
}

func insertTicketByCreateDate(db *gorm.DB, createDateOffset int) {
	userId := 1
	db.Create(&Ticket{
		UserId:              &userId,
		UserName:            "sdjksdj",
		PhoneNumber:         "091092392",
		Content:             "sklskdlsdkls",
		Origin:              "sdsdsd",
		HandlerUserId:       nil,
		State:               New,
		LastStateChangeDate: time.Now().AddDate(0, 0, createDateOffset),
	})
}

func insertTicketByPhoneNumberAndUserName(db *gorm.DB, phoneNumber string, userName string) {
	userId := 1
	db.Create(&Ticket{
		UserId:              &userId,
		UserName:            userName,
		PhoneNumber:         phoneNumber,
		Content:             "sklskdlsdkls",
		Origin:              "sdsdsd",
		HandlerUserId:       nil,
		State:               New,
		LastStateChangeDate: time.Now(),
	})
}

func openDb() *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	clearTickets(db)

	return db
}

func TestCreateTicket(t *testing.T) {
	userId := 1
	model := TicketCreateModel{
		UserId:      &userId,
		UserName:    "tester",
		PhoneNumber: "09123456789",
		Content:     "test ticket",
		Origin:      "test",
	}

	db := openDb()
	if err := db.AutoMigrate(&Ticket{}); err != nil {
		t.Fatal(err)
	}

	result := CreateTicket(model, db)

	assert.NotEqual(t, result.Id, 0)

	assert.Equal(t, result.State, New)
	assert.Equal(t, result.Origin, model.Origin)
	assert.Equal(t, result.Content, model.Content)
	assert.Equal(t, *result.UserId, *model.UserId)
	assert.Equal(t, result.UserName, model.UserName)
	assert.Equal(t, result.PhoneNumber, model.PhoneNumber)

	dbResult := Ticket{}
	db.First(&dbResult, result.Id)

	assert.Equal(t, dbResult.State, New)
	assert.Equal(t, dbResult.Origin, model.Origin)
	assert.Equal(t, dbResult.Content, model.Content)
	assert.Equal(t, *dbResult.UserId, *model.UserId)
	assert.Equal(t, dbResult.UserName, model.UserName)
	assert.Equal(t, dbResult.PhoneNumber, model.PhoneNumber)

}

func TestCloseTicket(t *testing.T) {
	userId := 1
	ticket := NewTicket(
		&userId,
		"tester",
		"09123456789",
		"test ticket",
		"test",
	)

	db := openDb()
	if err := db.AutoMigrate(&Ticket{}); err != nil {
		t.Fatal(err)
	}

	db.Create(&ticket)
	if err := CloseTicket(ticket.ID, true, db); err != nil {
		t.Fatal(err)
	}

	successDbResult := Ticket{}
	db.First(&successDbResult, ticket.ID)

	if err := CloseTicket(ticket.ID, false, db); err != nil {
		t.Fatal(err)
	}

	failedDbResult := Ticket{}
	db.First(&failedDbResult, ticket.ID)

	assert.Equal(t, successDbResult.State, SuccessfulClosed)
	assert.Equal(t, failedDbResult.State, UnsuccessfulClosed)
}

func TestGetAllTickets(t *testing.T) {
	model := TicketFilterModel{
		OrderCol:       nil,
		OrderAscending: nil,
		Page:           nil,
		Length:         nil,
		Search:         "",
		FromDateTime:   nil,
		ToDateTime:     nil,
		State:          nil,
	}
	db := openDb()
	insertTicket(db)
	insertTicket(db)
	tickets := GetTickets(model, db)

	assert.Equal(t, len(tickets), 2)
}

func TestPageTwoIsEmptyTickets(t *testing.T) {
	page := uint(2)
	model := TicketFilterModel{
		OrderCol:       nil,
		OrderAscending: nil,
		Page:           &page,
		Length:         nil,
		Search:         "",
		FromDateTime:   nil,
		ToDateTime:     nil,
		State:          nil,
	}
	db := openDb()
	insertTicket(db)
	tickets := GetTickets(model, db)

	assert.Equal(t, len(tickets), 0)
}

func TestLengthWorksOnTickets(t *testing.T) {
	length := uint(3)
	model := TicketFilterModel{
		OrderCol:       nil,
		OrderAscending: nil,
		Page:           nil,
		Length:         &length,
		Search:         "",
		FromDateTime:   nil,
		ToDateTime:     nil,
		State:          nil,
	}

	db := openDb()

	insertTicket(db)
	insertTicket(db)
	insertTicket(db)
	insertTicket(db)

	tickets := GetTickets(model, db)

	assert.Equal(t, len(tickets), 3)
}

func orderWorksOnTickets(t *testing.T, orderAscending bool) {
	orderCol := "last_state_change_date"

	model := TicketFilterModel{
		OrderCol:       &orderCol,
		OrderAscending: &orderAscending,
		Page:           nil,
		Length:         nil,
		Search:         "",
		FromDateTime:   nil,
		ToDateTime:     nil,
		State:          nil,
	}

	db := openDb()

	insertTicketByCreateDate(db, 0)
	insertTicketByCreateDate(db, 2)
	insertTicketByCreateDate(db, -5)
	insertTicketByCreateDate(db, 5)

	tickets := GetTickets(model, db)

	assert.Equal(t, len(tickets), 4)

	prev := tickets[0]
	for index, ticket := range tickets {
		if index == 0 {
			continue
		}

		assert.Equal(t, prev.LastStateChangeDate.Before(ticket.LastStateChangeDate), orderAscending)
		prev = ticket
	}
}

func TestOrderWorksOnTickets(t *testing.T) {
	t.Run("Desc", func(t *testing.T) {
		orderWorksOnTickets(t, false)
	})
	t.Run("Asc", func(t *testing.T) {
		orderWorksOnTickets(t, false)
	})
}

func TestDateTimeFilterWorksOnTickets(t *testing.T) {
	startDate := time.Now().AddDate(0, 0, -1)
	endDate := time.Now().AddDate(0, 0, 1)
	model := TicketFilterModel{
		OrderCol:       nil,
		OrderAscending: nil,
		Page:           nil,
		Length:         nil,
		Search:         "",
		FromDateTime:   &startDate,
		ToDateTime:     &endDate,
		State:          nil,
	}

	db := openDb()

	insertTicketByCreateDate(db, -5)
	insertTicketByCreateDate(db, 0)
	insertTicketByCreateDate(db, 2)

	ticketsBetween := GetTickets(model, db)
	model.FromDateTime = nil
	ticketsBefore := GetTickets(model, db)
	model.ToDateTime = nil
	model.FromDateTime = &startDate
	ticketsAfter := GetTickets(model, db)

	assert.Equal(t, len(ticketsBetween), 1)
	for _, ticket := range ticketsBetween {
		assert.Equal(t, ticket.LastStateChangeDate.Before(endDate), true)
		assert.Equal(t, ticket.LastStateChangeDate.After(startDate), true)
	}

	assert.Equal(t, len(ticketsBefore), 2)
	for _, ticket := range ticketsBefore {
		assert.Equal(t, ticket.LastStateChangeDate.Before(endDate), true)
	}

	assert.Equal(t, len(ticketsAfter), 2)
	for _, ticket := range ticketsAfter {
		assert.Equal(t, ticket.LastStateChangeDate.After(startDate), true)
	}
}

func TestSearchFilterWorksOnTickets(t *testing.T) {
	term := "match"
	model := TicketFilterModel{
		OrderCol:       nil,
		OrderAscending: nil,
		Page:           nil,
		Length:         nil,
		Search:         "match",
		FromDateTime:   nil,
		ToDateTime:     nil,
		State:          nil,
	}

	db := openDb()

	insertTicketByPhoneNumberAndUserName(db, "09121212123", "match1")
	insertTicketByPhoneNumberAndUserName(db, "09121212124", "other")
	insertTicketByPhoneNumberAndUserName(db, "09112121215", "also_match1")

	tickets := GetTickets(model, db)

	assert.Equal(t, len(tickets), 2)
	for _, ticket := range tickets {
		assert.Equal(t, strings.Contains(ticket.UserName, term) || strings.Contains(ticket.PhoneNumber, term), true)
	}
}
