package models

type Room struct {
	ID       int    `json:"id"`
	RoomName string `json:"room_name"`
}

type RoomList struct {
	Rooms []Room `json:"rooms"`
}

type RoomResponse struct {
	Status int      `json:"status"`
	Data   RoomList `json:"data"`
}

type Participant struct {
	ID        int    `json:"id"`
	AccountID int    `json:"account_id"`
	Username  string `json:"username"`
}

type RoomDetail struct {
	ID           int           `json:"id"`
	RoomName     string        `json:"room_name"`
	Participants []Participant `json:"participants"`
}

type RoomDetailResponse struct {
	Status int         `json:"status"`
	Data   *RoomDetail `json:"data"`
}

type InsertRoomRequest struct {
	AccountID int `json:"account_id"`
	RoomID    int `json:"room_id"`
}

type InsertRoomResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type LeaveRoomResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
