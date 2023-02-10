package main

import (
	"elastic/handler"
	"elastic/l"
	"elastic/store"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	render_chi "github.com/go-chi/render"
)

// Переписать не на Martini
func main() {
	// elastic logger
	elLogger, err := l.NewElasticLogger()
	if err != nil {
		fmt.Println("NewElasticLogger error: ", err)
		panic(err)
	}
	elLogger.Log("Starting...")

	//Sentry error handler
	//sentry.Init(sentry.Client(os.Getenv("SENTRY_DSN")))

	//Initialize Stores
	articleStore, err := store.NewArticleStore(elLogger)
	if err != nil {
		elLogger.Error("NewArticleStore creation failed: %v", err)
		panic(err)
	}

	//Initialize Handlers
	articleHandler := handler.NewArticleHandler(articleStore)

	// chi
	r := chi.NewRouter()
	r.Use(render_chi.SetContentType(render_chi.ContentTypeJSON))

	//routes
	//r.Get("Get", articleHandler.Id_chi)
	r.Post("/article/add", articleHandler.Add_chi)
	r.Post("/article/search", articleHandler.Search_chi)

	//panic
	panicHandler := handler.PanicHandler{}
	r.Get("/panic", panicHandler.Handle_chi)
	r.Post("/log/add", panicHandler.Log_chi)

	//listen
	elLogger.Log("Application started %v", 1)
	http.ListenAndServe(":3333", r)
	elLogger.Log("Application stopped %v", 2)
}
