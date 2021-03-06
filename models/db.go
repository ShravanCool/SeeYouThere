package models

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"

    "github.com/ShravanCool/SeeYouThere/helper"
)


func ConnectDB() *mongo.Database {

    uri := helper.GetEnvi("uri")

    //Set client options 
    //clientOptions := options.Client().ApplyURI(uri)
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    //Connect to MongoDB
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

    if err != nil {
        panic(err)
    }

    fmt.Println("Connected to MongoDB!")

    return client.Database("MeetingsAPI")
    //collection := client.Database("DB_name").Collection("collection_name")

    //return collection
}

//Error response model
type ErrorResponse struct {
    StatusCode int `json:"status"`
    ErrorMessage string `json:message`
}

func GetError(err error, w http.ResponseWriter) {
    log.Fatal(err.Error())
    var response = ErrorResponse{
        ErrorMessage: err.Error(),
        StatusCode: http.StatusInternalServerError,
    }

    message, _ := json.Marshal(response)
    w.WriteHeader(response.StatusCode)
    w.Write(message)
}
