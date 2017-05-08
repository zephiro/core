package store

import (
	"github.com/fernandez14/spartangeek-blacker/modules/mail"
	"gopkg.in/mgo.v2"
)

type Deps interface {
	Mgo() *mgo.Database
	Mailer() mail.Mailer
}

type Query func(col *mgo.Collection) *mgo.Query
