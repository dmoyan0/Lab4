package main

import (
	"context"
	"fmt"
	"log"
	"os"
	pb "https://github.com/dmoyan0/Lab4/blob/main/gRPC.proto"// Importa el paquete generado por protoc

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

func (s *datanodeServer) GetFileContent(ctx context.Context, req *pb.GetFileContentRequest) (res *pb.GetFileContentResponse, error) {
	filename := req.GetFileName()
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &pb.GetFileContentResponse{Content: string(fileContent)}, nil
}

func main() {
	// Establece la conexión gRPC con el Namenode
	conn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial Namenode: %v", err)
	}
	defer conn.Close()

	// Crea un cliente gRPC para el servicio Datanode
	client := pb.NewDatanodeClient(conn)
	//Falta establecer conexion con el servicio Namenode

	// Implementar un bucle que escuche por registros de decisiones
	// y los procese de alguna manera
}
