package handlers

import (
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hromov/jevelina/domain/contacts"
	"github.com/hromov/jevelina/domain/finances"
	"github.com/hromov/jevelina/domain/leads"
	"github.com/hromov/jevelina/storage/mysql/dao/models"
	"github.com/hromov/jevelina/useCases/tasks"
)

const defaultLimit = 50

func FilterFromQuery(u url.Values) models.ListFilter {
	filter := models.ListFilter{}
	if IDs := u.Get("ids"); IDs != "" {
		filter.IDs = make([]uint64, 0)
		slice := strings.Split(IDs, ",")
		for _, IDstring := range slice {
			if IDnumber, err := strconv.ParseUint(IDstring, 10, 64); err == nil {
				filter.IDs = append(filter.IDs, IDnumber)
			}

		}
	}
	if steps := u.Get("steps"); steps != "" {
		filter.Steps = make([]uint8, 0)
		slice := strings.Split(steps, ",")
		for _, stepString := range slice {
			if stepID, err := strconv.ParseUint(stepString, 10, 8); err == nil {
				filter.Steps = append(filter.Steps, uint8(stepID))
			}
		}
	}
	if query := u.Get("query"); query != "" {
		filter.Query = query
	}
	if limit := u.Get("limit"); limit != "" {
		filter.Limit, _ = strconv.Atoi(limit)
	} else {
		filter.Limit = defaultLimit
	}
	if offset := u.Get("offset"); offset != "" {
		filter.Offset, _ = strconv.Atoi(offset)
	}
	if leadID := u.Get("lead"); leadID != "" {
		filter.LeadID, _ = strconv.ParseUint(leadID, 10, 64)
	}
	if contactID := u.Get("contact"); contactID != "" {
		filter.ContactID, _ = strconv.ParseUint(contactID, 10, 64)
	}
	if parentID := u.Get("parent"); parentID != "" {
		filter.ParentID, _ = strconv.ParseUint(parentID, 10, 64)
	}
	if from := u.Get("from"); from != "" {
		from64, _ := strconv.ParseUint(from, 10, 64)
		filter.From = uint16(from64)
	}
	if completed := u.Get("completed"); completed != "" {
		filter.Completed, _ = strconv.ParseBool(completed)
	}
	if wallet := u.Get("wallet"); wallet != "" {
		wallet64, _ := strconv.ParseUint(wallet, 10, 64)
		filter.Wallet = uint16(wallet64)
	}
	if to := u.Get("to"); to != "" {
		to64, _ := strconv.ParseUint(to, 10, 64)
		filter.To = uint16(to64)
	}
	if respID := u.Get("responsible"); respID != "" {
		filter.ResponsibleID, _ = strconv.ParseUint(respID, 10, 64)
	}
	if active := u.Get("active"); active != "" {
		filter.Active = true
	}
	if tagID := u.Get("tag"); tagID != "" {
		tag64, _ := strconv.ParseUint(tagID, 10, 64)
		filter.TagID = uint8(tag64)
	}
	const timeForm = "Jan-02-2006"
	if minDate := u.Get("min_date"); minDate != "" {
		filter.MinDate, _ = time.Parse(timeForm, minDate)
	}
	if maxDate := u.Get("max_date"); maxDate != "" {
		filter.MaxDate, _ = time.Parse(timeForm, maxDate)
	}

	if stepID := u.Get("step"); stepID != "" {
		step64, _ := strconv.ParseUint(stepID, 10, 64)
		filter.StepID = uint8(step64)
	}
	return filter
}

func TasksFilter(u url.Values) tasks.Filter {
	filter := tasks.Filter{}
	if IDs := u.Get("ids"); IDs != "" {
		filter.IDs = make([]uint64, 0)
		slice := strings.Split(IDs, ",")
		for _, IDstring := range slice {
			if IDnumber, err := strconv.ParseUint(IDstring, 10, 64); err == nil {
				filter.IDs = append(filter.IDs, IDnumber)
			}

		}
	}
	if query := u.Get("query"); query != "" {
		filter.Query = query
	}
	if limit := u.Get("limit"); limit != "" {
		filter.Limit, _ = strconv.Atoi(limit)
	} else {
		filter.Limit = defaultLimit
	}
	if offset := u.Get("offset"); offset != "" {
		filter.Offset, _ = strconv.Atoi(offset)
	}
	if parentID := u.Get("parent"); parentID != "" {
		filter.ParentID, _ = strconv.ParseUint(parentID, 10, 64)
	}
	if respID := u.Get("responsible"); respID != "" {
		filter.ResponsibleID, _ = strconv.ParseUint(respID, 10, 64)
	}
	const timeForm = "Jan-02-2006"
	if minDate := u.Get("min_date"); minDate != "" {
		filter.MinDate, _ = time.Parse(timeForm, minDate)
	}
	if maxDate := u.Get("max_date"); maxDate != "" {
		filter.MaxDate, _ = time.Parse(timeForm, maxDate)
	}
	return filter
}

