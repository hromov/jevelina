package cdb

import (
	"log"
	"testing"
	"unicode"

	"github.com/hromov/jevelina/cdb/contacts"
	"github.com/hromov/jevelina/cdb/models"
)

func isInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

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

// func BenchmarkPhonesIsIntSearch(b *testing.B) {
// 	if err := Init(dsn); err != nil {
// 		log.Fatalf("Cant open data base error: %s", err.Error())
// 	}
// 	for i := 0; i < b.N; i++ {
// 		if !isInt("0674430") {
// 			panic("not numbers")
// 		}
// 		_, err := ContactsPhone(50, 0, "0674430")
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// }

// func BenchmarkPhonesRegexpSearch(b *testing.B) {
// 	if err := Init(dsn); err != nil {
// 		log.Fatalf("Cant open data base error: %s", err.Error())
// 	}
// 	for i := 0; i < b.N; i++ {
// 		var digitCheck = regexp.MustCompile(`^[0-9]+$`)
// 		if !digitCheck.MatchString("1212") {
// 			panic("not numbers")
// 		}
// 		_, err := ContactsPhone(50, 0, "0674430")
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// }

// func BenchmarkNamesAndPhonesSearch(b *testing.B) {
// 	if err := Init(dsn); err != nil {
// 		log.Fatalf("Cant open data base error: %s", err.Error())
// 	}
// 	for i := 0; i < b.N; i++ {
// 		_, err := ContactsNamesAndPhone(50, 0, "0674430")
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// }

// goos: linux
// goarch: amd64
// pkg: cdb
// BenchmarkFullSearch-12              	      14	  82883176 ns/op
// BenchmarkPhonesIsIntSearch-12       	      18	  64922847 ns/op
// BenchmarkPhonesRegexpSearch-12      	      19	  61434232 ns/op
// BenchmarkNamesAndPhonesSearch-12    	      18	  70388692 ns/op
// PASS
// ok  	cdb	7.025s
