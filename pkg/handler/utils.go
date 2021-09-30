package handler

import (
	"github.com/batroff/todo-back/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

func MakeURI(uriParams ...string) string {
	rawURI := strings.Join(uriParams, "/")
	re := regexp.MustCompile("/{2,}")
	return re.ReplaceAllString(rawURI, "/")
}

func GetIDFromURI(r *http.Request) (models.ID, error) {
	vars := mux.Vars(r)
	return uuid.Parse(vars["id"])
}

func ParseRequestFields(refReq, refUpd reflect.Value) {
	t := refReq.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		if refReq.Field(i).IsZero() && refUpd.FieldByName(f.Name).Kind() != reflect.Ptr {
			continue
		}

		var v reflect.Value

		switch refUpd.FieldByName(f.Name).Type().Kind() {
		case reflect.Ptr:
			v = refReq.Field(i)
		default:
			v = refReq.Field(i).Elem()
		}

		refUpd.FieldByName(f.Name).Set(v)
	}
}
