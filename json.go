package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Gfarf/Chirpy/internal/auth"
	"github.com/Gfarf/Chirpy/internal/database"
	"github.com/google/uuid"
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

/*
	func handler(w http.ResponseWriter, r *http.Request){
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
	}
*/
func mappingChirp(d *database.Chirp) StoringChirp {
	res := StoringChirp{}
	if d == nil {
		return res
	}
	res.Body = d.Body
	res.CreatedAt = d.CreatedAt
	res.UpdatedAt = d.UpdatedAt
	res.ID = d.ID
	res.UserID = d.UserID
	return res
}
func (cfg *apiConfig) handlerChirp(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
		//UserID uuid.UUID `json:"user_id"`
	}
	type response struct {
		//Body  string `json:"cleaned_body"`
		Error string `json:"error"`
		//Valid bool   `json:"valid"`
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
		//res.Valid = false
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
	//res.Body = stringReplace(chirp1.Body)
	//res.Valid = true
	jwt, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting Bearer Token: %s", err)
		w.WriteHeader(500)
		return
	}
	UserID, err := auth.ValidateJWT(jwt, cfg.secretString)
	if err != nil {
		log.Printf("Invalid token: %s", err)
		w.WriteHeader(401)
		return
	}
	fRes2 := StoringChirp{}
	fRes, err := cfg.dbQueries.CreateStoredChirp(r.Context(), database.CreateStoredChirpParams{Body: chirp1.Body, UserID: UserID})
	fRes2 = mappingChirp(&fRes)
	if err != nil {
		log.Printf("Error loading chirp to database: %s", err)
		w.WriteHeader(500)
		return
	}
	dat, err := json.Marshal(fRes2)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(dat)
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	fAll, err := cfg.dbQueries.GetChirpsAll(r.Context())
	if err != nil {
		log.Printf("Error getting chirps from database: %s", err)
		w.WriteHeader(500)
		return
	}
	allList := []StoringChirp{}
	for _, fRes := range fAll {
		fRes2 := mappingChirp(&fRes)
		allList = append(allList, fRes2)
	}

	dat, err := json.Marshal(allList)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}

func (cfg *apiConfig) handlerGetOneChirp(w http.ResponseWriter, r *http.Request) {
	chripID := r.PathValue("chirpID")
	uuidChirpID, err := uuid.Parse(chripID)
	if err != nil {
		log.Printf("Error parsing chirp uuid: %s", err)
		w.WriteHeader(404)
		return
	}
	c, err := cfg.dbQueries.GetOneChirpByID(r.Context(), uuid.UUID(uuidChirpID))
	if err != nil {
		log.Printf("Error getting chirp from database: %s", err)
		w.WriteHeader(404)
		return
	}
	chirp := mappingChirp(&c)
	dat, err := json.Marshal(chirp)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}
