package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/dmoyan0/Lab4/grpc" // proto
	"google.golang.org/grpc"
)

type DoshbankServer struct {
	montoAcumulado float64
}

// Función para obtener el monto acumulado
func (s *DoshbankServer) GetmontoAcumulado(ctx context.Context, req *pb.GetmontoAcumuladoRequest) (*pb.GetmontoAcumuladoResponse, error) {
	return &pb.GetmontoAcumuladoResponse{Monto: s.montoAcumulado}, nil
}

// Función que registra a un mercenario eliminado
func (s *DoshbankServer) MercenarioEliminado(ctx context.Context, req *pb.MercenarioEliminadoRequest) (*pb.MercenarioEliminadoResponse, error) {
	// Registra al mercenario eliminado y aumenta el monto acumulado cada vez que se llama esta función
	s.montoAcumulado += 100000000
	log.Printf("Mercenario %s eliminado en el piso %d. Monto acumulado: %f\n", req.Name, req.Piso, s.montoAcumulado)

	// Creamos e intentamos abrir  el txt
	file, err := os.OpenFile("mercenarios_eliminados.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error al abrir el archivo: %v", err)
	}
	defer file.Close()

	// Escritura en el txt con la info solicitada del mercenario
	line := fmt.Sprintf("Mercenario %s eliminado en el piso %d. Monto acumulado: %f\n", req.Name, req.Piso, s.montoAcumulado)
	if _, err := file.WriteString(line); err != nil {
		log.Fatalf("Error al escribir en el archivo: %v", err)
	}

	return &pb.MercenarioEliminadoResponse{}, nil
}

func main() {
	// Servidor gRPC
	server := grpc.NewServer()
	pb.RegisterDoshbankServer(server, &DoshbankServer{})

	// Iniciar el servidor gRPC
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("Doshbank gRPC server listening on port 50052...")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
