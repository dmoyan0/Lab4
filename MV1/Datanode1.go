package main

import (
	"context"
	"fmt"
	"log"
	"os"
	pb "path_to_your_proto_file" // Importa el paquete generado por protoc

	"google.golang.org/grpc"
)

type datanodeServer struct{}

func (s *datanodeServer) RegisterDecision(ctx context.Context, req *pb.RegisterDecisionRequest) (*pb.RegisterDecisionResponse, error) {
    //Se extrae los datos del registro recibido
    mercenaryName := req.mercenary_name
    decision := req.decision
    floor := req.floor



    fileName := fmt.Sprintf("%s_%s.txt", mercenaryName, floor)
    
    //desarrollar si es que existe el documento se debe escribir sobre este y no crear uno nuevo
    if _, err := os.Stat(fileName); err == nil {
        //Si existe el archivo, escribir sobre este
        file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0666)
        if err != nil {
            return nil, err
        }
        defer file.Close()

        _, err = file.WriteString(decision)
        if err != nil {
            return nil, err
        }
    }
    else {
        file, err := os.Create(fileName)
        if err != nil {
            log.Fatalf("Error creating file %v:", err)
        }
        defer file.Close()

        _,err = file.WriteString(decision)
        if err != nil {
            log.Fatalf("Failed to write decision")
        }
    }
    

    return &pb.RegisterDecisionResponse{"Decision escrita en un archivo correctamente"}, nil
}

func main() {
	// Establece la conexión gRPC con el Namenode
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial Namenode: %v", err)
	}
	defer conn.Close()

	// Crea un cliente gRPC para el servicio Namenode
	client := pb.NewNamenodeClient(conn)

	// Implementar un bucle que escuche por registros de decisiones
	// y los procese de alguna manera
}
