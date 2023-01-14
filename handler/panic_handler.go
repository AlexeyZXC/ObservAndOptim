package handler

import "net/http"

type PanicHandler struct{}

func (i PanicHandler) Handle() {
	panic("Panic")
}
func (i PanicHandler) Log() {
	//Отправить лог в sentry
}

//chi mux
func (i PanicHandler) Handle_chi(w http.ResponseWriter, r *http.Request) {
	panic("Panic")
}

func (i PanicHandler) Log_chi(w http.ResponseWriter, r *http.Request) {
	//Отправить лог в sentry
}
