package handler

import (
	"github.com/batroff/todo-back/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

const UUIDRegex = "{id:\\w{8}-(?:\\w{4}-){3}\\w{12}}"

func MakeRegexURI(uriParams ...string) string {
	return strings.Join(uriParams, "/")
}

func GetIDFromURI(r *http.Request) (models.ID, error) {
	vars := mux.Vars(r)
	return uuid.Parse(vars["id"])
}
