package main

import (
	"elastic/handler"
	"elastic/l"
	"elastic/store"
	"net/http"

	"github.com/go-chi/chi/v5"
	render_chi "github.com/go-chi/render"
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
	// m := martini.Classic()
	// m.Use(render.Renderer())
	// //Routes
	// m.Get("/article/id/:id", articleHandler.Id)
	// m.Post("/article/add", articleHandler.Add)
	// m.Post("/article/search", articleHandler.Search)
	// panicHandler := handler.PanicHandler{}
	// m.Get("/panic", panicHandler.Handle)
	// m.Post("/log/add", panicHandler.Log)
	// m.Run()

	// chi
	r := chi.NewRouter()
	r.Use(render_chi.SetContentType(render_chi.ContentTypeJSON))

	//routes
	r.Get("Get", articleHandler.Id_chi)
	r.Post("/article/add", articleHandler.Add_chi)
	r.Post("/article/search", articleHandler.Search_chi)
	//panic
	panicHandler := handler.PanicHandler{}
	r.Get("/panic", panicHandler.Handle_chi)
	r.Post("/log/add", panicHandler.Log_chi)

	http.ListenAndServe(":3333", r)
}

func parseErr(err error) {
	if err != nil {
		l.F(err)
	}
	l.Log.Log("Application started")
}
