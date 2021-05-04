package main

import (
    "context"
    "fmt"
    "net/http"
    "encoding/json"
    "log"
    //"io/ioutil"
    "time"
    "github.com/ShravanCool/SeeYouThere/models"
    "github.com/ShravanCool/SeeYouThere/helper"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    //"go.mongodb.org/mongo-driver/bson/primitive"
    //"go.mongodb.org/mongo-driver/mongo/options"
    //"go.mongodb.org/mongo-driver/mongo/readpref"
)

//func apiStatus(w http.ResponseWriter, r *http.Request) {
    //w.Header().Set("Content-Type", "application/json")
    //w.WriteHeader(http.StatusOK)
    //w.Write([]byte(`{"message": "Server is up and running.."}`))
//}

func MeetingSchedule(w http.ResponseWriter, r *http.Request) {
    //fmt.Println("Post method to add or schedule a meeting")
    w.Header().Set("Content-Type", "application/json")
    var meeting models.Meeting

    err := json.NewDecoder(r.Body).Decode(&meeting)

    meeting.ID = helper.GetID(32)
    meeting.CreatedAt = time.Now().Unix()
    fmt.Println("Meeting title- ", meeting.Title)

    if err != nil {
        fmt.Println(err.Error())
        http.Error(w, "Invalid Request Body", http.StatusBadRequest)
        return
    }

    if meeting.StartTime > meeting.EndTime {
        http.Error(w, "Start Time should be before End Time", http.StatusBadRequest)
    }

    emails := make(bson.A, len(meeting.Participants))
    for i := 0; i < len(meeting.Participants); i++ {
        emails[i] = meeting.Participants[i].Email

        _, found := helper.Find([]string{"Yes", "No", "Maybe", "Not Answered"}, meeting.Participants[i].RSVP)
        if !found {
            http.Error(w, "Invalid Request Body", http.StatusBadRequest)
            return
        }
    }

    unwind_stage := bson.D{{
        "$unwind", "$participants",
    }}

    match_stage := bson.D{{
        "$match", bson.M{
            "participants.email": bson.M{
                "$in": emails,
            },
            "participants.rsvp": bson.M{
                "$in": bson.A{"Yes", "Maybe"},
            },
        },
    }}

    check_time_stage := bson.D{{
        "$match", bson.M{
            "$or": bson.A{
                bson.M{
                    "startTime": bson.M{"$lte": meeting.EndTime},
                    "endTime": bson.M{"$gte": meeting.EndTime},
                },
                bson.M{
                    "startTime": bson.M{"$lte": meeting.StartTime},
                    "endTime": bson.M{"$gte": meeting.EndTime},
                },
                bson.M{
                    "startTime": bson.M{"$gte": meeting.StartTime},
                    "endTime": bson.M{"$lte": meeting.EndTime},
                },
            },
        },
    }}

    var dbase = models.ConnectDB()
    var meetingCollection = dbase.Collection("meetings")
    meetingCursor, err := meetingCollection.Aggregate(context.TODO(), mongo.Pipeline{unwind_stage, match_stage, check_time_stage})

    var meetings []bson.M
    if err = meetingCursor.All(context.TODO(), &meetings); err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }

    if len(meetings) != 0 {
        log.Println("Conflict occured with the meetings in the response.")
        json.NewEncoder(w).Encode(meeting)
        return
    }

    _, err = meetingCollection.InsertOne(context.TODO(), meeting)

    if err != nil {
        panic(err)
        return
    }

    json.NewEncoder(w).Encode(meeting)
}

func GetMeetingByID(w http.ResponseWriter, r *http.Request) {
    id := helper.GetParams(r)

    var meeting bson.M

    var dbase = models.ConnectDB()
    var meetingCollection = dbase.Collection("meetings")

    err := meetingCollection.FindOne(context.TODO(), bson.D{{
        "id", id,
    }}).Decode(&meeting)

    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        panic(err)
    }

    json.NewEncoder(w).Encode(meeting)

}

func main() {
    fmt.Println("Hello World!!")
    //http.HandleFunc("/", apiStatus)
    http.HandleFunc("/meetings", MeetingSchedule)
    http.HandleFunc("/meeting/", GetMeetingByID)
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        panic(err)
    }

}
