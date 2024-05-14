package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/dmoyan0/Lab4/tree/main/proto"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Servidor DoshBank
type DoshBankServer struct {
	MontoAcumulado int
}

// Servicio gRPC DoshBank
type DoshBankService struct {
	DoshBankServer *DoshBankServer
}

// Registra la muerte y el monto acumulado del mercenario
func (s *DoshBankService) Eliminado(ctx context.Context, req *pb.EliminadoRequest) (*empty.Empty, error) {
	log.Printf("Mercenario eliminado: %s en piso %d\n", req.Mercenario, req.Piso)
	s.DoshBankServer.MontoAcumulado += 10 // Aumenta el monto acumulado **
	log.Printf("Monto acumulado actualizado: %d\n", s.DoshBankServer.MontoAcumulado)
	s.ArchivoEliminado(req.Mercenario, req.Piso) // Archivo
	return &empty.Empty{}, nil
}

// Entrega el monto acumulado en el momento de la consulta
func (s *DoshBankService) ObtenerMonto(ctx context.Context, req *empty.Empty) (*pb.MontoAcumuladoResponse, error) {
	log.Printf("Consulta de monto acumulado recibida\n")
	return &pb.MontoAcumuladoResponse{TotalAmount: int32(s.DoshBankServer.MontoAcumulado)}, nil
}

// Trabajo de archivo en caso que muere el mercenario
func (s *DoshBankService) ArchivoEliminado(mercenario string, piso int32) {
	file, err := os.OpenFile("eliminado.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error al abrir el archivo: %v", err)
	}
	defer file.Close()

	if _, err := fmt.Fprintf(file, "%s Piso %d %d\n", mercenario, piso, s.DoshBankServer.MontoAcumulado); err != nil {
		log.Fatalf("Error al escribir en el archivo: %v", err)
	}
}

func main() {
	// Configuración del servidor gRPC
	server := &DoshBankServer{MontoAcumulado: 0}
	grpcServer := grpc.NewServer()
	doshBankService := &DoshBankService{DoshBankServer: server}
	pb.RegisterDoshBankServiceServer(grpcServer, doshBankService)

	// Comunicación con RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/") //*
	if err != nil {
		log.Fatalf("Error al conectar con RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error al abrir el canal: %v", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"cola_eliminados", // Nombre de la cola
		true,              // Durable
		false,             // AutoDelete
		false,             // Exclusive
		false,             // NoWait
		nil,               // Argumentos
	)
	if err != nil {
		log.Fatalf("Error al declarar la cola: %v", err)
	}

	msgs, err := ch.Consume(
		"cola_eliminados", // Nombre de la cola
		"",                // Consumer
		true,              // AutoAck
		false,             // Exclusive
		false,             // NoLocal
		false,             // NoWait
		nil,               // Argumentos
	)
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}

	go func() {
		for msg := range msgs {
			log.Printf("Mensaje recibido: %s\n", msg.Body)
		}
	}()

	// Terminar el servidor gRPC
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar el servidor gRPC
	reflection.Register(grpcServer)
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Error al servir: %v", err)
		}
	}()

	// Espera la señal y terminar
	<-stop
	log.Println("Deteniendo el servidor gRPC")
	grpcServer.GracefulStop()
}
