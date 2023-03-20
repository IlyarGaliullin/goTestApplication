package mongodb

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"testApplication/models"
	"testApplication/utils"
)

type mongodb struct {
	client *mongo.Client
}

func InitConnection() *mongodb {

	host := utils.Conf.Get("mongodb.host")
	port := utils.Conf.GetInt("mongodb.port")

	uri := fmt.Sprintf("mongodb://%s:%d", host, port)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	return &mongodb{client: client}
}

func (m mongodb) GetClients(ctx context.Context, offset int, limit int) []models.Client {

	coll := m.client.Database("admin").Collection("clients")

	filter := bson.D{}
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil
	}

	var clients []models.Client
	err = cursor.All(ctx, &clients)
	if err != nil {
		return nil
	}

	return clients
}

func (m mongodb) GetClientById(ctx context.Context, id int) (models.Client, error) {

	coll := m.client.Database("admin").Collection("clients")

	filter := bson.D{{"id", bson.D{{"&eq", id}}}}

	var client models.Client
	err := coll.FindOne(ctx, filter).Decode(&client)
	if err != nil {
		return models.Client{}, nil
	}

	return client, nil
}

func (m mongodb) CreateClient(ctx context.Context, client models.Client) (models.Client, error) {

	coll := m.client.Database("admin").Collection("clients")

	_, err := coll.InsertOne(ctx, client)
	if err != nil {
		return models.Client{}, nil
	}

	return client, nil
}

func (m mongodb) UpdateClient(ctx context.Context, client models.Client) error {

	coll := m.client.Database("admin").Collection("clients")

	filter := bson.D{{"id", bson.D{{"&eq", client.Id}}}}
	update := bson.D{{"name", client.Name}}

	updateResult, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if updateResult.UpsertedCount == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

func (m mongodb) DeleteClient(ctx context.Context, id int) error {

	coll := m.client.Database("admin").Collection("clients")

	filter := bson.D{{"id", bson.D{{"&eq", id}}}}

	deleteResult, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if deleteResult.DeletedCount == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
