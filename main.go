package main

import (
	"net/http"

	"github.com/go-redis/redis"
	"github.com/go-zoo/bone"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var (
	rd *render.Render
	db *redis.Client
)

func init() {
	// JSON response object
	rd = render.New()

	// init redis db
	db = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func main() {
	// router
	r := bone.New()

	buildRoutes(r)
	n := negroni.New()

	n.UseHandler(r)

	n.Run(":3000")

}

func buildRoutes(r *bone.Mux) {
	r.Get("/route/:token", http.HandlerFunc(GetRoute))
	r.Post("/route", http.HandlerFunc(CreateRoute))
}
