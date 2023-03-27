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
	client            *mongo.Client
	database          *mongo.Database
	clientsCollection *mongo.Collection
}

func InitConnection() *mongodb {

	host := utils.Conf.Get("mongodb.host")
	port := utils.Conf.GetInt("mongodb.port")
	database := utils.Conf.GetString("mongodb.database")

	uri := fmt.Sprintf("mongodb://%s:%d", host, port)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	mongoDatabase := client.Database(database)
	mongoClientsCollection := client.Database(database).Collection("clients")

	if err != nil {
		log.Fatal(err)
	}

	return &mongodb{client: client, database: mongoDatabase, clientsCollection: mongoClientsCollection}
}

func (m mongodb) GetClients(ctx context.Context, offset int, limit int) ([]models.Client, error) {

	var clients []models.Client

	filter := bson.D{}
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))

	cursor, err := m.clientsCollection.Find(ctx, filter, opts)
	if err != nil {
		log.Println(err)
		return clients, err
	}

	err = cursor.All(ctx, &clients)
	if err != nil {
		log.Println(err)
		return clients, err
	}

	return clients, nil
}

func (m mongodb) GetClientById(ctx context.Context, id int) (models.Client, error) {

	filter := bson.D{{"id", id}}

	var client models.Client
	err := m.clientsCollection.FindOne(ctx, filter).Decode(&client)
	if err != nil {
		log.Println(err)
		return models.Client{}, err
	}

	return client, nil
}

func (m mongodb) CreateClient(ctx context.Context, client models.Client) (models.Client, error) {

	_, err := m.clientsCollection.InsertOne(ctx, client)
	if err != nil {
		log.Println(err)
		return models.Client{}, err
	}

	return client, nil
}

func (m mongodb) UpdateClient(ctx context.Context, client models.Client) error {

	filter := bson.D{{"id", client.Id}}
	update := bson.D{{"$set", bson.D{{"name", client.Name}}}}

	updateResult, err := m.clientsCollection.UpdateOne(ctx, filter, update)
	if err != nil {

		fmt.Print(err)
		log.Println(err)
		return err
	}
	if updateResult.ModifiedCount == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

func (m mongodb) DeleteClient(ctx context.Context, id int) error {

	filter := bson.D{{"id", id}}

	deleteResult, err := m.clientsCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if deleteResult.DeletedCount == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
