package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo
	return Models{
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID       string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string    `bson:"name" json:"name"`
	Data     string    `bson:"data" json:"data"`
	CreateAt time.Time `bson:"created_at" json:"created_at"`
	UpdateAt time.Time `bson:"update_at" json:"update_at"`
}

// Insert insert one data to mongo Db (structure => LogEntry)
func (l *LogEntry) Insert(entry LogEntry) error {
	// connect to logger database use logs collection
	collection := client.Database("logger").Collection("logs")

	// insert one in database logs
	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:     entry.Name,
		Data:     entry.Data,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	})

	if err != nil {
		return err
	}
	return nil
}

func (l *LogEntry) All() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// connect to logger database use logs collection
	collection := client.Database("logger").Collection("logs")

	// opts like statment in sql
	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	// cursor are result in pointer
	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println()
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*LogEntry

	// fetch each all cursor writer
	for cursor.Next(ctx) {
		var item LogEntry
		// Decode each data
		err := cursor.Decode(&item)
		if err != nil {
			log.Println("error decoding log into slice item", err)
			return nil, err
		}
		// append items
		logs = append(logs, &item)
	}
	return logs, nil
}

func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// connect to logger database use logs collection
	collection := client.Database("logger").Collection("logs")
	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var entry LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docId}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (l *LogEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// connect to logger database use logs collection
	collection := client.Database("logger").Collection("logs")

	err := collection.Drop(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// connect to logger database use logs collection
	collection := client.Database("logger").Collection("logs")
	docId, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		return nil, err
	}
	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docId},
		bson.D{
			{"$set", bson.D{
				{"name", l.Name},
				{"data", l.Data},
				{"updated_at", time.Now()},
			}},
		},
	)

	if err != nil {
		return nil, err
	}
	return result, nil
}