func ContactsFilter(u url.Values) contacts.Filter {
	filter := contacts.Filter{}

	if query := u.Get("query"); query != "" {
		filter.Query = query
	}
	if limit := u.Get("limit"); limit != "" {
		filter.Limit, _ = strconv.Atoi(limit)
	} else {
		filter.Limit = defaultLimit
	}
	if offset := u.Get("offset"); offset != "" {
		filter.Offset, _ = strconv.Atoi(offset)
	}

	if tagID := u.Get("tag"); tagID != "" {
		tag64, _ := strconv.ParseUint(tagID, 10, 64)
		filter.TagID = uint8(tag64)
	}
	return filter
}

func LeadsFilter(u url.Values) leads.Filter {
	filter := leads.Filter{}
	if IDs := u.Get("ids"); IDs != "" {
		filter.IDs = make([]uint64, 0)
		slice := strings.Split(IDs, ",")
		for _, IDstring := range slice {
			if IDnumber, err := strconv.ParseUint(IDstring, 10, 64); err == nil {
				filter.IDs = append(filter.IDs, IDnumber)
			}

		}
	}
	if steps := u.Get("steps"); steps != "" {
		filter.Steps = make([]uint8, 0)
		slice := strings.Split(steps, ",")
		for _, stepString := range slice {
			if stepID, err := strconv.ParseUint(stepString, 10, 8); err == nil {
				filter.Steps = append(filter.Steps, uint8(stepID))
			}
		}
	}
	if query := u.Get("query"); query != "" {
		filter.Query = query
	}
	if limit := u.Get("limit"); limit != "" {
		filter.Limit, _ = strconv.Atoi(limit)
	} else {
		filter.Limit = defaultLimit
	}
	if offset := u.Get("offset"); offset != "" {
		filter.Offset, _ = strconv.Atoi(offset)
	}
	if contactID := u.Get("contact"); contactID != "" {
		filter.ContactID, _ = strconv.ParseUint(contactID, 10, 64)
	}
	if completed := u.Get("completed"); completed != "" {
		filter.Completed, _ = strconv.ParseBool(completed)
	}
	log.Println(u.Get("by_date"))
	if byDate := u.Get("by_date"); byDate != "" {
		filter.ByCreationDate, _ = strconv.ParseBool(byDate)
	}
	log.Println(filter.ByCreationDate)
	if respID := u.Get("responsible"); respID != "" {
		filter.ResponsibleID, _ = strconv.ParseUint(respID, 10, 64)
	}
	if active := u.Get("active"); active != "" {
		filter.Active = true
	}
	const timeForm = "Jan-02-2006"
	if minDate := u.Get("min_date"); minDate != "" {
		filter.MinDate, _ = time.Parse(timeForm, minDate)
	}
	if maxDate := u.Get("max_date"); maxDate != "" {
		filter.MaxDate, _ = time.Parse(timeForm, maxDate)
	}

	if stepID := u.Get("step"); stepID != "" {
		step64, _ := strconv.ParseUint(stepID, 10, 64)
		filter.StepID = uint8(step64)
	}
	return filter
}

func FinFilter(u url.Values) finances.Filter {
	filter := finances.Filter{}
	if IDs := u.Get("ids"); IDs != "" {
		filter.IDs = make([]uint64, 0)
		slice := strings.Split(IDs, ",")
		for _, IDstring := range slice {
			if IDnumber, err := strconv.ParseUint(IDstring, 10, 64); err == nil {
				filter.IDs = append(filter.IDs, IDnumber)
			}

		}
	}

	if query := u.Get("query"); query != "" {
		filter.Query = query
	}
	if limit := u.Get("limit"); limit != "" {
		filter.Limit, _ = strconv.Atoi(limit)
	} else {
		filter.Limit = defaultLimit
	}
	if offset := u.Get("offset"); offset != "" {
		filter.Offset, _ = strconv.Atoi(offset)
	}
	if parentID := u.Get("parent"); parentID != "" {
		filter.ParentID, _ = strconv.ParseUint(parentID, 10, 64)
	}
	if from := u.Get("from"); from != "" {
		from64, _ := strconv.ParseUint(from, 10, 64)
		filter.From = uint16(from64)
	}
	if completed := u.Get("completed"); completed != "" {
		filter.Completed, _ = strconv.ParseBool(completed)
	}
	if wallet := u.Get("wallet"); wallet != "" {
		wallet64, _ := strconv.ParseUint(wallet, 10, 64)
		filter.Wallet = uint16(wallet64)
	}
	if to := u.Get("to"); to != "" {
		to64, _ := strconv.ParseUint(to, 10, 64)
		filter.To = uint16(to64)
	}
	const timeForm = "Jan-02-2006"
	if minDate := u.Get("min_date"); minDate != "" {
		filter.MinDate, _ = time.Parse(timeForm, minDate)
	}
	if maxDate := u.Get("max_date"); maxDate != "" {
		filter.MaxDate, _ = time.Parse(timeForm, maxDate)
	}
	return filter
}
