package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"geekbrains/internal/models"
	"geekbrains/store"
	"net/http"
	"net/http/pprof"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DatabaseURL = "postgres://usr:pwd@localhost:5432/example?sslmode=disable"
)

type App struct {
	logger *zap.Logger
	// pool       *pgxpool.Pool
	repository Repository
}

func (a *App) parseUserID(r *http.Request) (*uuid.UUID, error) {
	strUserID := chi.URLParam(r, "id")
	if strUserID == "" {
		return nil, nil
	}
	userID, err := uuid.Parse(strUserID)
	if err != nil {
		a.logger.Debug(
			fmt.Sprintf("failed to parse userID (uuid) from: '%s'", strUserID),
			zap.Field{Key: "error", String: err.Error(), Type: zapcore.StringType},
		)
		return nil, err
	}
	a.logger.Debug(fmt.Sprintf("userID parsed: %s", userID))
	return &userID, nil
}

func (a *App) usersHandler(w http.ResponseWriter, r *http.Request) {
	a.logger.Info("usersHandler called", zap.Field{Key: "method", String: r.Method, Type: zapcore.StringType})
	ctx := r.Context()
	users, err := a.repository.GetUsers(ctx)
	a.logger.Info("users len: " + strconv.Itoa(len(users)))

	if err != nil {
		msg := fmt.Sprintf(`failed to get users: %s`, err)
		a.logger.Error(msg)
		writeResponse(w, http.StatusInternalServerError, msg)
		return
	}
	writeJsonResponse(w, http.StatusOK, users)
}

func (a *App) userHandler(w http.ResponseWriter, r *http.Request) {
	a.logger.Info("userHandler called", zap.Field{Key: "method", String: r.Method, Type: zapcore.StringType})
	ctx := r.Context()

	if r.Method == "GET" {
		userID, err := a.parseUserID(r)
		if err != nil {
			writeResponse(w, http.StatusBadRequest, fmt.Sprintf(`failed to parse user's	id: %s`, err))
			return
		}
		user, err := a.repository.GetUser(ctx, *userID)
		if err != nil {
			status := http.StatusInternalServerError
			switch {
			case errors.Is(err, ErrNotFound):
				status = http.StatusNotFound
			}
			writeResponse(w, status, fmt.Sprintf(`failed to get user with id %s: %s`, userID, err))
			return
		}
		writeJsonResponse(w, http.StatusOK, user)
	}

	if r.Method == "POST" {
		defer r.Body.Close()
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, err)
			return
		}

		user.ID, err = a.repository.AddUser(ctx, user)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, err)
			return
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, user)
	}
}

func (a *App) articleHandler(w http.ResponseWriter, r *http.Request) {
	a.logger.Info("articleHandler called", zap.Field{Key: "method", String: r.Method, Type: zapcore.StringType})
	ctx := r.Context()

	//todo: user -> article
	if r.Method == "GET" {
		userID, err := a.parseUserID(r)
		if err != nil {
			writeResponse(w, http.StatusBadRequest, fmt.Sprintf(`failed to parse user's	id: %s`, err))
			return
		}
		user, err := a.repository.GetUser(ctx, *userID)
		if err != nil {
			status := http.StatusInternalServerError
			switch {
			case errors.Is(err, ErrNotFound):
				status = http.StatusNotFound
			}
			writeResponse(w, status, fmt.Sprintf(`failed to get user with id %s: %s`, userID, err))
			return
		}
		writeJsonResponse(w, http.StatusOK, user)
	}

	if r.Method == "POST" {
		defer r.Body.Close()
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, err)
			return
		}

		user.ID, err = a.repository.AddUser(ctx, user)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, err)
			return
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, user)
	}
}

func (a *App) userArticlesHandler(w http.ResponseWriter, r *http.Request) {
	a.logger.Info("userArticlesHandler called", zap.Field{Key: "method", String: r.Method, Type: zapcore.StringType})
	userID, err := a.parseUserID(r)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, fmt.Sprintf(`failed to parse user'sid: %s`, err))
		return
	}
	articles, err := a.repository.GetUserArticles(r.Context(), *userID)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, fmt.Sprintf(`failed to getuser's (id: %s) articles: %s`, userID, err))
		return
	}
	writeJsonResponse(w, http.StatusOK, articles)
}

func (a *App) panicHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = recover()
		writeResponse(w, http.StatusOK, "panic logged, see server log")
	}()
	a.logger.Panic("panic!!!")
}

// func isNil(a interface{}) bool {
// 	defer func() { recover() }()
// 	return a == nil || reflect.ValueOf(a).IsNil()
// }

func (a *App) Init(ctx context.Context, logger *zap.Logger) error {
	config, err := pgxpool.ParseConfig(DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to parse conn string (%s): %w", DatabaseURL, err)
	}
	config.ConnConfig.LogLevel = pgx.LogLevelDebug
	config.ConnConfig.Logger = zapadapter.NewLogger(logger) // логгер запросов в БД

	// pool, err := pgxpool.ConnectConfig(ctx, config)
	// if err != nil {
	// 	return fmt.Errorf("unable to connect to database: %w", err)
	// }

	a.logger = logger

	// a.pool = pool
	// a.repository = NewCachedRepository(NewRepository(a.pool))
	//a.repository = NewCachedRepository(NewRepository(pool))

	//a.repository = NewCachedRepository(store.NewStore())
	a.repository, err = store.NewStore() // no cache
	if err != nil {
		return errors.New("Failed to create repository: " + err.Error())
	}

	return a.repository.InitSchema(ctx)
}

func (a *App) Serve() error {
	r := chi.NewRouter()
	//TODO: сделать полнотекстовым
	r.Get("/users", http.HandlerFunc(a.usersHandler))
	r.Get("/user/{id}", http.HandlerFunc(a.userHandler))
	//TODO: сделать полнотекстовым
	// r.Get("/user/{id}/articles", http.HandlerFunc(a.userArticlesHandler))
	r.Get("/user/articles", http.HandlerFunc(a.userArticlesHandler))
	r.Get("/panic", http.HandlerFunc(a.panicHandler))

	//TODO: сделать запрос на создание пользователя
	//TODO: сделать запрос на создание статьи
	r.Post("/user/add", http.HandlerFunc(a.userHandler))
	r.Post("/article/add", http.HandlerFunc(a.articleHandler))

	// profiling
	r.Mount("/debug", Profiler())
	return http.ListenAndServe("0.0.0.0:9000", r)
}

func Profiler() http.Handler {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, r.RequestURI+"/pprof/", http.StatusMovedPermanently)
	})
	r.HandleFunc("/pprof", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, r.RequestURI+"/", http.StatusMovedPermanently)
	})
	r.HandleFunc("/pprof/*", pprof.Index)
	r.HandleFunc("/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/pprof/profile", pprof.Profile)
	r.HandleFunc("/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/pprof/trace", pprof.Trace)
	r.Handle("/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/pprof/threadcreate", pprof.Handler("threadcreate"))
	r.Handle("/pprof/mutex", pprof.Handler("mutex"))
	r.Handle("/pprof/heap", pprof.Handler("heap"))
	r.Handle("/pprof/block", pprof.Handler("block"))
	r.Handle("/pprof/allocs", pprof.Handler("allocs"))
	return r
}
