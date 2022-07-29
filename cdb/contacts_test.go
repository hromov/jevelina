package cdb

import (
	"log"
	"testing"

	"github.com/hromov/jevelina/cdb/contacts"
	"github.com/hromov/jevelina/cdb/models"
)

func BenchmarkContacts(b *testing.B) {
	db, err := OpenTest()
	if err != nil {
		log.Fatalf("Cant open data base error: %s", err.Error())
	}
	contacts := &contacts.Contacts{DB: db.DB}
	for i := 0; i < b.N; i++ {
		_, err := contacts.List(models.ListFilter{Limit: 50, Offset: 0, Query: ""})
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkContactsPhoneSearch(b *testing.B) {
	db, err := OpenTest()
	if err != nil {
		log.Fatalf("Cant open data base error: %s", err.Error())
	}
	contacts := &contacts.Contacts{DB: db.DB}
	for i := 0; i < b.N; i++ {
		_, err := contacts.List(models.ListFilter{Limit: 50, Offset: 0, Query: "067"})
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkContactsNameSearch(b *testing.B) {
	db, err := OpenTest()
	if err != nil {
		log.Fatalf("Cant open data base error: %s", err.Error())
	}
	contacts := &contacts.Contacts{DB: db.DB}
	for i := 0; i < b.N; i++ {
		_, err := contacts.List(models.ListFilter{Limit: 50, Offset: 0, Query: "Петр"})
		if err != nil {
			panic(err)
		}
	}
}
