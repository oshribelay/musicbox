package server

import (
	"net/http"

	"github.com/oshribelay/musicbox/internal/spotify"
)

func Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /spotify/login", spotify.LoginHandler)
	mux.HandleFunc("GET /spotify/callback", spotify.CallbackHandler)
	
	return http.ListenAndServe("localhost:8080", mux) 
 }