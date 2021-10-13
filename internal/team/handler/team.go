package handler

import (
	"encoding/json"
	"fmt"
	"github.com/batroff/todo-back/internal/models"
	"github.com/batroff/todo-back/internal/team"
	relMaker "github.com/batroff/todo-back/internal/team_relation_maker"
	"github.com/batroff/todo-back/internal/user"
	"github.com/batroff/todo-back/pkg/handler"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
	"reflect"
)

const teamRoute = "/teams"
const teamUsersRoute = "/users"

func teamCreateHandler(teamCase team.UseCase) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		responseWriter := handler.NewResponseWriter(rw)

		t := new(models.Team)
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			responseWriter.Write(http.StatusBadRequest, fmt.Errorf("%s: %s", models.ErrBadRequest, err))
			return
		}

		if err := teamCase.CreateTeam(t); err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}

		rw.Header().Set("Location", handler.MakeURI(r.RequestURI, t.ID.String()))
		rw.WriteHeader(http.StatusCreated)
	})
}

func teamListHandler(teamCase team.UseCase) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		responseWriter := handler.NewResponseWriter(rw)
		responseWriter.SetHeaders(map[string]string{
			"Cache-Control": "no-store, no-cache, must-revalidate",
			"Content-Type":  "application/json; charset=utf-8",
			"Pragma":        "no-cache",
		})

		// TODO : implement filtered list
		//if err := r.ParseForm(); err != nil {
		//	responseWriter.Write(http.StatusBadRequest, fmt.Errorf("%s: %s", models.ErrBadRequest, err))
		//	return
		//}

		teams, err := teamCase.SelectTeamsList()
		if err == models.ErrNotFound {
			responseWriter.Write(http.StatusOK, make([]string, 0))
			return
		} else if err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}

		responseWriter.Write(http.StatusOK, teams)
	})
}

func teamGetHandler(teamCase team.UseCase) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		responseWriter := handler.NewResponseWriter(rw)

		id, err := handler.GetIDFromURI(r)
		if err != nil {
			responseWriter.Write(http.StatusBadRequest, fmt.Errorf("%s: %s", models.ErrBadRequest, err))
			return
		}

		t, err := teamCase.SelectTeamByID(id)
		if err == models.ErrNotFound {
			responseWriter.Write(http.StatusNotFound, err)
			return
		} else if err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}

		if err := json.NewEncoder(rw).Encode(&t); err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}
	})
}

// TODO: add tests
func teamDeleteHandler(teamCase team.UseCase, relCase relMaker.UseCase) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		responseWriter := handler.NewResponseWriter(rw)

		id, err := handler.GetIDFromURI(r)
		if err != nil {
			responseWriter.Write(http.StatusBadRequest, fmt.Errorf("%s: %s", models.ErrBadRequest, err))
			return
		}

		if _, err := relCase.SelectRelationsByTeamID(id); err != models.ErrNotFound && err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		} else if err == nil {
			if err = relCase.DeleteRelationsByTeamID(id); err != nil {
				responseWriter.Write(http.StatusInternalServerError, err)
				return
			}
		}

		if err := teamCase.DeleteTeam(id); err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}

		rw.WriteHeader(http.StatusNoContent)
	})
}

// TODO : add tests
func teamUpdateHandler(teamCase team.UseCase) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		responseWriter := handler.NewResponseWriter(rw)

		id, err := handler.GetIDFromURI(r)
		if err != nil {
			responseWriter.Write(http.StatusBadRequest, fmt.Errorf("%s: %s", models.ErrBadRequest, err))
			return
		}

		var reqTeam models.RequestTeam
		if err := json.NewDecoder(r.Body).Decode(&reqTeam); err != nil {
			responseWriter.Write(http.StatusBadRequest, fmt.Errorf("%s: %s", models.ErrBadRequest, err))
			return
		}

		t, err := teamCase.SelectTeamByID(id)
		if err == models.ErrNotFound {
			responseWriter.Write(http.StatusNotFound, models.ErrNotFound)
			return
		} else if err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}

		refReq := reflect.ValueOf(reqTeam)
		refUpd := reflect.ValueOf(t).Elem()
		handler.ParseRequestFields(refReq, refUpd)

		if err := teamCase.UpdateTeam(t); err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}

		// Encode response
		if err := json.NewEncoder(rw).Encode(&reqTeam); err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}
	})
}

