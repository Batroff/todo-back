package handler

import (
	"fmt"
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
	uri := re.ReplaceAllString(rawURI, "/")
	if uri[len(uri)-1] == '/' {
		uri = uri[:len(uri)-1]
	}

	return uri
}

func MakeQuery(uri string, queryParams map[string]interface{}) string {
	if len(queryParams) == 0 {
		return uri
	}

	query := make([]string, len(queryParams))
	for k, v := range queryParams {
		query = append(query, fmt.Sprintf("%s=%v", k, v))
	}

	return uri + "?" + strings.Join(query, "&")
}

// GetIDFromURI : deprecated - TODO: replace with ParseQueryIDs
func GetIDFromURI(r *http.Request) (models.ID, error) {
	vars := mux.Vars(r)
	return uuid.Parse(vars["id"])
}

func ParseQueryIDs(r *http.Request) (map[string]models.ID, error) {
	vars := mux.Vars(r)

	ids := make(map[string]models.ID, 0)
	for k, v := range vars {
		if !strings.Contains(strings.ToLower(k), "id") {
			continue
		}

		id, err := uuid.Parse(v)
		if err != nil {
			return nil, err
		}

		ids[k] = id
	}

	return ids, nil
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
