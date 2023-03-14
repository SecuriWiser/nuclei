package mongo

import (
	"context"
	"github.com/projectdiscovery/gologger"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Connect() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	uri := os.Getenv("MONGO_URI")
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	var err error
	Client, err = mongo.Connect(context.TODO(), opts)

	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	var result bson.M
	if err := Client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	gologger.Info().Msg("Pinged...! Successfully connected to MongoDB!\n")
}

func Disconnect() {
	if Client != nil {
		if err := Client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
		gologger.Info().Msg("Disconnected to mongo!\n")
	}
}
