package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"math/rand"
	"sync"
	"time"
	"strings"
	pb "https://github.com/dmoyan0/Lab4/blob/main/gRPC.proto"// Importa el paquete generado por protoc

	"google.golang.org/grpc"
)

type Mercenary struct {
	name string
	client pb.DirectorClient
	conn *grpc.Clientconn
	floors []string
	datanodeIP string
}

//Funcion que establece un nuevo mercenario con conexion dir
func NewMercenary(name string, dir string) (*Mercenary, error) {
	conn, err := grpc.Dial(dir, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := &grpc.NewDirectorClient(conn)
	
	//Elegir un datanode al azar para el mercenario
	datanodes := []string{"localhost:50052", "localhost:50053", "localhost:50054"}
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(datanodes))
	randomDatanode := datanodes[randomIndex]

	return &Mercenary{
		name: name,
		client: client,
		conn: conn,
		floors: []string{"Piso 1", "Piso 2", "Piso 3"},
		datanodeIP: randomDatanode,
	}, nil
}

//Ejecutar mercenario
func (m *Mercenary) Run() {
	defer m.conn.Close()

	//Informar al director del estado de preparacion
	req := &pb.MercenaryReadyRequest{
		Name: m.name,
		Ready: true,
	}
	resp, err := m.client.MercenaryReady(context.Background(), req)
	if err != nil {
		log.Fatalf("Error al confirmar estado de preparacion: %v", err)
	}
	fmt.Printf("Mercenario %s esta listo: %s\n", m.name, resp.Message)

	//Implementar el resto de la logica
}


func (m *Mercenary) Player() {
	defer m.conn.Close()

	var ready bool

	//Loop para esperar a que se confirme el estado de preparacion
	for !ready {
		var readyString string

		fmt.Print("¿Está listo para entrar al piso? [Si/No]: ")
		fmt.Scanf("%s", &readyString)

		readyString = strings.ToLower(strings.TrimSpace(readyString))

		if readyString == "si" {
			ready = true

			req := &pb.MercenaryReadyRequest{
				Name: m.name,
				Ready: ready,
			}

			resp, err := m.client.MercenaryReady(context.Background(), req)

			if err != nil {
				log.Fatalf("Failed to confirm: %v", err)
			}
		}
	}
	fmt.Printf("Mercenario %s esta listo: %s\n", m.name, resp.Message)
}

func main() {
	Director_dir := "localhost:50050"
	var wg sync.WaitGroup

	mercenaries := []string{"mercenary1", "mercenary2", "mercenary3", "mercenary4", "mercenary5", "mercenary6", "mercenary7"}

	var Player string
	fmt.Printf("Ingrese el nombre de su mercenario: ")
	fmt.Scan(&Player)

	for i:=0; i<7; i++ {
		wg.Add(1)
		go func(i int) {
			m, err := NewMercenary(mercenaries[i], Director_dir)
			if err != nil {
				log.Fatalf("error creating Mercenary: %v", err)
			}
			m.Run()
			wg.Done()
		}(i)
	}

	wg.Add(1)
	go func() {//Para el mercenario del jugador
		m, err := NewMercenary(Player, Director_dir)
		if err != nil {
			log.Fatalf("error creating Mercenary: %v", err)
		}
		m.Player()//Se envia a una funcion 
		wg.Done()
	}
	wg.Wait()

}