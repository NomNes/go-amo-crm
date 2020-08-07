package go_amo_crm

import (
	"os"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/joho/godotenv"
	. "github.com/smartystreets/goconvey/convey"
)

var amo *AmoCrm

func init() {
	_, file, _, _ := runtime.Caller(0)
	dir := path.Dir(file)
	if err := godotenv.Load(dir + "/test.env"); err != nil {
		panic("No .env file found")
	}
	amo = &AmoCrm{storage: &FileStorage{Path: dir + "/amo.storage"}}
	amo.Setup(os.Getenv("SUBDOMAIN"), os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"), os.Getenv("REDIRECT_URI"))
}

func auth() {
	Convey("Auth", func() {
		err := amo.Restore()
		if err != nil {
			err = amo.Auth(os.Getenv("AUTHORIZATION_CODE"))
		}
		So(err, ShouldBeNil)
	})
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestAmoCrm_Auth(t *testing.T) {
	Convey("Failed", t, func() {
		c := &AmoCrm{storage: &RuntimeStorage{}}
		c.Setup("", "", "", "")
		So(c.Auth(""), ShouldBeError)
	})
	Convey("Success", t, func() {
		auth()
	})
}

func TestAmoCrm_Fields(t *testing.T) {
	Convey("Fields", t, func() {
		auth()
		fields, err := amo.GetFields()
		So(err, ShouldBeNil)
		So(fields, ShouldNotBeNil)
	})
}

func TestAmoCrm_GetContacts(t *testing.T) {
	Convey("Contacts", t, func() {
		auth()
		Convey("get list", func() {
			contacts, err := amo.GetContacts()
			So(err, ShouldBeNil)
			So(contacts, ShouldNotBeNil)
			if len(contacts) > 0 {
				Convey("get item", func() {
					item, err := amo.GetContact(contacts[0].Id)
					So(err, ShouldBeNil)
					So(item, ShouldNotBeNil)
					So(item.Id, ShouldEqual, contacts[0].Id)
				})
			}
		})

		Convey("add item", func() {
			fields, err := amo.GetFields()
			if err != nil {
				panic(err)
			}
			var phoneField AddContactCustomField
			var emailField AddContactCustomField
			if fields.Contacts != nil {
				for _, field := range fields.Contacts {
					switch field.Code {
					case "PHONE":
						phoneField = AddContactCustomField{
							Id: field.Id,
							Values: []AddContactCustomFieldValue{
								{Value: "+79876543210", Enum: "WORK"},
							},
						}
					case "EMAIL":
						emailField = AddContactCustomField{
							Id: field.Id,
							Values: []AddContactCustomFieldValue{
								{Value: "test@example.com", Enum: "WORK"},
							},
						}
					}
				}
			}

			contact := AddContact{
				Name:         "Контакт Тестов",
				CustomFields: []AddContactCustomField{emailField},
			}

			ids, err := amo.AddContacts([]AddContact{contact})
			So(err, ShouldBeNil)
			So(ids, ShouldHaveLength, 1)
			if len(ids) > 0 {
				contact.Id = ids[0]
				contact.UpdatedAt = int(time.Now().Unix())
				contact.CustomFields = append(contact.CustomFields, phoneField)
				updatedIds, err := amo.UpdateContacts([]AddContact{contact})
				So(err, ShouldBeNil)
				So(updatedIds, ShouldHaveLength, 1)
			}
		})
	})
}
