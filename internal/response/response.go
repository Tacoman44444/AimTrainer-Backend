package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	fmt.Println(msg)
}

func RespondWithJSON(w http.ResponseWriter, code int, body interface{}) {
	data, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(500)
		fmt.Println("ERROR: could not convert data to json, responded with code 500 --- RespondWithJSON()")
	}
	w.WriteHeader(code)
	w.Write(data)
}
