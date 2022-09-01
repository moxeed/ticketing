package app

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"log"
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
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})

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
