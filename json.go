package main

import (
	"encoding/json"
	"log"
	"net/http"
)

/*func handler(w http.ResponseWriter, r *http.Request){
    type parameters struct {
        // these tags indicate how the keys in the JSON should be mapped to the struct fields
        // the struct fields must be exported (start with a capital letter) if you want them parsed
        Name string `json:"name"`
        Age int `json:"age"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        // an error will be thrown if the JSON is invalid or has the wrong types
        // any missing fields will simply have their values in the struct set to their zero value
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
    }
    // params is a struct with data populated successfully
    // ...
}*/

/*func handler(w http.ResponseWriter, r *http.Request){
    // ...

    type returnVals struct {
        // the key will be the name of struct field unless you give it an explicit JSON tag
        CreatedAt time.Time `json:"created_at"`
        ID int `json:"id"`
    }
    respBody := returnVals{
        CreatedAt: time.Now(),
        ID: 123,
    }
    dat, err := json.Marshal(respBody)
	if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
	}
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)
    w.Write(dat)
}*/

func handlerChirp(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}
	type response struct {
		Body  string `json:"cleaned_body"`
		Error string `json:"error"`
		Valid bool   `json:"valid"`
	}
	decoder := json.NewDecoder(r.Body)
	chirp1 := chirp{}
	err := decoder.Decode(&chirp1)
	if err != nil {
		log.Printf("Error decoding chirp: %s", err)
		w.WriteHeader(500)
		return
	}
	res := response{}
	if len(chirp1.Body) > 140 {
		res.Error = "Chirp is too long"
		res.Valid = false
		dat, err := json.Marshal(res)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}
	res.Body = stringReplace(chirp1.Body)
	res.Valid = true
	dat, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}
