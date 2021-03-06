// Written by Taishi
// audience controller
package controller

import (
	"net/http"
	"strconv"
	"errors"
	"database/sql"
	"encoding/json"

	"github.com/jmoiron/sqlx"
	"github.com/gorilla/mux"

    "github.com/shortintern2020-B-frontier/TeamD/repository"
)

type Audience struct {
	db *sqlx.DB
}

//constructor
func NewAudience(db *sqlx.DB) *Audience {
	return &Audience{db: db}
}

//return audience.count from audience table
func (audience *Audience) FindAudience(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	// get parameters from url
	ellapsed_time, err := strconv.Atoi(r.FormValue("ellapsed_time"))
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	if ellapsed_time % 60000 != 0 {
		return http.StatusBadRequest, nil, err
	}
    vars := mux.Vars(r)
	room_id, err := strconv.Atoi(vars["room_id"])
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	//check if room exists
	if _, err := repository.FindRoomDB(audience.db, room_id); err != nil {
		return http.StatusNotFound, nil, err
	}

	//check room.end_time
	endtime, err := repository.FindRoomEndTimeDB(audience.db, room_id)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	if endtime != nil && *endtime < ellapsed_time {
		return http.StatusNotFound, nil, errors.New("page not found")
	}

	//retrieve audience.count from database
	a, err := repository.FindAudience(audience.db, room_id, ellapsed_time)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, a, nil
}

func (audience *Audience) UpdateAudience(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	// get parameters
	type EllapsedTime struct{
		Time int `json:"ellapsed_time"`
	}
	ellapsed_time := &EllapsedTime{}
	err := json.NewDecoder(r.Body).Decode(&ellapsed_time)
	if err != nil {
        return http.StatusBadRequest, nil, err
	}
	if ellapsed_time.Time % 60000 != 0 {
		return http.StatusBadRequest, nil, err
	}
    vars := mux.Vars(r)
	room_id, err := strconv.Atoi(vars["room_id"])
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	//check if room exists
	if _, err := repository.FindRoomDB(audience.db, room_id); err != nil {
		return http.StatusNotFound, nil, err
	}

	//check room.end_time
	endtime, err := repository.FindRoomEndTimeDB(audience.db, room_id)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	if endtime != nil && *endtime < ellapsed_time.Time {
		return http.StatusNotFound, nil, errors.New("page not found")
	}

	//retrieve audience.count from database
	a, err := repository.FindAudience(audience.db, room_id, ellapsed_time.Time)
	var result sql.Result
	if err != nil {
		return http.StatusInternalServerError, nil, nil
	} else {
		result, err = repository.UpdateAudience(audience.db, a.Audience+1, room_id, ellapsed_time.Time)
	}
	rows, _ := result.RowsAffected()
	return http.StatusOK, rows, nil
}
