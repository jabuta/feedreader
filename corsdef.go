package main

import "github.com/go-chi/cors"

var corsOptions = cors.Options{
	AllowedOrigins:   []string{"*"},
	AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT"},
	AllowCredentials: true,
}
