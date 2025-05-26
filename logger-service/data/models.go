// Package data contains the data models and database operations
package data

import (
	// Import context for handling timeouts and cancellations
	"context"
	// Import log for logging errors
	"log"
	// Import time for timestamps
	"time"

	// Import MongoDB driver packages
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Global MongoDB client
var client *mongo.Client

// New creates a new Models instance with MongoDB client
func New(mongo *mongo.Client) Models {
	// Store MongoDB client in global variable
	client = mongo

	// Return a new Models instance with an empty LogEntry
	return Models{
		LogEntry: LogEntry{},
	}
}

// Models holds the different model types
type Models struct {
	// LogEntry model for log data
	LogEntry LogEntry
}

// LogEntry represents a log entry in the database
type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// Insert adds a new log entry to the database
func (l *LogEntry) Insert(entry LogEntry) error {
	// Get the logs collection from the logs database
	collection := client.Database("logs").Collection("logs")

	// Insert the log entry with current timestamps
	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	// Check for errors during insertion
	if err != nil {
		log.Println("Error inserting into logs:", err)
		return err
	}
	// Return nil if successful
	return nil
}

// AllLogs retrieves all log entries from the database
func (l *LogEntry) AllLogs() ([]*LogEntry, error) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	// Ensure context is canceled when function completes
	defer cancel()

	// Get the logs collection
	collection := client.Database("logs").Collection("logs")

	// Set options for sorting results
	opts := options.Find()
	// Sort by created_at in descending order (newest first)
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	// Find all documents with empty filter
	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	// Check for errors during query
	if err != nil {
		log.Println("Finding all docs error:", err)
		return nil, err
	}
	// Ensure cursor is closed when function completes
	defer cursor.Close(ctx)

	// Slice to hold the log entries
	var logs []*LogEntry
	// Iterate through the cursor
	for cursor.Next(ctx) {
		// Create a variable for each log entry
		var item LogEntry
		// Decode the current document into the item
		err := cursor.Decode(&item)
		// Check for decoding errors
		if err != nil {
			log.Println("Error decoding log:", err)
			return nil, err
		} else {
			// Add the log entry to the slice
			logs = append(logs, &item)
		}
	}
	// Return the log entries and nil error
	return logs, nil
}

// GetOne retrieves a single log entry by ID
func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	// Ensure context is canceled when function completes
	defer cancel()

	// Get the logs collection
	collection := client.Database("logs").Collection("logs")

	// Convert the ID string to MongoDB ObjectID
	docID, err := primitive.ObjectIDFromHex(id)
	// Check for conversion errors
	if err != nil {
		return nil, err
	}

	// Variable to hold the found entry
	var entry LogEntry
	// Find one document with the matching ID and decode it
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	// Check for errors during find or decode
	if err != nil {
		return nil, err
	}
	// Return the found entry and nil error
	return &entry, nil
}

func (l *LogEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	if err := collection.Drop(ctx); err != nil {
		return err
	}
	return nil
}

func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(l.ID)
	// Check for conversion errors
	if err != nil {
		return nil, err
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
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
