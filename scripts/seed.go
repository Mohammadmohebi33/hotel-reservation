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

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client, db.Dbname)
	roomStore := db.NewMongoRoomStore(client, db.Dbname, hotelStore)

	hotel := types.Hotel{
		Name:     "test",
		Location: "thran",
		Rooms:    []primitive.ObjectID{},
	}

	insertedHotel, err := hotelStore.Insert(ctx, &hotel)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(insertedHotel)

	rooms := []types.Room{
		{
			Type:      types.SingleRoomType,
			BasePrice: 120.1,
		},
		{
			Type:      types.DoubleRoomType,
			BasePrice: 88,
		},
		{
			Type:      types.DoubleRoomType,
			BasePrice: 90,
		},
		{
			Type:      types.SingleRoomType,
			BasePrice: 200,
		},
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