// TODO : add tests
func teamAddUserHandler(relCase relMaker.UseCase, userCase user.UseCase, teamCase team.UseCase) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		responseWriter := handler.NewResponseWriter(rw)

		ids, err := handler.ParseQueryIDs(r)
		if err != nil {
			responseWriter.Write(http.StatusBadRequest, fmt.Errorf("%s: %s", models.ErrBadRequest, err))
			return
		}

		// Find teamID in db
		teamID, ok := ids["teamID"]
		if !ok {
			responseWriter.Write(http.StatusBadRequest, fmt.Errorf("%s: %s", models.ErrBadRequest, "teamID wasn't found in query"))
			return
		}

		if _, err := teamCase.SelectTeamByID(teamID); err == models.ErrNotFound {
			responseWriter.Write(http.StatusNotFound, err)
			return
		} else if err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}

		// Find userID in db
		userID, ok := ids["userID"]
		if !ok {
			responseWriter.Write(http.StatusBadRequest, fmt.Errorf("%s: %s", models.ErrBadRequest, "userID wasn't found in query"))
			return
		}

		if _, err := userCase.GetUser(userID); err == models.ErrNotFound {
			responseWriter.Write(http.StatusNotFound, err)
			return
		} else if err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}

		// Add user to team
		relation := models.NewUserTeamRel(userID, teamID)
		if err := relCase.CreateRelation(relation); err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}

		rw.WriteHeader(http.StatusCreated)
	})
}

// TODO : add tests
func teamDeleteUserHandler(relCase relMaker.UseCase) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		responseWriter := handler.NewResponseWriter(rw)

		ids, err := handler.ParseQueryIDs(r)
		if err != nil {
			responseWriter.Write(http.StatusBadRequest, fmt.Errorf("%s: %s", models.ErrBadRequest, err))
			return
		}

		// Parse teamID
		teamID, ok := ids["teamID"]
		if !ok {
			responseWriter.Write(http.StatusBadRequest, fmt.Errorf("%s: %s", models.ErrBadRequest, "teamID wasn't found in query"))
			return
		}

		// Parse userID
		userID, ok := ids["userID"]
		if !ok {
			responseWriter.Write(http.StatusBadRequest, fmt.Errorf("%s: %s", models.ErrBadRequest, "userID wasn't found in query"))
			return
		}

		// Delete relation
		if err := relCase.DeleteRelationByIDs(teamID, userID); err != nil {
			responseWriter.Write(http.StatusInternalServerError, err)
			return
		}

		rw.WriteHeader(http.StatusNoContent)
	})
}

func MakeTeamHandlers(r *mux.Router, n negroni.Negroni,
	teamCase team.UseCase, relCase relMaker.UseCase, userCase user.UseCase) {
	// TEAM OPERATIONS /teams/{id}/users/
	// Add user to team
	r.Handle(handler.MakeURI(teamRoute, "{teamID}", teamUsersRoute, "{userID}"), n.With(
		negroni.Wrap(teamAddUserHandler(relCase, userCase, teamCase)),
	)).Methods("PUT").
		Name("TeamAddUserHandler")

	// Delete user from team
	r.Handle(handler.MakeURI(teamRoute, "{teamID}", teamUsersRoute, "{userID}"), n.With(
		negroni.Wrap(teamDeleteUserHandler(relCase)),
	)).Methods("DELETE").
		Name("TeamDeleteUserHandler")

	// TEAM OPERATIONS /teams
	// Create
	r.Handle(teamRoute, n.With(
		negroni.Wrap(teamCreateHandler(teamCase)),
	)).Methods("POST").
		Headers("Content-Type", "application/json").
		Name("TeamCreateHandler")

	// List
	r.Handle(teamRoute, n.With(
		negroni.Wrap(teamListHandler(teamCase)),
	)).Methods("GET").
		Name("TeamListHandler")

	// Get
	r.Handle(handler.MakeURI(teamRoute, "{id}"), n.With(
		negroni.Wrap(teamGetHandler(teamCase)),
	)).Methods("GET").
		Name("TeamGetHandler")

	// Delete
	r.Handle(handler.MakeURI(teamRoute, "{id}"), n.With(
		negroni.Wrap(teamDeleteHandler(teamCase, relCase)),
	)).Methods("DELETE").
		Name("TeamDeleteHandler")

	// Update
	r.Handle(handler.MakeURI(teamRoute, "{id}"), n.With(
		negroni.Wrap(teamUpdateHandler(teamCase)),
	)).Methods("PATCH").
		Headers("Content-Type", "application/json").
		Name("TeamUpdateHandler")
}
