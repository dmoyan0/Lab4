package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	pb "github.com/dmoyan0/Lab4/blob/main/gRPC.proto" // proto

	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)

// Servidor del Director
type DirectorServer struct {
	DoshBankClient  pb.DoshbankClient
	RabbitChannel   *amqp.Channel
	RabbitQueueName string
}

// Se inicializa el cliente Doshbank
func InitDoshbankClient() (pb.DoshbankClient, error) {
	conn, err := grpc.Dial("localhost:50056", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return pb.NewDoshbankClient(conn), nil
}

// Se inicializa el canal RabbitMQ y se declara la cola
func InitRabbitMQChannel() (*amqp.Channel, string, error) {
	//Establece conexión con el servidor rabbitMQ en el puerto 5672
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, "", err
	}
	ch, err := conn.Channel() //Si la conexión es exitosa -> Crea un canal rabbit para crear la cola
	if err != nil {
		return nil, "", err
	}

	//Se declara la cola para los mensajescon el Doshbank
	queue, err := ch.QueueDeclare(
		"eliminacion_queue", // Nombre de la cola
		false,               // Durabilidad
		false,               // Eliminar al finalizar
		false,               // Exclusividad
		false,               // No esperar confirmación
		nil,                 // Argumentos adicionales
	)
	if err != nil {
		return nil, "", err
	}

	return ch, queue.Name, nil
}

// Función para recibir el monto acumulado actual que entrega el doshbank
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

// Función que informa al DoshBank y a RabbitMQ cuando un mercenario es eliminado
func (s *DirectorServer) ReportarEliminacion(ctx context.Context, mercenario string, piso int) {
	// Informa al DoshBank
	_, err := s.DoshBankClient.MercenarioEliminado(ctx, &pb.MercenarioEliminadoRequest{
		Name: mercenario,
		Piso: piso,
	})
	if err != nil {
		log.Printf("Error al informar eliminación al DoshBank: %v", err)
	}

	mensaje := fmt.Sprintf("Mercenario %s eliminado en el piso %d", mercenario, piso)

	// Publica un mensaje en RabbitMQ
	err = s.RabbitChannel.Publish(
		"",                // Exchange
		s.RabbitQueueName, // Key
		false,             // Mandatory
		false,             // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(mensaje),
		})
	if err != nil {
		log.Printf("Error al publicar mensaje en RabbitMQ: %v", err)
	}
}

