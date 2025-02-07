package main

import "net/http"

func (app *aplicacion) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok\n"))
}
