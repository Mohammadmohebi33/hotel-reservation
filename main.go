package main

import (
	"context"
	"flag"
	"github.com/Mohammadmohebi33/hotel-reservation/api/middleware"
	"github.com/Mohammadmohebi33/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"

	"github.com/Mohammadmohebi33/hotel-reservation/api"
	"github.com/gofiber/fiber/v2"
)

const dburi = "mongodb://localhost:27017"

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {

	listenAddr := flag.String("listenAddr", ":8080", "the listen address of api server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal("errro")
	}

	app := fiber.New(config)
	auth := app.Group("/api")
	apiv1 := app.Group("/api/v1", middleware.JWTAuthentication)

	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))
	hotelStore := db.NewMongoHotelStore(client, "hotel")
	roomStore := db.NewMongoRoomStore(client, "hotel", hotelStore)
	userStore := db.NewMongoUserStore(client)
	authHandler := api.NewAuthHandler(userStore)
	store := db.Store{
		User:  userStore,
		Hotel: hotelStore,
		Room:  roomStore,
	}
	hotelHandler := api.NewHotelHandler(store)

	//auth
	auth.Post("/auth", authHandler.HandleAuthentication)
	//users
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("user/:id", userHandler.HandleDeleteUser)
	apiv1.Post("/user", userHandler.HandlPostUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)

	//hotels handlers
	apiv1.Get("/hotel", hotelHandler.HandlerGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandlerGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)
	app.Listen(*listenAddr)
}
