package main

import (
    "context"
    "log"
    "net"
    "google.golang.org/grpc"
    pb "https://github.com/dmoyan0/Lab4/blob/main/gRPC.proto"// Importa el paquete generado por protoc
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
