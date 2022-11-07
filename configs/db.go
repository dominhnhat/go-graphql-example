package configs
import (
	"context"
	"fmt"
	"log"
	"go-graphql/graph/model"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
type DB struct {
	client *mongo.Client
}
func ConnectDB() *DB {
	client, err := mongo.NewClient(options.Client().ApplyURI(EnvMongoURI()))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")
	return &DB{client: client}
}
func colHelper(db *DB, collectionName string) *mongo.Collection {
	return db.client.Database("test").Collection(collectionName)
}
func (db *DB) CreateLink(input *model.NewLink) (*model.Link, error) {
	collection := colHelper(db, "link")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, input)

	if err != nil {
		return nil, err
	}

	link := &model.Link{
		ID:          res.InsertedID.(primitive.ObjectID).Hex(),
		Title: input.Title,
		Address: input.Address,
	}

	return link, err
}
func (db *DB) GetLink() ([]*model.Link, error) {
	collection := colHelper(db, "link")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var links []*model.Link
	defer cancel()

	res, err := collection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	defer res.Close(ctx)
	for res.Next(ctx) {
		var singleLink *model.Link
		if err = res.Decode(&singleLink); err != nil {
			log.Fatal(err)
		}
		links = append(links, singleLink)
	}

	return links, err
}