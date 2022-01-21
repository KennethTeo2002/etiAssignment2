package main

import(
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func InsertSchedule(collection *mongo.Collection, ctx context.Context, classCode string, scheduleDateTime string){
	doc := bson.M{
		{Key: "ClassCode" , Value: classCode},
		{Key: "Schedule" , Value: scheduleDateTime},
	}
	result,err := collection.InsertOne(ctx.TODO(), doc)
}

func GetSchedule(collection *mongo.Collection, ctx context.Context) []scheduleInfo{
	
}