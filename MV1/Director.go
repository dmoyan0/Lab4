package main

import (
	"context"
	"log"
	"net"
	"math/rand"

	pb "github.com/dmoyan0/Lab4/grpc"

	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)

type DirectorServer struct {
	DoshBankClient  pb.DoshbankClient
	RabbitChannel   *amqp.Channel
	RabbitQueueName string
}


func (s *DirectorServer) MercenaryDecision(ctx context.Context, req *pb.MercenaryDecisionRequest) (*pb.MercenaryDecisionResponse, error) {
	// Implementa la lógica para manejar la decisión del mercenario
	switch req.Floor {
	case 1:
		// Piso 1 Entrada al infierno
		arma := req.armaChoice

		// Numeros aleatorios entre 0 y 100 
		rand.Seed(time.Now().UnixNano())
		X := rand.Intn(101)
		Y := rand.Intn(101)
		for X == Y {
			Y = rand.Intn(101)
		}//Caso que sean iguales, se elige otro

		// 	Probabilidades
		prob1 := float64(X)
		prob2 := float64(Y - X)
		prob3 := float64(100 - Y)

		// Determinar el resultado según el arma elegida
		switch arma {
		case 1: // Escopeta
			if prob1 >= 50 {
				// Mercenario vive
				log.Printf("El mercenario %s sobrevivió al infierno con la escopeta", req.Name)
				return &pb.MercenaryDecisionResponse{Message: "El mercenario sobrevivió al infierno con la escopeta"}, nil
			} else {
				// Eliminadp
				log.Printf("El mercenario %s no sobrevivió al infierno con la escopeta", req.Name)
				s.ReportarEliminacion(ctx, req.Name, req.Floor)
				return &pb.MercenaryDecisionResponse{Message: "El mercenario no sobrevivió al infierno con la escopeta"}, nil
			}
		case 2: // Rifle 
			if prob2 >= 50 {
				log.Printf("El mercenario %s sobrevivió al infierno con el rifle automático", req.Name)
				return &pb.MercenaryDecisionResponse{Message: "El mercenario sobrevivió al infierno con el rifle automático"}, nil
			} else {
				log.Printf("El mercenario %s no sobrevivió al infierno con el rifle automático", req.Name)
				s.ReportarEliminacion(ctx, req.Name, req.Floor)
				return &pb.MercenaryDecisionResponse{Message: "El mercenario no sobrevivió al infierno con el rifle automático"}, nil
			}
		case 3: // Puños 
			if prob3 >= 50 {
				log.Printf("El mercenario %s sobrevivió al infierno con los puños eléctricos", req.Name)
				return &pb.MercenaryDecisionResponse{Message: "El mercenario sobrevivió al infierno con los puños eléctricos"}, nil
			} else {
				log.Printf("El mercenario %s no sobrevivió al infierno con los puños eléctricos", req.Name)
				s.ReportarEliminacion(ctx, req.Name, req.Floor)
				return &pb.MercenaryDecisionResponse{Message: "El mercenario no sobrevivió al infierno con los puños eléctricos"}, nil
			}
		default:
			// Si la elección de arma no corresponde
			return nil, fmt.Errorf("arma no válida")
		}
	case 2:
		//Piso 2 
	case 3:
		//Piso 3
	}
}

func (s *DirectorServer) PlayerDecision(ctx context.Context, req *pb.MercenaryDecisionRequest) (*pb.MercenaryDecisionResponse, error) {
	var client pb.MercenaryClient

	conn, err := grpc.Dial("localhost:50055", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	player := pb.NewMercenaryClient(conn)

	


}

// Informa al DoshBank y a RabbitMQ cuando un mercenario es eliminado
func (s *DirectorServer) ReportarEliminacion(ctx context.Context, mercenario string, piso int) {
	// Informa al DoshBank
	_, err := s.DoshBankClient.MercenarioEliminado(ctx, &pb.MercenarioEliminadoRequest{
		Name: mercenario,
		Piso: piso,
	})
	if err != nil {
		log.Printf("Error al informar eliminación al DoshBank: %v", err)
	}

	// Publica un mensaje en RabbitMQ
	err = s.RabbitChannel.Publish(
		"",                // Exchange
		s.RabbitQueueName, // Key
		false,             // Mandatory
		false,             // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(mercenario),
		})
	if err != nil {
		log.Printf("Error al publicar mensaje en RabbitMQ: %v", err)
	}
}

// Se solicita al DoshBank el monto acumulado actual
func (s *DirectorServer) ObtenerMontoDoshBank(ctx context.Context) (float64, error) {
	// Envía una solicitud al DoshBank para obtener el monto acumulado
	resp, err := s.DoshBankClient.GetmontoAcumulado(ctx, &pb.GetmontoAcumuladoRequest{})
	if err != nil {
		log.Printf("Error al obtener monto acumulado del DoshBank: %v", err)
		return 0, err
	}
	// Retorna el monto acumulado recibido
	return resp.Monto, nil
}


func (s *DirectorServer) enviarDecisionANamenode(ctx context.Context, req *pb.MercenaryDecisionRequest) (*pb.MercenaryDecisionResponse, error) {
	//Se crea una conexion con el namenode
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error en la conexion: %v", err)
	}
	defer conn.Close()

	//Se crea el cliente para el servicio del namenode
	client := pb.NewClient(conn)

	//Registrar decision en el cliente
	_, err = Namenode.RegisterDecision(ctx, &pb.MercenaryDecisionRequest) {
		MercenaryName: req.name,
		Decision: req.decision,
		Floor: req.floor,
		IP: req.datanode_ip
	}
	if err != nil {
		log.Fatalf("Falla al registrar la decision en el datanode: %v", err)
	}

	return &pb.MercenaryDecisionResponse{Message: "Decision enviada al namenode"}
}

func main() {
	// Conexión con el servidor DoshBank
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error al conectar con el servidor DoshBank: %v", err)
	}
	defer conn.Close()

	// Cliente para el servicio DoshBank
	doshBankClient := pb.NewDoshbankClient(conn)

	// Conexión con RabbitMQ
	rabbitConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Error al conectar con RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()

	// Canal con RabbitMQ
	rabbitChannel, err := rabbitConn.Channel()
	if err != nil {
		log.Fatalf("Error al abrir el canal con RabbitMQ: %v", err)
	}
	defer rabbitChannel.Close()

	// Se declara la cola en RabbitMQ
	queue, err := rabbitChannel.QueueDeclare(
		"eliminacion_queue", // Nombre de la cola
		false,               // Durabilidad
		false,               // Eliminar al finalizar
		false,               // Exclusividad
		false,               // No esperar confirmación
		nil,                 // Argumentos adicionales
	)
	if err != nil {
		log.Fatalf("Error al declarar la cola en RabbitMQ: %v", err)
	}

	//Servidor Director
	directorServer := &DirectorServer{
		DoshBankClient:  doshBankClient,
		RabbitChannel:   rabbitChannel,
		RabbitQueueName: queue.Name,
	}

	// Se inicia el servidor gRPC Director
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDirectorServer(s, directorServer)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
