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
	usersCollection   *mongo.Collection
}

func InitConnection() *mongodb {

	host := utils.Conf.Get("mongodb.host")
	port := utils.Conf.GetInt("mongodb.port")
	database := utils.Conf.GetString("mongodb.database")

	uri := fmt.Sprintf("mongodb://%s:%d", host, port)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	mongoDatabase := client.Database(database)
	mongoClientsCollection := mongoDatabase.Collection("clients")
	mongoUsersCollection := mongoDatabase.Collection("users")

	if err != nil {
		log.Fatal(err)
	}

	return &mongodb{client: client, database: mongoDatabase, clientsCollection: mongoClientsCollection, usersCollection: mongoUsersCollection}
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

func (m mongodb) List(ctx context.Context, offset int, limit int) ([]models.User, error) {
	var users []models.User

	filter := bson.D{}
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset)).SetProjection(bson.D{{"roles", 0}})

	cursor, err := m.usersCollection.Find(ctx, filter, opts)
	if err != nil {
		log.Println(err)
		return users, err
	}

	err = cursor.All(ctx, &users)
	if err != nil {
		log.Println(err)
		return users, err
	}

	return users, nil
}

func (m mongodb) ById(ctx context.Context, id int) (models.User, error) {

	filter := bson.D{{"id", id}}

	var user models.User
	opts := options.FindOne().SetProjection(bson.D{{"roles", 0}})
	err := m.usersCollection.FindOne(ctx, filter, opts).Decode(&user)
	if err != nil {
		log.Println(err)
		return models.User{}, err
	}

	return user, nil
}

func (m mongodb) ByEmail(ctx context.Context, email string) (models.User, error) {
	filter := bson.D{{"email", email}}

	var user models.User
	opts := options.FindOne().SetProjection(bson.D{{"roles", 0}})
	err := m.usersCollection.FindOne(ctx, filter, opts).Decode(&user)
	if err != nil {
		log.Println(err)
		return models.User{}, err
	}

	return user, nil
}

func (m mongodb) CreateUser(ctx context.Context, newUser models.User) (models.User, error) {

	_, err := m.clientsCollection.InsertOne(ctx, newUser)
	if err != nil {
		log.Println(err)
		return models.User{}, err
	}

	return newUser, nil
}

func (m mongodb) UpdateUser(ctx context.Context, user models.User) (models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m mongodb) DeleteUser(ctx context.Context, id int) (models.User, error) {

	filter := bson.D{{"id", id}}

	var user models.User
	opts := options.FindOneAndDelete().SetProjection(bson.D{{"roles", 0}})
	err := m.clientsCollection.FindOneAndDelete(ctx, filter, opts).Decode(&user)

	if err != nil {
		return user, err
	}
	if err == mongo.ErrNoDocuments {
		return user, errors.New("no rows affected")
	}

	return user, nil
}

func (m mongodb) UpdateRoles(ctx context.Context, user models.User, roles []models.Role) (models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m mongodb) CheckUserGrant(ctx context.Context, userId int, table string, operation string) (found bool, err error) {
	operation = "$$item." + operation
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"id", userId}}}},
		{{"$lookup", bson.D{
			{"from", "roles"},
			{"localField", "roles"},
			{"foreignField", "name"},
			{"as", "roles"},
		}}},
		{{"$project", bson.D{{
			"grants", bson.D{{
				"$reduce", bson.D{
					{"input", "$roles.grants"},
					{"initialValue", bson.A{}},
					{"in", bson.D{{"$concatArrays", bson.A{"$$value", "$$this"}}}},
				},
			}},
		}}}},
		{{"$project", bson.D{{"grant", bson.D{{
			"$filter", bson.D{
				{"input", "$grants"},
				{"as", "item"},
				{"cond", bson.D{{
					"$and", bson.A{
						bson.D{{"$eq", bson.A{"$$item.table", table}}},
						bson.D{{"$eq", bson.A{operation, true}}},
					}},
				}},
			},
		}}}}}},
	}

	cursor, err := m.usersCollection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Println(err)
		return false, err
	}

	var grants []models.Grant
	err = cursor.All(ctx, &grants)
	if err != nil {
		log.Println(err)
		return false, err
	}

	if len(grants) == 0 {
		return false, nil
	}

	return true, nil
}

func (m mongodb) GetFieldsPermissions(ctx context.Context, userId int, tableName string) (map[string]bool, error) {
	pipeline := mongo.Pipeline{
		bson.D{{"$match", bson.D{{"id", userId}}}},
		bson.D{{"$lookup", bson.D{
			{"from", "roles"},
			{"localField", "roles"},
			{"foreignField", "name"},
			{"as", "roles"},
		}}},
		bson.D{{"$project", bson.D{
			{"_id", 0},
			{"fields", bson.D{
				{"$reduce", bson.D{
					{"input", "$roles.fields"},
					{"initialValue", bson.A{}},
					{"in", bson.D{{"$concatArrays", bson.A{"$$value", "$$this"}}}},
				}},
			}},
		}}},
		bson.D{{"$unwind", bson.D{{"path", "$fields"}}}},
	}
	cursor, err := m.usersCollection.Aggregate(context.TODO(), pipeline)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	result := struct {
		Fields map[string]map[string]bool `json:"fields"`
	}{}

	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&result)
		if err != nil {
			fmt.Println(err)
		}
	}
	return result.Fields[tableName], nil
}
