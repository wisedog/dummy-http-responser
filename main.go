package main

import (
	"net/http"
	"os"

	mgo "github.com/globalsign/mgo"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

var db *mgo.Database

const apiVersion = "v1"

func init() {
	// connect database
	mongoDBURI := os.Getenv("MONGODB_URI")
	if mongoDBURI == "" {
		panic("mongoDB URI is empty")
	}
	session, err := mgo.Dial(mongoDBURI)

	if err != nil {
		log.Fatal(err)
	}
	dbName := os.Getenv("MONGODB_DATABASE")
	db = session.DB(dbName)
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

func main() {
	// TODO load env

	router := createRoute()
	// to support for CORS
	handler := cors.Default().Handler(router)

	// Start the server with port. default port is 3000 for Heroku
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("Starting Dummy Http Responser on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func createRoute() *httprouter.Router {
	// setup router
	router := httprouter.New()
	router.GET("/echo", handleEcho)
	router.POST("/echo", handleEcho)
	router.PUT("/echo", handleEcho)
	router.DELETE("/echo", handleEcho)

	router.POST("/create", handleV1CreateDummy)

	// fast http router is not support chaining or multiple methods setting at once
	router.GET("/v1/:id", handleV1Custom)
	router.POST("/v1/:id", handleV1Custom)
	router.PUT("/v1/:id", handleV1Custom)
	router.DELETE("/v1/:id", handleV1Custom)

	return router
}
