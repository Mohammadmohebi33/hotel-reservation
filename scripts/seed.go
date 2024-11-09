package main

import (
	"context"
	"fmt"
	"github.com/Mohammadmohebi33/hotel-reservation/db"
	"github.com/Mohammadmohebi33/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const dburi = "mongodb://localhost:27017"

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	userStore  db.UserStore
	ctx        = context.Background()
)

func seedUser(fname, lname, email string) {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: fname,
		LastName:  lname,
		Email:     email,
		Password:  "password",
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = userStore.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
}
func seedHotel(name, location string, rating int) {

	hotel := types.Hotel{
		Name:     "test",
		Location: "thran",
		Rooms:    []primitive.ObjectID{},
		Rating:   4,
	}

	rooms := []types.Room{
		{
			Size:      "small",
			BasePrice: 120.1,
		},
		{
			Size:      "large",
			BasePrice: 88,
		},
		{
			Size:      "small",
			BasePrice: 90,
		},
		{
			Size:      "medium",
			BasePrice: 200,
		},
	}

	insertedHotel, err := hotelStore.Insert(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		insertedRoom, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(insertedRoom)
	}

}

func main() {
	seedHotel("hoteltest1", "thran", 5)
	seedHotel("hoteltest2", "rasht", 3)
	seedHotel("hoteltest3", "tabriz", 2)
	seedHotel("hoteltest4", "mashad", 1)

	seedUser("mohamamd", "mohebbi", "m@gmail.com")
}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.Dbname).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client, db.Dbname)
	roomStore = db.NewMongoRoomStore(client, db.Dbname, hotelStore)
	userStore = db.NewMongoUserStore(client)
}
