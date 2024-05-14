package main

import (
	"context"
	pb "https://github.com/dmoyan0/Lab4/blob/main/gRPC.proto" // Importa el paquete generado por protoc
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

// Implementa el método MercenaryDecision del servicio Director
func (s *server) MercenaryDecision(ctx context.Context, req *pb.MercenaryDecisionRequest) (*pb.MercenaryDecisionResponse, error) {
	//implementar la lógica para manejar la decisión del mercenario
	log.Printf("Received decision from mercenary %s: %s\n", req.Name, req.Decision)

	// Retorna una respuesta indicando que la decisión se procesó correctamente
	return &pb.MercenaryDecisionResponse{Message: "Respuesta recibida correctamente"}, nil
}

// Implementa el método GetAccumulatedAmount del servicio Director
func (s *server) GetAccumulatedAmount(ctx context.Context, req *pb.GetAccumulatedAmountRequest) (*pb.GetAccumulatedAmountResponse, error) {
	// implementar la lógica para obtener el monto acumulado

	// Retorna una respuesta con el monto acumulado
	return &pb.GetAccumulatedAmountResponse{Amount: amount}, nil
}

// ReportarEliminacion informa al DoshBank cuando un mercenario es eliminado
func (s *DirectorServer) ReportarEliminacion(ctx context.Context, mercenario string) {
	_, err := s.DoshBankClient.Eliminado(ctx, &pb.EliminadoRequest{
		Mercenario: mercenario,
		Piso:       0, //*
	})
	if err != nil {
		log.Printf("Error al informar eliminación al DoshBank: %v", err)
	}
}

// ObtenerMontoDoshBank obtiene el monto acumulado actual del DoshBank
func (s *DirectorServer) ObtenerMontoDoshBank(ctx context.Context) int {
	resp, err := s.DoshBankClient.ObtenerMonto(ctx, &pb.Empty{})
	if err != nil {
		log.Printf("Error al obtener monto acumulado del DoshBank: %v", err)
		return 0
	}
	return int(resp.TotalAmount)
}

func main() {
	lis, err := net.Listen("tcp", ":")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDirectorServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
