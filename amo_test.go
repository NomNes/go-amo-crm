package go_amo_crm

import (
	"log"
	"os"
	"path"
	"runtime"
	"testing"

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

func TestAmoCrm_Contacts(t *testing.T) {
	Convey("Catalogs", t, func() {
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
					log.Println(item)
				})
			}
		})
	})
}
