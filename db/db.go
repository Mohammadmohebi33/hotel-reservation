package db

type Store struct {
	User  UserStore
	Hotel HotelStore
	Room  RoomStore
}
