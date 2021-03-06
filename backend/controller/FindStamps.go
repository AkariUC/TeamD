// Written by Yuto Kojima

package controller

import (
    "fmt"
    "net/http"
    "strconv"

    "github.com/jmoiron/sqlx"
    "github.com/gorilla/mux"
    "github.com/shortintern2020-B-frontier/TeamD/model"
    "github.com/shortintern2020-B-frontier/TeamD/repository"
)

type Stamp struct {
    db *sqlx.DB
}

// controller
func NewStamp(db *sqlx.DB) *Stamp{
    return &Stamp{db: db}
}

// find stamps that mutch the conditions
func (stamp *Stamp) FindStamps(w http.ResponseWriter, r *http.Request) (int, interface{}, error){

    // convert int -> string
    ellapsed_time, _ := strconv.Atoi(r.FormValue("ellapsed_time"))
    vars := mux.Vars(r)
    room_id, _ := strconv.Atoi(vars["room_id"])

    // throw error when the rooms does't exist
    rooms := make([]model.RoomTemp, 0)
    if err := stamp.db.Select(&rooms, "select id, end_time from room where id = ?", room_id); err != nil {
        fmt.Printf("%s", err)
        return http.StatusInternalServerError, nil, err
    }

    switch {
        // no rooms
        case len(rooms) <= 0:
            return http.StatusNotFound, nil, nil

        // two or more rooms exist
        case len(rooms) >= 2:
            return http.StatusInternalServerError, nil, nil
    }
    
    room := rooms[0]
    _, end_time := room.RoomID, room.EndTime

    stamps, err := repository.SelectStamps(stamp.db, ellapsed_time, room_id, end_time)

    if err != nil {
        return http.StatusInternalServerError, nil, err
    }

    return http.StatusOK, stamps, nil
}
