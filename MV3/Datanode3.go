package main

import (
    "context"
    "log"
    "google.golang.org/grpc"
    "net"

    pb "https://github.com/dmoyan0/Lab4/blob/main/gRPC.proto"// Importa el paquete generado por protoc
)

type datanodeServer struct {}

func (s *datanodeServer) RegisterDecision(ctx context.Context, req *pb.RegisterDecisionRequest) (*pb.RegisterDecisionResponse, error) {
    //Se extrae los datos del registro recibido
    mercenaryName := req.mercenary_name
    decision := req.decision
    floor := req.floor

    fileName := fmt.Sprintf("%s_%s.txt", mercenaryName, floor)
    file, err := os.Create(fileName)
    if err != nil {
        log.Fatalf("Error creating file %v:", err)
    }
    defer file.Close()

    _,err = file.WriteString(decision)
    if err != nil {
        log.Fatalf("Failed to write decision")
    }

    return &pb.RegisterDecisionResponse{"Decision escrita en un archivo correctamente"}, nil
}

func main() {
    // Establece la conexión gRPC con el Namenode
    conn, err := grpc.Dial("localhost:50054", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to dial Namenode: %v", err)
    }
    defer conn.Close()

    // Crea un cliente gRPC para el servicio Namenode
    client := pb.NewClient(conn)

    // Implementar un bucle que escuche por registros de decisiones
    // y los procese de alguna manera
}
