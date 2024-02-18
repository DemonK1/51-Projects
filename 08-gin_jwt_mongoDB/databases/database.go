package databases

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var Client = DBInstance()

// DBInstance 用于创建并返回 MongoDB 客户端对象，以便在代码中进行 MongoDB 数据库操作
func DBInstance() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	MongoDB := os.Getenv("MONGODB_URL")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MongoDB))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")

	return client
}

// OpenCollection 用于从给定的 MongoDB 客户端对象中获取指定名称的集合对象，以便在代码中进行集合操作。
func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection = client.Database("cluster0").Collection(collectionName)
	return collection
}
