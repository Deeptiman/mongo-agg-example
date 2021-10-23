package main

import (
	"context"
	"encoding/json"
	"fmt"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	ctx    = context.Background()
	client *mongo.Client
	err    error
	log    hclog.Logger
)

func main() {
	fmt.Println("Mongo Agg Example!")

	PORT := ":5000"
	log = hclog.Default()

	router := mux.NewRouter()
	router.StrictSlash(true)

	// GET
	request := router.Methods(http.MethodGet).Subrouter()
	request.HandleFunc("/api/movies/filter-genere", FilterGenre).
		Queries("genre", "{genre:[A-Za-z,]+}")
	request.HandleFunc("/api/movies/bucket-query", BucketQuery)
	request.HandleFunc("/api/movies/geonear", GeoNearQuery).
		Queries("lat", "{lat:[-0-9.]+}").
		Queries("long", "{long:[0-9.,]+}")	
	request.HandleFunc("/api/movies/graph-lookup", GraphLookup).
		Queries("email", "{email:[a-zA-Z0-9+_.-]+@[a-zA-Z0-9.-]+$}")

	mongoConnection()

	cors := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	// create the http server
	server := http.Server{
		Addr:         PORT,
		Handler:      cors(router),
		ErrorLog:     log.StandardLogger(&hclog.StandardLoggerOptions{}),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Info("Starting server on", "PORT", PORT)

	err = server.ListenAndServe()
	if err != nil {
		log.Error("Unable to start server", "error", err)
		os.Exit(1)
	}
}

func mongoConnection() {

	wMajority := writeconcern.New(writeconcern.WMajority())

	poolMonitor, cmdMonitor, serverMonitor := MongoMonitors()

	// mongodb://pd123:pd123@localhost:27017/product?authSource=product&replicaSet=rs0

	// mongodb+srv://admin:admin123@cluster0.p5zvm.mongodb.net/product?authSource=admin&replicaSet=atlas-v7g7gu-shard-0&readPreference=primary&appname=MongoDB%20Compass&ssl=true

	// mongosh mongodb+srv://cluster0.p5zvm.mongodb.net/sample_mflix --username admin --password admin123 --authenticationDatabase admin

	client, err = mongo.NewClient(options.Client().
		ApplyURI("mongodb+srv://admin:admin123@cluster0.p5zvm.mongodb.net/sample_mflix?authSource=admin&replicaSet=atlas-v7g7gu-shard-0&readPreference=primary&appname=MongoDB%20Compass&ssl=true").
		SetPoolMonitor(poolMonitor).
		SetMonitor(cmdMonitor).
		SetServerMonitor(serverMonitor).
		SetWriteConcern(wMajority))
	if err != nil {
		fmt.Println("Error initializing to MongoDB : " + err.Error())
		return
	}
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		fmt.Println("Error connecting to MongoDB : " + err.Error())
		return
	}
	log.Info("MongoDB Connected!")

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Error("Mongo Ping Error : ", err.Error())
		return
	}
	log.Info("MongoDB Connection Pinged!")

	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Error("Mongo ListDatabase Error : ", err.Error())
		return
	}
	log.Info("database", "list", databases)
}

func FilterGenre(w http.ResponseWriter, r *http.Request) {

	ctx, _ := context.WithTimeout(context.Background(), 50*time.Second)

	genre := strings.TrimSpace(r.URL.Query().Get("genre"))

	info := filterGenreAgg(ctx, client, genre)

	RespondJSON(w, http.StatusOK, info)
}

func BucketQuery(w http.ResponseWriter, r *http.Request) {

	ctx, _ := context.WithTimeout(context.Background(), 50*time.Second)

	graph := bucketAggQuery(ctx, client)

	RespondJSON(w, http.StatusOK, graph)
}

func GeoNearQuery(w http.ResponseWriter, r *http.Request) {

	ctx, _ := context.WithTimeout(context.Background(), 50*time.Second)

	lat := strings.TrimSpace(r.URL.Query().Get("lat"))
	long := strings.TrimSpace(r.URL.Query().Get("long"))

	locations := geoNearAggQuery(ctx, client, lat, long)

	RespondJSON(w, http.StatusOK, locations)

}

func GraphLookup(w http.ResponseWriter, r *http.Request) {

	ctx, _ := context.WithTimeout(context.Background(), 50*time.Second)

	email := strings.TrimSpace(r.URL.Query().Get("email"))

	graph := graphLookup(ctx, client, email)

	RespondJSON(w, http.StatusOK, graph)
}

func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "POST, GET, OPTIONS, PUT, DELETE")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