func (s *DirectorServer) MercenaryDecision(ctx context.Context, req *pb.MercenaryDecisionRequest) (*pb.MercenaryDecisionResponse, error) {
	var readyGame string
	fmt.Print("¿Dar inicio al juego? [Si/No]: ")
	fmt.Scanf("%s", &readyGame)

	if readyGame == "Si" {
		// Lógica para manejar la decisión del mercenario dependiendo del piso
		switch req.Floor {

		case 1:
			var readyPiso1 string
			fmt.Print("¿Dar inicio al Piso 1? [Si/No]: ")
			fmt.Scanf("%s", &readyPiso1)

			if readyPiso1 == "Si" {
				// Piso 1 Entrada al infierno
				arma := req.armaChoice

				// Numeros aleatorios entre 0 y 100
				rand.Seed(time.Now().UnixNano())
				X := rand.Intn(101)
				Y := rand.Intn(101)
				for X == Y {
					Y = rand.Intn(101)
				} //Caso que sean iguales, se elige otro

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
						s.enviarDecisionANamenode(ctx, req)
						return &pb.MercenaryDecisionResponse{Message: "El mercenario sobrevivió al infierno con la escopeta", estado: true}, nil
					} else {
						// Eliminadp
						log.Printf("El mercenario %s no sobrevivió al infierno con la escopeta", req.Name)
						s.ReportarEliminacion(ctx, req.Name, req.Floor)
						s.enviarDecisionANamenode(ctx, req)
						return &pb.MercenaryDecisionResponse{Message: "El mercenario no sobrevivió al infierno con la escopeta", estado: false}, nil
					}
				case 2: // Rifle
					if prob2 >= 50 {
						log.Printf("El mercenario %s sobrevivió al infierno con el rifle", req.Name)
						s.enviarDecisionANamenode(ctx, req)
						return &pb.MercenaryDecisionResponse{Message: "El mercenario sobrevivió al infierno con el rifle", estado: true}, nil
					} else {
						log.Printf("El mercenario %s no sobrevivió al infierno con el rifle", req.Name)
						s.ReportarEliminacion(ctx, req.Name, req.Floor)
						s.enviarDecisionANamenode(ctx, req)
						return &pb.MercenaryDecisionResponse{Message: "El mercenario no sobrevivió al infierno con el rifle", estado: false}, nil
					}
				case 3: // Puños
					if prob3 >= 50 {
						log.Printf("El mercenario %s sobrevivió al infierno con los puños eléctricos", req.Name)
						s.enviarDecisionANamenode(ctx, req)
						return &pb.MercenaryDecisionResponse{Message: "El mercenario sobrevivió al infierno con los puños eléctricos", estado: true}, nil
					} else {
						log.Printf("El mercenario %s no sobrevivió al infierno con los puños eléctricos", req.Name)
						s.ReportarEliminacion(ctx, req.Name, req.Floor)
						s.enviarDecisionANamenode(ctx, req)
						return &pb.MercenaryDecisionResponse{Message: "El mercenario no sobrevivió al infierno con los puños eléctricos", estado: false}, nil
					}
				default:
					// Si la elección de arma no corresponde
					return nil, fmt.Errorf("arma no válida")
				}
			} else {
				log.Printf("El director no ha comenzado con el piso 1")
			}

		case 2:
			var readyPiso2 string
			fmt.Print("¿Dar inicio al Piso 2? [Si/No]: ")
			fmt.Scanf("%s", &readyPiso2)

			if readyPiso2 == "Si" {
				//Piso 2: Trampas y traiciones
				decisionMercenario := req.Decision //Estas pueden ser o el camino A o B

				rand.Seed(time.Now().UnixNano())
				decisionDirector := ""
				if rand.Intn(2) == 0 {
					decisionDirector = "A"
				} else {
					decisionDirector = "B"

				}

				//Comparamos la decision del mercenario vs la del director
				if decisionMercenario == decisionDirector {
					log.Printf("El mercenario %s eligio eligio el pasillo %s y pasa al piso final!", req.Name, &decisionMercenario)
					s.enviarDecisionANamenode(ctx, req)
					return &pb.MercenaryDecisionResponse{Mensaje: fmt.Sprintf("El mercesario eligioel pasillo %s y pasa al piso final", &decisionMercenario), estado: true}, nil
				} else {
					log.Printf("El mercenario %s eligio eligio el pasillo %s y quedo eliminado", req.Name, &decisionMercenario)
					s.enviarDecisionANamenode(ctx, req)
					return &pb.MercenaryDecisionResponse{Mensaje: fmt.Sprintf("El mercesario eligioel pasillo %s y quedo eliminado", &decisionMercenario), estado: false}, nil
				}
			} else {
				log.Printf("El director no ha comenzado con el piso 2")
			}
		case 3:
			var readyPiso3 string
			fmt.Print("¿Dar inicio al Piso 3? [Si/No]: ")
			fmt.Scanf("%s", &readyPiso3)

			if readyPiso3 == "Si" {
				//Piso 3: Confrontación Final
				rand.Seed(time.Now().UnixNano())
				patriarcaNumero := rand.Intn(15) + 1

				mercenarioNumero := req.Decision
				aciertos := req.Aciertos

				// Comparamos el número del mercenario con el número del patriarca
				if mercenarioNumero == int(patriarcaNumero) {
					aciertos++
				}

				if aciertos >= 2 {
					log.Printf("El mercenario %s acerto dos o más veces el número del Patriarca, gano! %d ", req.Name, patriarcaNumero)
					s.enviarDecisionANamenode(ctx, req)
					return &pb.MercenaryDecisionResponse{
						Mensaje:  fmt.Sprintf("El mercesario ha salido victorioso "),
						estado:   true,
						Aciertos: true,
					}, nil
				} else {
					log.Printf("El mercenario %s no acerto dos o más veces el número del Patriarca %d ", req.Name, patriarcaNumero)
					s.enviarDecisionANamenode(ctx, req)
					return &pb.MercenaryDecisionResponse{
						Mensaje:  fmt.Sprintf("El mercesario ha sido eliminado "),
						estado:   false,
						Aciertos: false,
					}, nil

				}
			} else {
				log.Printf("El director no ha comenzado con el piso 3")
			}
		default:
			return nil, fmt.Errorf("Piso invalido")
		}

	} else {
		log.Printf("El Director ha decidido no iniciar el juego")
	}

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

	conn, err = grpc.Dial("localhost:50055", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failer to connect: %v", nil)
	}
	defer conn.Close()

	mercenary := pb.NewMercenaryClient(conn)

	//Registrar decision en el cliente
	_, err = Namenode.RegisterDecision(ctx, &pb.MercenaryDecisionRequest){
		MercenaryName: req.name,
		Decision:      req.decision,
		Floor:         req.floor,
		IP:            req.datanode_ip,
	}
	if err != nil {
		log.Fatalf("Falla al registrar la decision en el datanode: %v", err)
	}

	return &pb.MercenaryDecisionResponse{Message: "Decision enviada al namenode"}
}

func main() {
	//DoshBank y RabbitMQ
	doshBankClient, err := InitDoshbankClient() //Se llama función para inicializar cliente Doshbank
	if err != nil {
		log.Fatalf("Error al conectar con  DoshBank: %v", err)
	}
	defer doshBankClient.Close()

	rabbitChannel, rabbitQueueName, err := InitRabbitMQChannel() //Se llama función para inicializar el canal RabbitMq
	if err != nil {
		log.Fatalf("Error al conectar con RabbitMQ: %v", err)
	}
	defer rabbitChannel.Close()

	directorServer := &DirectorServer{
		DoshBankClient:  doshBankClient,
		RabbitChannel:   rabbitChannel,
		RabbitQueueName: rabbitQueueName,
	}

	//Se iniciliza el servidor gRCP del director en el puerto 50053
	lis, err := net.Listen("tcp", ":50050")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDirectorServer(s, directorServer)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
