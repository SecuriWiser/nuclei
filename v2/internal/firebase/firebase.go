package firebase

import (
	"cloud.google.com/go/firestore"
	"context"
	_ "embed"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"log"
)

var Client *firestore.Client

//go:embed serviceAccount.json
var serviceAccount []byte

func Connect() {
	ctx := context.Background()
	//sa := option.WithCredentialsFile("path/to/serviceAccount.json")
	sa := option.WithCredentialsJSON(serviceAccount)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	Client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
}

func Disconnect() {
	err := Client.Close()
	if err != nil {
		return
	}
}
