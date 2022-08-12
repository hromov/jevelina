package handlers

import (
	"net/url"
	"reflect"
	"time"

	"github.com/gorilla/schema"
	"github.com/hromov/jevelina/domain/contacts"
	"github.com/hromov/jevelina/domain/finances"
	"github.com/hromov/jevelina/domain/leads"
	"github.com/hromov/jevelina/services/events"
	"github.com/hromov/jevelina/useCases/analytics"
	"github.com/hromov/jevelina/useCases/tasks"
)

const defaultLimit = 50

type listFilter struct {
	IDs            []uint64  `schema:"ids[]"`
	Limit          int       `schema:"limit"`
	Offset         int       `schema:"offset"`
	LeadID         uint64    `schema:"lead"`
	ContactID      uint64    `schema:"contact"`
	Query          string    `schema:"query"`
	ParentID       uint64    `schema:"parent"`
	Active         bool      `schema:"active"`
	ByCreationDate bool      `schema:"by_date"`
	StepID         uint8     `schema:"step"`
	Steps          []uint8   `schema:"steps[]"`
	ResponsibleID  uint64    `schema:"responsible"`
	MinDate        time.Time `schema:"min_date"`
	MaxDate        time.Time `schema:"max_date"`
	From           uint16    `schema:"from"`
	To             uint16    `schema:"to"`
	Wallet         uint16    `schema:"wallet"`
	Completed      bool      `schema:"completed"`
}

var timeConverter = func(value string) reflect.Value {
	if v, err := time.Parse("Jan-02-2006", value); err == nil {
		return reflect.ValueOf(v)
	}
	return reflect.Value{} // this is the same as the private const invalidType
}

func parseFilter(u url.Values) (listFilter, error) {
	filter := listFilter{}
	d := schema.NewDecoder()
	d.RegisterConverter(time.Time{}, timeConverter)
	if err := d.Decode(&filter, u); err != nil {
		return listFilter{}, err
	}
	if filter.Limit == 0 {
		filter.Limit = defaultLimit
	}
	return filter, nil
}

func (f *listFilter) toFinances() finances.Filter {
	return finances.Filter{
		IDs:       f.IDs,
		Limit:     f.Limit,
		Offset:    f.Offset,
		Query:     f.Query,
		ParentID:  f.ParentID,
		MinDate:   f.MinDate,
		MaxDate:   f.MaxDate,
		From:      f.From,
		To:        f.To,
		Wallet:    f.Wallet,
		Completed: f.Completed,
	}
}

func (f *listFilter) toLeads() leads.Filter {
	return leads.Filter{
		IDs:            f.IDs,
		Limit:          f.Limit,
		Offset:         f.Offset,
		ContactID:      f.ContactID,
		Query:          f.Query,
		Active:         f.Active,
		ByCreationDate: f.ByCreationDate,
		StepID:         f.StepID,
		Steps:          f.Steps,
		ResponsibleID:  f.ResponsibleID,
		MinDate:        f.MinDate,
		MaxDate:        f.MaxDate,
		Completed:      f.Completed,
	}
}

func (f *listFilter) toAnalytics() analytics.Filter {
	return analytics.Filter{
		ResponsibleID: f.ResponsibleID,
		MinDate:       f.MinDate,
		MaxDate:       f.MaxDate,
	}
}

func (f *listFilter) toTasks() tasks.Filter {
	return tasks.Filter{
		IDs:           f.IDs,
		Limit:         f.Limit,
		Offset:        f.Offset,
		Query:         f.Query,
		ParentID:      f.ParentID,
		ResponsibleID: f.ResponsibleID,
		MinDate:       f.MinDate,
		MaxDate:       f.MaxDate,
	}
}

func (f *listFilter) toContacts() contacts.Filter {
	return contacts.Filter{
		Limit:  f.Limit,
		Offset: f.Offset,
		Query:  f.Query,
	}
}

func (f *listFilter) toEvents() events.EventFilter {
	return events.EventFilter{
		ParentID: f.ParentID,
		UserID:   f.ResponsibleID,
		Limit:    f.Limit,
		Offset:   f.Offset,
	}
}
