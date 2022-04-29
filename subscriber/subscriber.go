package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type Game struct {
	Juegoid        string `json:"juegoid"`
	Cantjugadores  string `json:"cantjugadores"`
	Nombrejuego    string `json:"nombrejuego"`
	Jugadorganador int    `json:"jugadorganador"`
	Queue          string `json:"queue"`
}

type GameRedis struct {
	Nombre_Juego   string `json:"nombre_juego"`
	Nombre_Ganador int    `json:"nombre_ganador"`
}

type GamesTiDB struct {
	Nombrejuego    string `json:"nombrejuego"`
	Jugadorganador int    `json:"jugadorganador"`
}

var ADDRMONGO = "mongodb://admindb:1234@35.209.156.204:27017"
var NAMEDB = "Practica2"
var NAMECOLL = "Practica2"

func main() {

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "my-cluster-kafka-bootstrap.kafka:9092",
		"group.id":          "my-topic",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Consumer Startup... %v\n", c)

	err = c.SubscribeTopics([]string{"my-topic"}, nil)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to Suscribed Topic: %s\n", err)
		os.Exit(1)
	}
	run := true

	for run {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %+v\n", msg.TopicPartition, string(msg.Value))
			//MONGO
			SaveLogMongo(msg.Value)

		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	fmt.Printf("Closing consumer\n")
	c.Close()
}

func FromJSON(data []byte) GameRedis {
	game := GameRedis{}
	err := json.Unmarshal(data, &game)
	if err != nil {
		panic(err)
	}
	return game
}

func FromJSON2(data []byte) GamesTiDB {
	game := GamesTiDB{}
	err := json.Unmarshal(data, &game)
	if err != nil {
		panic(err)
	}
	return game
}

func SaveLogMongo(bjsonLog []byte) {

	// Get Struct Log
	logsobj := Game{}
	err := json.Unmarshal(bjsonLog, &logsobj)
	if err != nil {
		fmt.Println("Error al decodificar")
		return
	}

	// Database connection
	client, err := mongo.NewClient(options.Client().ApplyURI(ADDRMONGO))
	if err != nil {
		fmt.Println("Error de Conexion-1")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		fmt.Println("Error de Conexion-2")
		return
	}
	defer client.Disconnect(ctx)

	dbLogs := client.Database(NAMEDB).Collection(NAMECOLL)

	// Insert Value
	_, err = dbLogs.InsertOne(ctx, logsobj)
	if err != nil {
		fmt.Println("Not inserted")
		fmt.Println(err)
	} else {
		fmt.Println("New data inserted")
	}
}
