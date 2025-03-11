package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Tacoman44444/AimTrainer-Backend/internal/auth"
	"github.com/Tacoman44444/AimTrainer-Backend/internal/database"
	"github.com/Tacoman44444/AimTrainer-Backend/internal/response"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	db *database.Queries
}

func (cfg *apiConfig) recieveLoginInfoHandler(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		Username string `json:"username"`
		Password string `json:"password"`
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
		response.RespondWithJSON(w, 500, "unable to hash password")
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

func (cfg *apiConfig) recieveSessionInfoHandler(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Score    int    `json:"score"`
		Accuracy string `json:"accuracy"`
	}

	params := Params{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		fmt.Println("ERROR: unable to decode JSON data")
		fmt.Println(err)
		response.RespondWithJSON(w, 500, "unable to decode JSON data")
		return
	}

	userData, err := cfg.db.FindUserByUsername(r.Context(), params.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("ERROR: no user by this name")
			fmt.Println(err)
			response.RespondWithJSON(w, 400, "no user goes by "+params.Username)
			return
		}
		fmt.Println("ERROR: could not execute SQL query!")
		fmt.Println(err)
		response.RespondWithJSON(w, 500, "could not execute SQL query")
		return
	}

	err = auth.CheckPasswordHash(params.Password, userData.HashedPassword)
	if err != nil {
		fmt.Println("ERROR: incorrect password")
		fmt.Println(err)
		response.RespondWithJSON(w, 500, "incorrect password -- "+params.Password) //params.password only included for testing
		return
	}

	_, err = cfg.db.CreateSession(r.Context(), database.CreateSessionParams{
		Score:    int32(params.Score),
		Accuracy: params.Accuracy,
		PlayerID: userData.ID,
	})

	if err != nil {
		fmt.Println("ERROR: could not execute SQL query")
		fmt.Println(err)
		response.RespondWithJSON(w, 500, "could not execute SQL query")
	}

	response.RespondWithJSON(w, 201, "session created successfully")

}

func main() {
	godotenv.Load()
	servMux := http.NewServeMux()
	dbUrl := os.Getenv("DB_URL")
	fmt.Println(dbUrl)
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		fmt.Println("ERROR: could not open database")
		fmt.Println(err)
		return
	}
	dbQueries := database.New(db)

	currentState := apiConfig{}
	currentState.db = dbQueries

	servMux.HandleFunc("POST /api/users", currentState.recieveLoginInfoHandler)
	servMux.HandleFunc("POST /api/sessions", currentState.recieveSessionInfoHandler)

	server := http.Server{}
	server.Handler = servMux
	server.Addr = ":8080"

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
