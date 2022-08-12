package analytics

import (
	"time"

	"github.com/hromov/jevelina/domain/leads"
)

type Filter struct {
	ResponsibleID uint64
	MinDate       time.Time
	MaxDate       time.Time
}

func (f *Filter) toLeadsFilter() leads.Filter {
	return leads.Filter{
		ResponsibleID:  f.ResponsibleID,
		MinDate:        f.MinDate,
		MaxDate:        f.MaxDate,
		ByCreationDate: true,
	}
}
