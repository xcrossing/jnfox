package util

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/xcrossing/jnfo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MgInstance struct {
	client     *mongo.Client
	collection *mongo.Collection
}

type MgDoc struct {
	Bango string
	Stars []string
}

const ext = ".jpg"

func NewMgInstance(cfgMongo ConfigMongo) (*MgInstance, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(cfgMongo.Uri))
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
	db := client.Database(cfgMongo.Db)
	collection := db.Collection(cfgMongo.Collection)

	return &MgInstance{client: client, collection: collection}, nil
}

func (mg *MgInstance) Close() {
	mg.client.Disconnect(aCtx())
}

func (mg *MgInstance) Fetch(bango string) (*MgDoc, error) {
	doc := &MgDoc{}
	result := mg.collection.FindOne(aCtx(), bson.D{{"bango", bango}})
	if err := result.Err(); err != nil {
		return nil, err
	}
	if err := result.Decode(doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func (mg *MgInstance) BatchFetch(bangos []string) (*[]MgDoc, error) {
	docs := make([]MgDoc, len(bangos))
	cursor, err := mg.collection.Find(aCtx(), bson.D{{"bango", bson.D{{"$in", bangos}}}})
	if err != nil {
		return nil, err
	}
	err = cursor.All(aCtx(), &docs)
	if err != nil {
		return nil, err
	}
	return &docs, nil
}

func (mg *MgInstance) InsertOne(nfo *jnfo.Jnfo) error {
	doc := bson.D{
		{"title", nfo.Title},
	}
	if duration, err := strconv.Atoi(nfo.Duration); err == nil {
		doc = append(doc, primitive.E{"duration", duration})
	}
	if nfo.Studio != "" {
		doc = append(doc, primitive.E{"studio", nfo.Studio})
	}
	if nfo.Label != "" {
		doc = append(doc, primitive.E{"label", nfo.Label})
	}
	if nfo.Serie != "" {
		doc = append(doc, primitive.E{"series", nfo.Serie})
	}
	if nfo.Director != "" {
		doc = append(doc, primitive.E{"director", nfo.Director})
	}
	doc = append(doc, primitive.E{"categories", nfo.Categories})
	doc = append(doc, primitive.E{"stars", nfo.Cast})
	doc = append(doc, primitive.E{"prefix", nfo.Prefix()})
	doc = append(doc, primitive.E{"created_at", time.Now().Format("2006-01-02T15:04:05Z")})
	doc = append(doc, primitive.E{"count_stars", len(nfo.Cast)})
	doc = append(doc, primitive.E{"published_at", nfo.Date})
	doc = append(doc, primitive.E{"bango", nfo.Num})

	_, err := mg.collection.InsertOne(aCtx(), doc)
	return err
}

func aCtx() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return ctx
}

func (d *MgDoc) PicName() string {
	if len(d.Stars) > 0 {
		return fmt.Sprintf("%s-%s%s", d.Bango, strings.Join(d.Stars, " "), ext)
	}
	return fmt.Sprintf("%s%s", d.Bango, ext)
}
