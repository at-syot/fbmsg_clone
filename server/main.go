package main

import (
	"log"
	"net/http"

	"github.com/at-syot/msg_clone/db"
	"github.com/at-syot/msg_clone/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
)

func main() {
	err := db.Init()
	if err != nil {
		log.Fatal(err.Error())
	}

	r := chi.NewRouter()
	r.Use(cors.Handler(
		cors.Options{
			AllowedOrigins: []string{"http://*"},
			AllowedMethods: []string{
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
				http.MethodOptions,
			},
		},
	))

	r.Post("/user", handlers.CreateUserHandler)
	r.Get("/users", handlers.GetUsersHandler)
	r.Get("/users/{uid}/channels", handlers.GetUserChannelsHandler)
	r.Post("/channels", handlers.CreateChannel)
	r.Get("/ws", handlers.WebsocketHandler)

	log.Println("server start pn port :3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
