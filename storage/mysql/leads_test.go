package mysql

import (
	"log"
	"testing"

	"github.com/hromov/jevelina/storage/mysql/dao/leads"
	"github.com/hromov/jevelina/storage/mysql/dao/models"
)

func BenchmarkLeads(b *testing.B) {
	db, err := OpenTest()
	if err != nil {
		log.Fatalf("Can't open DB error")
	}
	l := &leads.Leads{DB: db.DB}
	for i := 0; i < b.N; i++ {
		_, err := l.List(models.ListFilter{Limit: 50, Offset: 0, Query: ""})
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkActiveLeads(b *testing.B) {
	db, err := OpenTest()
	if err != nil {
		log.Fatalf("Can't open DB error")
	}
	l := &leads.Leads{DB: db.DB}
	for i := 0; i < b.N; i++ {
		leadList, err := l.List(models.ListFilter{Active: true, Limit: 50, Offset: 0, Query: ""})
		if err != nil {
			panic(err)
		}
		for _, lead := range leadList.Leads {
			if lead.ClosedAt != nil {
				log.Fatalf("should return only active leads")
			}
		}
	}
}
