package analytics

import (
	"context"

	"github.com/hromov/jevelina/domain/leads"
)

type Service interface {
	LeadsBySource(ctx context.Context, filter Filter) ([]LeadsBySource, error)
}

type service struct {
	ls leads.Service
}

func NewService(ls leads.Service) *service {
	return &service{ls}
}

func (s *service) LeadsBySource(ctx context.Context, filter Filter) ([]LeadsBySource, error) {
	leadsResponse, err := s.ls.List(ctx, filter.toLeadsFilter())
	if err != nil {
		return nil, err
	}

	resMap := map[string]int{}
	for _, l := range leadsResponse.Leads {
		if l.Analytics.Domain != "" {
			resMap[l.Analytics.Domain] += 1
		} else {
			resMap[l.Source.Name] += 1
		}
	}
	res := make([]LeadsBySource, 0)
	for k, v := range resMap {
		res = append(res, LeadsBySource{Source: k, Count: v})
	}

	return res, nil
}
