package httpResponse

import (
	"encoding/json"
	"log"
	"net/http"
)

const errHeader = "httpResponse"

// ReturnInternalError 500
func ReturnInternalError(w http.ResponseWriter, r *http.Request, error error, errFrom string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	js, err := json.MarshalIndent(error, "", "\t")
	if err != nil {
		log.Fatalln("handlers returnInternalError MarshalIndent error: \n", err)
	}
	w.Write(js)
	log.Fatal(errFrom, error)
}

// ReturnSuccessStatus 200
func ReturnSuccessStatus(w http.ResponseWriter, r *http.Request, v interface{}) {
	js, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		ReturnInternalError(w, r, err, errHeader + ".ReturnSuccessStatus MarshalIndent error: \n")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func ReturnBadRequest(w http.ResponseWriter, r *http.Request) {

}
