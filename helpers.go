package main

import (
	"slices"
	"strings"
)

func stringReplace(s string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	e := strings.Split(s, " ")
	for i, word := range e {
		if slices.Contains(badWords, strings.ToLower(word)) {
			e[i] = "****"
		}
	}
	return strings.Join(e, " ")
}
