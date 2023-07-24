package configs

import (
    "context"
    "fmt"
    "log"
    "time"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// Connects Mongo DB using EnvMongoURI function while also disconnecting 
// when there is trouble connecting within 10 seconds
func ConnectDB() *mongo.Client  {
	// context that has a 10 second timeout connected to it
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// use 'client' as the variable that connects the DB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(EnvMongoURI()))
  if err != nil {
      log.Fatal(err)
  }
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

  //ping the database
  err = client.Ping(ctx, nil)
	if err != nil {
    log.Fatal(err)
  }

  fmt.Println("Connected to MongoDB")
  return client
}

//Client instance
var DB *mongo.Client = ConnectDB()

//getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
    collection := client.Database("golangAPI").Collection(collectionName)
    return collection
}