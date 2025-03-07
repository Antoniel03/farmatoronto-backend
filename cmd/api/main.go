package main

import (
	"log"

	"github.com/Antoniel03/farmatoronto-backend/internal/db"
	"github.com/Antoniel03/farmatoronto-backend/internal/env"
	"github.com/Antoniel03/farmatoronto-backend/internal/store"
)

func main() {

	cfg := config{
		addr: env.GetString("ADDR", ":8081"),
		db: dbConfig{
			addr: env.GetString("DB_ADDR", "../../internal/db/farma_db.db"),
		},
		jwtAuth: jwtConfig{
			tokenAuth:  NewJWTAuth("a-very-ultra-super-secure-secret!", "HS256"),
			expiration: env.GetInt64("JWT_EXP", 3600),
		},
	}

	database, err := db.OpenDB(cfg.db.addr)
	if err != nil {
		log.Panic(err)
	}
	defer db.CloseDB(database)

	err = db.SetupDB(database)
	if err != nil {
		log.Panic(err)
	}

	store := store.NewSQLiteStorage(database)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	log.Println("Servidor iniciado en el puerto" + app.config.addr)
	log.Fatal(app.run(mux))

}
