package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"

	"github.com/colebaileygit/basic-api-server/orders"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)

	router.Mount("/orders", orders.Routes())

	return router	
}

func main() {
	router := Routes()

	// TODO: Update port to use ENV variable
	log.Fatal(http.ListenAndServe(":8080", router))
}