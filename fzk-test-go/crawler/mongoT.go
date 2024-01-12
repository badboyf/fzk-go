package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func main() {
}

func DeleteByIds(ids *[]string) {
	Exec(func(c *mongo.Collection) {
		filter := bson.M{"id": bson.M{"$in": ids}}
		many, err := c.DeleteMany(context.TODO(), filter)
		if err != nil {
			fmt.Printf("insert error: %v \n", err)
		}
		fmt.Printf("删除结果：%v \n", many)
	})
}
func InsertBatch(v interface{}) {
	var realV = v.([]interface{})
	Exec(func(c *mongo.Collection) {
		many, err := c.InsertMany(context.TODO(), realV)
		if err != nil {
			fmt.Printf("insert error: %v \n", err)
		}
		fmt.Printf("插入结果：%v \n", many)
	})
}


func Exec(f func(c *mongo.Collection)) {
	param := fmt.Sprintf("mongodb://127.0.0.1:27017")
	clientOptions := options.Client().ApplyURI(param)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println(err)
	}

	collection := client.Database("golangcc").Collection("articleCrawler")
	f(collection)

	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

func conn()  {
	param := fmt.Sprintf("mongodb://127.0.0.1:27017")
	clientOptions := options.Client().ApplyURI(param)

	// 建立客户端连接
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println(err)
	}

	// 检查连接情况
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Connected to MongoDB!")

	//指定要操作的数据集
	collection := client.Database("golangcc").Collection("articleCrawler")

	//执行增删改查操作
	fmt.Printf("%v", collection)

	// 断开客户端连接
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}
