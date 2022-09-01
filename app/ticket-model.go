package app

import "time"

type TicketModel struct {
	Id                  uint
	UserId              *int
	UserName            string
	PhoneNumber         string
	Content             string
	Origin              string
	HandlerUserId       *int
	State               TicketState
	CreatedAt           time.Time
	LastStateChangeDate time.Time
}

type TicketCreateModel struct {
	UserId      *int
	UserName    string
	PhoneNumber string
	Content     string
	Origin      string
}

type TicketFilterModel struct {
	OrderCol       *string
	OrderAscending *bool
	Page           *uint
	Length         *uint
	Search         string
	FromDateTime   *time.Time `time_format:"2006-01-02"`
	ToDateTime     *time.Time
	State          *TicketState
}
