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
	datanodeslist := []string{"localhost:50052", "localhost:50053", "localhost:50054"}//Modificar de acuerdo a nombres e IPs

	rand.Seed(time.Now().UnixNano())//Se elije una direccion al azar
	chosenDatanode := datanodeslist[rand.Intn(len(datanodeslist))]

	//Establece una conexion gRPC con el Datanode elegido
	conn, err := net.Dial(chosenDatanode, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error al conectarse: %v", err)
	}
	defer conn.Close()

	//Crea un cliente gRPC para el servicio Datanode
	client := pb.NewClient(conn)

	_,err = client.RegisterDecision(ctx, &pb.RegisterDecisionRequest{
		MercenaryName: req.MercenaryName,
		Floor: req.Floor,
		DatanodeIP: req.DatanodeIP,
	})
	if err != nil {
		log.Fatalf("Error al enviar el registro al Datanode": %v, err)
	}
    // Retorna una respuesta indicando que el registro se procesó correctamente
    return &pb.RegisterDecisionResponse{Message: "Decision registered successfully"}, nil
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
    s := grpc.NewServer()
    pb.RegisterNamenodeServer(s, &server{})
    if err := s.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
