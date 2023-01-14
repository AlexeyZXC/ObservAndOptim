package main

import (
	"elastic/handler"
	"elastic/l"
	"elastic/store"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

// Переписать не на Martini
func main() {
	//Sentry error handler
	//sentry.Init(sentry.Client(os.Getenv("SENTRY_DSN")))
	//Initialize Stores
	articleStore, err := store.NewArticleStore()
	parseErr(err)
	//Initialize Handlers
	articleHandler := handler.NewArticleHandler(articleStore)
	//Initialize Router
	m := martini.Classic()
	m.Use(render.Renderer())
	//Routes
	m.Get("/article/id/:id", articleHandler.Id)
	m.Post("/article/add", articleHandler.Add)
	m.Post("/article/search", articleHandler.Search)
	panicHandler := handler.PanicHandler{}
	m.Get("/panic", panicHandler.Handle)
	m.Post("/log/add", panicHandler.Log)
	m.Run()
}

func parseErr(err error) {
	if err != nil {
		l.F(err)
	}
	l.Log.Log("Application started")
}
