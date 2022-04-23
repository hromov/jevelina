package amoimport

import (
	"crypto/sha1"
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/hromov/cdb/contacts"
	"github.com/hromov/jevelina/base"
)

var mysqlErr *mysql.MySQLError

//key = hash, val = id
var contactsMap map[string]uint = map[string]uint{}

func hashIt(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return string(bs)
}

func Push_Contacts(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return errors.New("Unable to read input file " + path + ". Error: " + err.Error())
	}
	defer f.Close()

	r := csv.NewReader(f)
	for i := 0; i < 10000; i++ {
		record, err := r.Read()
		// Stop at EOF.
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}
		// Display record.
		// ... Display record length.
		// ... Display all individual elements of the slice.

		// fmt.Println(record)
		// fmt.Println(len(record))
		// for value := range record {
		// 	fmt.Printf(" %d = %v\n", value, record[value])
		// }
		db := base.GetDB()
		if contact := recordToContact(record); contact != nil {
			if _, err := db.Create(contact); err != nil {
				if !errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
					log.Printf("Can't create contact for record # = %d error: %s", i, err.Error())
				}
			}

			// } else {
			// 	log.Printf("contacts for record # = %d created: %+v", i, c)
			// }
			//notices 1-5, fullname, contact responsible, records[21:30], records[30:44]
			str := record[2] + record[6] + strings.Join(record[18:27], ",") + strings.Join(record[29:43], ",")
			// log.Println(str)
			hashed := hashIt(str)
			if _, exist := contactsMap[hashed]; exist {
				log.Println("WTF!!!!!!! contact exist with hash = ", hashed)
				log.Println(contact)
			} else {
				contactsMap[hashed] = contact.ID
			}
		}

	}
	return nil

	// csvReader := csv.NewReader(f)
	// records, err := csvReader.ReadAll()
	// if err != nil {
	// 	return errors.New("Error parsing file: " + err.Error())
	// }

	// return records
}

func recordToContact(record []string) *contacts.Contact {
	if len(record) == 0 {
		return nil
	}
	if len(record) != 43 {
		log.Println("Wrong record schema? len(record) = ", len(record))
		log.Println(record)
		return nil
	}
	contact := &contacts.Contact{}
	id, err := strconv.ParseUint(record[0], 10, 64)
	if err != nil || id == 0 {
		log.Println("ID parse error: " + err.Error())
		return nil
	}
	contact.ID = uint(id)
	if record[1] == "контакт" {
		contact.IsPerson = true
	}
	if record[3] != "" {
		contact.Name = record[3]
	} else {
		contact.Name = record[2]
	}
	contact.SecondName = record[4]
	if !contact.IsPerson && record[5] != "" {
		contact.Name = record[5]
	}
	//implement real user by get func
	contact.ResponsibleID = nil
	contact.CreatedID = nil

	const timeForm = "02-01-2006 15:04:05"
	contact.CreatedAt, _ = time.Parse(timeForm, record[7])

	//contact.tags = getTags
	//contact.notices = getNotices record[13:18]

	// Phones start
	dc := regexp.MustCompile(`[^\d|,]`)
	str := dc.ReplaceAllString(strings.Join(record[18:24], ","), "")
	digits := regexp.MustCompile(`(\d){6,13}`)
	// log.Println(str)
	phones := digits.FindAllString(str, -1)
	// log.Println(phones, len(phones))
	switch len(phones) {
	case 0:
		log.Printf("no phones found for contact: %d\n", contact.ID)
		break
	case 1:
		contact.Phone = phones[0]
	default:
		contact.Phone = phones[0]
		contact.SecondPhone = strings.Join(phones[1:], ",")
	}
	// Phones End

	// Email start
	mx := regexp.MustCompile(`[\w-\.]+@([\w-]+\.)+[\w-]{2,4}`)
	emails := mx.FindAllString(strings.Join(record[24:27], ","), -1)
	switch len(emails) {
	case 0:
		break
	case 1:
		contact.Email = emails[0]
	default:
		contact.SecondEmail = strings.Join(emails[1:], ",")
	}
	// Email end
	contact.URL = record[27]
	contact.Address = record[28]
	contact.City = record[29]

	// implements real source
	contact.SourceID = nil

	contact.Position = record[31]

	contact.Analytics.CID = record[40]
	contact.Analytics.UID = record[41]
	contact.Analytics.TID = record[42]

	// log.Printf("all ok: %+v", contact)
	return contact
}

// 0 = ID
//  1 = Тип
//  2 = Полное имя контакта
//  3 = Имя
//  4 = Фамилия
//  5 = Название компании
//  6 = Ответственный
//  7 = Дата создания контакта
//  8 = Кем создан контакт
//  9 = Сделки
//  10 = Дата редактирования
//  11 = Кем редактирован
//  12 = Теги
//  13 = Примечание 1
//  14 = Примечание 2
//  15 = Примечание 3
//  16 = Примечание 4
//  17 = Примечание 5
//  18 = Рабочий телефон
//  19 = Рабочий прямой телефон
//  20 = Мобильный телефон
//  21 = Факс
//  22 = Домашний телефон
//  23 = Другой телефон
//  24 = Рабочий email
//  25 = Личный email
//  26 = Другой email
//  27 = Web
//  28 = Адрес
//  29 = Город
//  30 = Источник
//  31 = Должность
//  32 = Товар
//  33 = Skype
//  34 = ICQ
//  35 = Jabber
//  36 = Google Talk
//  37 = MSN
//  38 = Другой IM
//  39 = Пользовательское соглашение
//  40 = cid
//  41 = uid
//  42 = tid

//  21 = Рабочий телефон
//  22 = Рабочий прямой телефон
//  23 = Мобильный телефон
//  24 = Факс
//  25 = Домашний телефон
//  26 = Другой телефон
//  27 = Рабочий email
//  28 = Личный email
//  29 = Другой email
//  30 = Город
//  31 = Источник
//  32 = Должность
//  33 = Товар
//  34 = Skype
//  35 = ICQ
//  36 = Jabber
//  37 = Google Talk
//  38 = MSN
//  39 = Другой IM
//  40 = Пользовательское соглашение
//  41 = cid
//  42 = uid
//  43 = tid
