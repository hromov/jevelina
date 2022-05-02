package api

import (
	"net/url"
	"strconv"
	"time"

	"github.com/hromov/jevelina/cdb/models"
)

const defaultLimit = 50

func filterFromQuery(u url.Values) models.ListFilter {
	filter := models.ListFilter{}
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
