package app

import "time"

type TicketModel struct {
	Id                  uint        `json:"id,omitempty"`
	UserId              *int        `json:"userId,omitempty"`
	ClientId            string      `json:"clientId,omitempty"`
	UserName            string      `json:"userName,omitempty"`
	PhoneNumber         string      `json:"phoneNumber,omitempty"`
	Content             string      `json:"content,omitempty"`
	Origin              string      `json:"origin,omitempty"`
	HandlerUserId       *int        `json:"handlerUserId,omitempty"`
	State               TicketState `json:"state,omitempty"`
	CreatedAt           time.Time   `json:"createdAt"`
	LastStateChangeDate time.Time   `json:"lastStateChangeDate"`
}

type TicketCreateModel struct {
	UserId      *int
	ClientId    string
	UserName    string
	PhoneNumber string
	Content     string
	Origin      string
}

type CloseTicketModel struct {
	Reason      CloseReason
	Description string
}
