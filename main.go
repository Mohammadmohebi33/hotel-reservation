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

	// stores
	hotelStore := db.NewMongoHotelStore(client, "hotel")
	roomStore := db.NewMongoRoomStore(client, "hotel", hotelStore)
	userStore := db.NewMongoUserStore(client)
	bookingStore := db.NewMongoBookingStore(client)

	store := &db.Store{
		User:    userStore,
		Hotel:   hotelStore,
		Room:    roomStore,
		Booking: bookingStore,
	}

	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))
	authHandler := api.NewAuthHandler(userStore)
	hotelHandler := api.NewHotelHandler(store)
	roomHandler := api.NewRoomHandler(store)
	bookHandler := api.NewBookHandler(store)

	auth := app.Group("/api")
	apiv1 := app.Group("/api/v1", middleware.JWTAuthentication(userStore))
	admin := apiv1.Group("/admin", middleware.AdminAuth)

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

	//booking
	apiv1.Get("/room", roomHandler.HandleGetRooms)
	apiv1.Post("/room/:id/booking", roomHandler.HandleRoomBook)
	admin.Get("/booking", bookHandler.HandleGetBookings)
	apiv1.Get("/booking/:id", bookHandler.HandleGetBooking)
	apiv1.Get("/booking/:id/cancel", bookHandler.HandleCancelBook)
	app.Listen(*listenAddr)
}
