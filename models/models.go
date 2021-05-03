package models

//import (
    //"go.mongodb.org/mongo-driver/bson/primitive"
//)

type Participant struct {
    Name string `json:"name" bson:"name"`
    Email string `json:"email" bson:"email"`
    RSVP string `json:"rsvp" bson:"rsvp"`
}

type Meeting struct {
    ID string `json:"id" bson:"id"`
    Title string `json:"title" bson:"title"`
    Participants []Participant `json:"participants" bson:"participants"`
    StartTime int64 `json:"startTime" bson:"startTime"`
    EndTime int64 `json:"endTime" bson:"endTime"`
    CreatedAt int64 `json:"creationTime" bson:"creationTime"`
}


