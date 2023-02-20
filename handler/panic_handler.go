package handler

import (
	"net/http"
)

type PanicHandler struct {
	Logger iLogger
}

func (i PanicHandler) Handle() {
	panic("Panic")
}
func (i PanicHandler) Log() {
	//Отправить лог в sentry
}

//chi mux
func (i PanicHandler) Handle_chi(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query().Get("msg")
	i.Logger.Error("Panic !!! from request: %v", msg)
	panic("Panic")
}

func (i PanicHandler) Log_chi(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query().Get("msg")
	i.Logger.Info("Log from request: %v", msg)
	//Отправить лог в sentry
}
