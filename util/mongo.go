package util

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mgInstance struct {
	client     *mongo.Client
	collection *mongo.Collection
}

type mgDoc struct {
	Bango string
	Stars []string
}

const ext = ".jpg"

func NewMgInstance() (*mgInstance, error) {
	mongoCfg := viper.Sub("mongo")
	if mongoCfg == nil {
		return nil, errors.New("no mongo config !")
	}

	uri := mongoCfg.GetString("uri")
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// connect
	err = client.Connect(aCtx())
	if err != nil {
		return nil, err
	}

	// ping
	if err := client.Ping(aCtx(), readpref.Primary()); err != nil {
		return nil, err
	}

	// get collection
	db := client.Database(mongoCfg.GetString("db"))
	collection := db.Collection(mongoCfg.GetString("collection"))

	return &mgInstance{client: client, collection: collection}, nil
}

func (mg *mgInstance) Close() {
	mg.client.Disconnect(aCtx())
}

func (mg *mgInstance) Fetch(bango string) (*mgDoc, error) {
	doc := &mgDoc{}
	err := mg.collection.FindOne(aCtx(), bson.D{{"bango", bango}}).Decode(doc)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func aCtx() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return ctx
}

func (d *mgDoc) PicName() string {
	if len(d.Stars) > 0 {
		return fmt.Sprintf("%s-%s%s", d.Bango, strings.Join(d.Stars, " "), ext)
	}
	return fmt.Sprintf("%s%s", d.Bango, ext)
}
