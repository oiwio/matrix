package db

import (
	"fmt"
	"reflect"

	"gopkg.in/mgo.v2"
)

var (
	Collections map[string]string
)

func Collection(s *mgo.Session, m interface{}) *mgo.Collection {
	typ := reflect.TypeOf(m).Elem()
	n := typ.Name()
	return s.DB("aladdin").C(fmt.Sprintf("%s", n))
}
