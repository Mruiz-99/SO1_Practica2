package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"

	pb "github.com/racarlosdavid/demo-gRPC/proto"
	"google.golang.org/grpc"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type server struct {
	pb.UnimplementedServicioGolangServer
}

type Games struct {
	Juegoid        string `json:"juegoid"`
	Cantjugadores  string `json:"cantjugadores"`
	Nombrejuego    string `json:"nombrejuego"`
	Jugadorganador int    `json:"jugadorganador"`
	Queue          string `json:"queue"`
}

var Ganador int

func (s *server) IniciarJuego(ctx context.Context, in *pb.PlayerGameRequest) (*pb.PlayerGameReply, error) {
	//log.Printf("El id del juego es %v con la cantidad de %v jugadores", in.GetGame(), in.GetPlayers())

	cant_Jugadores, err := strconv.Atoi(in.GetPlayers())
	if err != nil {
		log.Fatalf("Error al convertir a int: %v", err)
	}

	num_Juego, err := strconv.Atoi(in.GetGame())
	if err != nil {
		log.Fatalf("Error al convertir a int: %v", err)
	}

	if num_Juego == 1 {
		Ganador = Juego1(cant_Jugadores)
		juego := new(Games)
		juego.Juegoid = "1"
		juego.Cantjugadores = in.GetPlayers()
		juego.Nombrejuego = "Juego_Random"
		juego.Jugadorganador = Ganador
		juego.Queue = "queue_kafka"
		sendKafka(juego)

	} else if num_Juego == 2 {
		Ganador = Juego2(cant_Jugadores)
		juego := new(Games)
		juego.Juegoid = "2"
		juego.Cantjugadores = in.GetPlayers()
		juego.Nombrejuego = "Pelea_Impares_Pares"
		juego.Jugadorganador = Ganador
		juego.Queue = "queue_kafka"
		sendKafka(juego)

	} else if num_Juego == 3 {
		Ganador = Juego3(cant_Jugadores)
		juego := new(Games)
		juego.Juegoid = "3"
		juego.Cantjugadores = in.GetPlayers()
		juego.Nombrejuego = "Ruleta_no_rusa"
		juego.Jugadorganador = Ganador
		juego.Queue = "queue_kafka"
		sendKafka(juego)

	} else if num_Juego == 4 {
		Ganador = Juego4(cant_Jugadores)
		juego := new(Games)
		juego.Juegoid = "4"
		juego.Cantjugadores = in.GetPlayers()
		juego.Nombrejuego = "La_posicion_de_la_Suerte"
		juego.Jugadorganador = Ganador
		juego.Queue = "queue_kafka"
		sendKafka(juego)

	} else if num_Juego == 5 {
		Ganador = Juego5(cant_Jugadores)
		juego := new(Games)
		juego.Juegoid = "5"
		juego.Cantjugadores = in.GetPlayers()
		juego.Nombrejuego = "La_Ultima_Bala"
		juego.Jugadorganador = Ganador
		juego.Queue = "queue_kafka"
		sendKafka(juego)

	} else {
		Ganador = 0
	}

	strResultado := strconv.Itoa(Ganador)
	//log.Printf("Received: %v", in.GetGame())
	return &pb.PlayerGameReply{Mensajeganador: "El Jugador ganador del Juego " + in.GetGame() + " es: " + strResultado}, nil
}

func main() {
	lis, err := net.Listen("tcp", PORT())

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterServicioGolangServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func PORT() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "50051"
	}
	return ":" + port
}

//Juego_Random
func Juego1(LimJugadores int) int {

	var ale int

	ale = rand.Intn(LimJugadores) + 1

	//fmt.Print(ale)

	return ale
}

//Pelea_Impares_Pares
func Juego2(LimJugadores int) int {

	var ale int

	num1 := 1

	for num1 != 0 {
		ale = rand.Intn(LimJugadores) + 1
		fmt.Println(ale)
		num1 = (ale % 2)
		fmt.Println(num1)
	}

	par := ale

	num2 := 0

	for num2 == 0 {
		ale = rand.Intn(LimJugadores) + 1
		fmt.Println(ale)
		num2 = (ale % 2)
		fmt.Println(num2)
	}

	impar := ale

	if par > impar {
		return par
	} else {
		return impar
	}
}

//Ruleta_no_rusa
func Juego3(LimJugadores int) int {

	var ale int
	var ruleta []int

	for i := 0; i <= 5; i++ {
		ale = rand.Intn(LimJugadores) + 1
		//fmt.Println(ale)
		ruleta = append(ruleta, ale)
	}

	fmt.Println(ruleta)

	ale = rand.Intn(5)

	fmt.Println(ale)

	return ruleta[ale]
}

//La_posicion_de_la_Suerte
func Juego4(LimJugadores int) int {

	var ale int
	var Lista []int

	for i := 0; i <= 9; i++ {
		ale = rand.Intn(LimJugadores) + 1
		//fmt.Println(ale)
		Lista = append(Lista, ale)
	}

	fmt.Println(Lista)

	return Lista[6]
}

//La_Ultima_Bala
func Juego5(LimJugadores int) int {

	var ale int
	var Lista []int

	for i := 11; i >= 0; i-- {
		ale = rand.Intn(LimJugadores) + 1
		//fmt.Println(ale)
		Lista = append(Lista, ale)
	}

	fmt.Println(Lista)

	return Lista[11]
}

func sendKafka(word *Games) {

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": brokerAddr()})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	topic := "my-topic"

	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          word.ToJSON(),
	}, nil)
	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
}

func brokerAddr() string {
	brokerAddr := os.Getenv("BROKER_ADDR")
	if len(brokerAddr) == 0 {
		brokerAddr = "my-cluster-kafka-bootstrap.kafka:9092"
	}
	return brokerAddr
}
func (g Games) ToJSON() []byte {
	ToJSON, err := json.Marshal(g)
	if err != nil {
		panic(err)
	}
	return ToJSON
}
