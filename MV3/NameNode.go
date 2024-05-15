package main

import (
    "context"
    "log"
    "net"
    "google.golang.org/grpc"
	"math/rand"
	"time"

    pb "https://github.com/dmoyan0/Lab4/blob/main/gRPC.proto"
)

type server struct{}

// Implementa el método RegisterDecision del servicio Namenode
func (s *server) RegisterDecision(ctx context.Context, req *pb.RegisterDecisionRequest) (*pb.RegisterDecisionResponse, error) {
	datanode_ip := req.datanode_ip

	//Establece una conexion gRPC con el Datanode elegido
	conn, err := net.Dial(datanode_ip, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error al conectarse: %v", err)
	}
	defer conn.Close()

	//Crea un cliente gRPC para el servicio Datanode
	client := pb.NewClient(conn)

	_,err = client.RegisterDecision(ctx, &pb.RegisterDecisionRequest{
		MercenaryName: req.name,
		Decision: req.decision,
		Floor: req.floor,
		IP: req.datanode_ip
	})
	if err != nil {
		log.Fatalf("Error al enviar el registro al Datanode": %v, err)
	}
    // Retorna una respuesta indicando que el registro se procesó correctamente
    return &pb.RegisterDecisionResponse{Message: "Decision registrada"}, nil
}

func getFileContent(filename string, client pb.DatanodeClient, director pb.DirectorClient) error{
	req := &pb.GetFileContentRequest{
		Filename: filename,
	}

	resp, err := client.GetFileContent(context.Background(), req)
	if err != nil {
		return err
	}

	sendFileContentToDirector(resp.getContent(), filename, director)

	return nil
}

func sendFileContentToDirector(content string, filename string, director pb.DirectorClient) error {
	req := &pb.SendFileContentRequest{
		Filename: filename,
		Content: content,
	}

	_, err := client.SendFileContent(context.Background(), req)
	if err != nil {
		return err
	}

	return nil
}

func main() {
    conn, err := net.Listen("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
	defer conn.Close()
	
	client := pb.NewDatanodeClient(conn)
	connDirector, err := net.Dial("localhost:50050", grpc.WithInsecure())
	
	if err != nil {
		log.Fatalf("Failed to dial Director: %v", err)
	}
	defer connDirector.Close()

	//Cliente para el director
	directorClient := pb.NewDirectorClient(connDirector)
	//Implementar logica para la distributcion de los registros

}
