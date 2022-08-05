package contacts

import (
	"time"

	"github.com/hromov/jevelina/domain/misc"
	"github.com/hromov/jevelina/domain/users"
)

type Contact struct {
	ID        uint64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time

	Name        string
	SecondName  string
	Responsible users.User
	Created     users.User
	Phone       string
	SecondPhone string
	Email       string
	SecondEmail string
	URL         string

	City    string
	Address string

	Source   misc.Source
	Position string

	Analytics misc.Analytics
}

type ContactRequest struct {
	ID            uint64
	Name          string
	SecondName    string
	ResponsibleID uint64
	CreatedID     uint64
	Phone         string
	SecondPhone   string
	Email         string
	SecondEmail   string
	URL           string

	City    string
	Address string

	SourceID uint8
	Position string

	Analytics misc.Analytics
}

type Filter struct {
	Limit  int
	Offset int
	TagID  uint8
	Query  string
}

type ContactsResponse struct {
	Contacts []Contact
	Total    int64
}
