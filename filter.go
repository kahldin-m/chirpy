package main

import (
	"strings"
)


func cleanChirp (s string) string {
	badWords := map[string]struct{}{
	"kerfuffle": {},
	"sharbert": {},
	"fornax": {},
	}

	soapedChirp := strings.Split(s, " ")

	for i, word := range soapedChirp {
		_, ok := badWords[strings.ToLower(word)]
		if ok {
			soapedChirp[i] = "****"
		}
	}

	return strings.Join(soapedChirp, " ")
}