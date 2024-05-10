package main

import (
    "context"
    "log"
    "google.golang.org/grpc"
    pb "path_to_your_proto_file" // Importa el paquete generado por protoc
)

func main() {
    // Establece la conexión gRPC con el Namenode
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to dial Namenode: %v", err)
    }
    defer conn.Close()

    // Crea un cliente gRPC para el servicio Namenode
    client := pb.NewNamenodeClient(conn)

    // Implementar un bucle que escuche por registros de decisiones
    // y los procese de alguna manera
}
