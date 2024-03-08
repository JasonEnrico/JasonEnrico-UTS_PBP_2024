package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	m "uts/models"

	"github.com/gorilla/mux"
)

func GetAllRooms(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	rows, err := db.Query("SELECT id, room_name FROM rooms")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(m.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Internal Server Error",
		})
		return
	}
	defer rows.Close()

	var rooms []m.Room
	for rows.Next() {
		var room m.Room
		if err := rows.Scan(&room.ID, &room.RoomName); err != nil {
			// Send Error Response
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(m.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "Internal Server Error",
			})
			return
		} else {
			rooms = append(rooms, room)
		}
	}

	response := m.RoomResponse{
		Status: http.StatusOK,
		Data: m.RoomList{
			Rooms: rooms,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetDetailRooms(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	params := mux.Vars(r)
	roomID := params["id"]

	query := "SELECT rooms.id, rooms.room_name, participants.id AS participant_id, participants.id_account, accounts.username FROM rooms LEFT JOIN participants ON rooms.id = participants.id_room LEFT JOIN accounts ON participants.id_account = accounts.id WHERE rooms.id = ?"
	rows, err := db.Query(query, roomID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(m.RoomDetailResponse{
			Status: http.StatusInternalServerError,
			Data:   nil,
		})
		return
	}
	defer rows.Close()

	room := &m.RoomDetail{}

	for rows.Next() {
		var participant m.Participant
		var username sql.NullString

		err := rows.Scan(&room.ID, &room.RoomName, &participant.ID, &participant.AccountID, &username)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(m.RoomDetailResponse{
				Status: http.StatusInternalServerError,
				Data:   nil,
			})
			return
		}

		if username.Valid {
			participant.Username = username.String
		}

		room.Participants = append(room.Participants, participant)
	}

	response := m.RoomDetailResponse{
		Status: http.StatusOK,
		Data:   room,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func InsertRoom(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	var request m.InsertRoomRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(m.InsertRoomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
		})
		return
	}

	var maxPlayer int
	var gameID int
	err = db.QueryRow("SELECT id_game FROM rooms WHERE id = ?", request.RoomID).Scan(&gameID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(m.InsertRoomResponse{
			Status:  http.StatusNotFound,
			Message: "Room not found",
		})
		return
	}

	err = db.QueryRow("SELECT max_player FROM games WHERE id = ?", gameID).Scan(&maxPlayer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(m.InsertRoomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Internal Server Error",
		})
		return
	}

	var currentParticipants int
	err = db.QueryRow("SELECT COUNT(*) FROM participants WHERE id_room = ?", request.RoomID).Scan(&currentParticipants)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(m.InsertRoomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Internal Server Error",
		})
		return
	}

	if currentParticipants >= maxPlayer {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(m.InsertRoomResponse{
			Status:  http.StatusBadRequest,
			Message: "Room is full",
		})
		return
	}

	_, err = db.Exec("INSERT INTO participants (id_room, id_account) VALUES (?, ?)", request.RoomID, request.AccountID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(m.InsertRoomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to insert into room",
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(m.InsertRoomResponse{
		Status:  http.StatusCreated,
		Message: "Successfully inserted into room",
	})
}

func LeaveRoom(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	params := mux.Vars(r)
	participantID := params["id"]

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM participants WHERE id = ?", participantID).Scan(&count)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(m.LeaveRoomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Internal Server Error",
		})
		return
	}

	if count == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(m.LeaveRoomResponse{
			Status:  http.StatusNotFound,
			Message: "Participant not found",
		})
		return
	}

	_, err = db.Exec("DELETE FROM participants WHERE id = ?", participantID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(m.LeaveRoomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to leave room",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(m.LeaveRoomResponse{
		Status:  http.StatusOK,
		Message: "Successfully left room",
	})
}
