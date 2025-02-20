package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Tacoman44444/AimTrainer-Backend/internal/auth"
	"github.com/Tacoman44444/AimTrainer-Backend/internal/database"
	"github.com/Tacoman44444/AimTrainer-Backend/internal/response"
)

type apiConfig struct {
	db *database.Queries
}

func (cfg *apiConfig) recieveLoginInfoHandler(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		Username string `json:"username"`
		Password string `json:"hashed_password"`
	}
	params := Params{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Println("ERROR: unable to decode json data")
		fmt.Println(err)
		response.RespondWithJSON(w, 500, "unable to decode json data")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		fmt.Println("ERROR: could not hash password")
		fmt.Println(err)
		response.RespondWithJSON(w, 500, "unable to has password")
		return
	}

	_, err = cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Username:       params.Username,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		fmt.Println("ERROR: could not execute SQL query CreateUser")
		fmt.Println(err)
		response.RespondWithJSON(w, 500, "could not execute SQL query CreateUser")
		return
	}

	response.RespondWithJSON(w, 201, "user created successfully!")
}

func main() {
	servMux := http.NewServeMux()

	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		fmt.Println("ERROR: could not open database")
	}
	dbQueries := database.New(db)

	currentState := apiConfig{}
	currentState.db = dbQueries

	servMux.HandleFunc("POST /api/users", currentState.recieveLoginInfoHandler)

	server := http.Server{}
	server.Handler = servMux
	server.Addr = ":8080"

	http.ListenAndServe(server.Addr, server.Handler)
}
