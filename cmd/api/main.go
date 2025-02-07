package main

import "log"
import "github.com/Antoniel03/farmatoronto-backend/internal/env"

func main() {
	cfg := config{addr: env.GetString("ADDR", ":8081")}
	app := &aplicacion{config: cfg}
	mux := app.mount()
	log.Println("Servidor iniciado en el puerto" + app.config.addr)
	log.Fatal(app.run(mux))
}
