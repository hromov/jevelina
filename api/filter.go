package api

import (
	"net/url"
	"strconv"

	"github.com/hromov/cdb/models"
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
	if tagID := u.Get("tag"); tagID != "" {
		tag64, _ := strconv.ParseUint(tagID, 10, 64)
		filter.TagID = uint8(tag64)
	}
	return filter
}
