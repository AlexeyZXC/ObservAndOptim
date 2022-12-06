package handler

type PanicHandler struct{}

func (i PanicHandler) Handle() {
	panic("Panic")
}
func (i PanicHandler) Log() {
	//Отправить лог в sentry
}
