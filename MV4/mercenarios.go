package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	pb "github.com/dmoyan0/Lab4/grpc" // Importa el paquete generado por protoc

	"google.golang.org/grpc"
)

type Mercenary struct {
	name       string
	client     pb.DirectorClient
	conn       *grpc.Clientconn
	floors     int
	datanodeIP string
	arma       int
}

// Funcion que establece un nuevo mercenario con conexion dir
func NewMercenary(name string, dir string) (*Mercenary, error) {
	conn, err := grpc.Dial(dir, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := &pb.NewDirectorClient(conn)

	//Elegir un datanode al azar para el mercenario
	datanodes := []string{"localhost:50052", "localhost:50053", "localhost:50054"}
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(datanodes))
	randomDatanode := datanodes[randomIndex]

	return &Mercenary{
		name:       name,
		client:     client,
		conn:       conn,
		floors:     1, //Cada mercenario comienza en el piso 1
		datanodeIP: randomDatanode,
	}, nil
}

// Ejecutar mercenario
func (m *Mercenary) Run(int floor) {
	defer m.conn.Close()

	//Informar al director del estado de preparacion
	req := &pb.MercenaryReadyRequest{
		Name:  m.name,
		Ready: true,
	}
	resp, err := m.client.MercenaryReady(context.Background(), req)
	if err != nil {
		log.Fatalf("Error al confirmar estado de preparacion: %v", err)
	}
	fmt.Printf("Mercenario %s esta listo: %s\n", m.name, resp.Message)

	//Implementar el resto de la logica
	if floor <= 3 {
		m.Decision(floor)
		m.Run(floor++)
	}	
}

func (m *Mercenary) Player_Ready(int floor) {
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
				Name:   m.name,
				Ready:  ready,
				Floors: floor,
			}

			resp, err := m.client.MercenaryReady(context.Background(), req)

			if err != nil {
				log.Fatalf("Failed to confirm: %v", err)
			}
		}
	}
	fmt.Printf("Mercenario %s esta listo: %s\n", m.name, resp.Message)

	if floor <= 3 {
		m.PlayerDecision(floor)
		m.Player_Ready(floor++, done)
	}

}

func (m *Mercenary) Decision(floor int) {
	switch floor {
	case 1:
		armas := []int{1, 2, 3}
		rand.Seed(time.Now().UnixNano())
		randomIndex := rand.Intn(len(armas))
		arma := armas[randomIndex]

		req := &pb.MercenaryDecisionRequest{
			Name:       m.name,
			Floor:      int(floor),
			ArmaChoice: int(arma),
		}

		resp, err := m.client.MercenaryDecision(context.Background(), req)
		if err != nil {
			log.Fatalf("Error al tomar decisión en el piso %d: %v", floor, err)
		}
		fmt.Printf("Decisión del mercenario %s en el piso %d: %s\n", m.name, floor, resp.Message)
		if resp.estado == false {//Muere el mercenario
			wg.Done();
			log.Printf("Murio el mercenario: %s", m.name)
		}
		else {
			fmt.Printf("Mercenario %s sobrevivio el piso", m.name)
		}

	case 2:
		pasillo := []string{"A", "B"}
		rand.Seed(time.Now().UnixNano())
		randomIndex := rand.Intn(len(pasillo))
		decision := pasillo[randomIndex]

		req := &pb.MercenaryDecisionRequest{
			Name:     m.name,
			Floor:    int(floor),
			Decision: decision,
		}

		resp, err := m.client.MercenaryDecision(context.Background(), req)
		if err != nil {
			log.Fatalf("Error al tomar decisión en el piso %d: %v", floor, err)
		}
		fmt.Printf("Decisión del mercenario %s en el piso %d: %s\n", m.name, floor, resp.Message)
		if resp.estado == false {//Muere el mercenario
			wg.Done();
			log.Printf("Murio el mercenario: %s", m.name)
		}
		else {
			fmt.Printf("Mercenario %s sobrevivio el piso", m.name)
		}


	case 3:
		aciertosMercenario := 0

		//5 rondas cada uno de los mercenarios que quedan
		for i := 0; i < 5; i++ {
			//Elige numero
			mercenarioNumero := rand.Intn(15) + 1
			req := &pb.MercenaryDecisionRequest{
				Name:     m.name,
				Floor:    int(floor),
				Decision: mercenarioNumero,
			}
			resp, err := m.client.MercenaryDecision(context.Background(), req)
			if err != nil {
				log.Fatalf("Error al tomar decisión en el piso %d: %v", floor, err)
			}
			if resp.estado == false {//Muere el mercenario
				wg.Done();
				log.Printf("Murio el mercenario: %s", m.name)
			}
			else {
				fmt.Printf("Mercenario %s sobrevivio el piso", m.name)
			}

			if resp.Acierto {
				aciertosMercenario++
			}
			fmt.Printf("Ronda %d: Decisión del mercenario %s en el piso %d: %s , aciertos: %d\n", i+1, m.name, floor, resp.Message, aciertosMercenario)
		}
		fmt.Printf("Total de aciertos del mercenario %s en el piso %d: %d\n", m.name, floor, aciertosMercenario)

	default:
		log.Printf("Piso inválido: %d", floor)
	}
}

func (m *Mercenary) PlayerDecision(floor int) {
	fmt.Printf("Solicitar monto?(si)")
	var choice string
	fmt.Scanf("%s", &choice)
	if choice == "si"{
		GetAcumulatedAmount()
	}

	switch floor {
	case 1:
		var arma int
		fmt.Print("Seleccione un arma (1: Escopeta, 2: Rifle, 3: Puños electricos)")
		fmt.Scanf("%d", &arma)

		req := &pb.MercenaryDecisionRequest{
			Name:       m.name,
			Floor:      int(floor),
			ArmaChoice: arma,
		}

		resp, err := m.client.MercenaryDecision(context.Background(), req)
		if err != nil {
			log.Fatalf("Error al tomar decisión en el piso %d: %v", floor, err)
		}
		fmt.Printf("Decisión del mercenario %s en el piso %d: %s\n", m.name, floor, resp.Message)
		if resp.estado == false {//Muere el mercenario
			wg.Done();
			log.Printf("Murio el mercenario: %s", m.name)
		}
		else {
			fmt.Printf("Mercenario %s sobrevivio el piso", m.name)
		}

	case 2:
		var decision string
		fmt.Print("Seleccione un pasillo (A o B): ")
		fmt.Scanf("%s, &decision")

		req := &pb.MercenaryDecisionRequest{
			Name:     m.name,
			Floor:    int(floor),
			Decision: decision,
		}

		resp, err := m.client.MercenaryDecision(context.Background(), req)
		if err != nil {
			log.Fatalf("Error al tomar decisión en el piso %d: %v", floor, err)
		}
		fmt.Printf("Decisión del mercenario %s en el piso %d: %s\n", m.name, floor, resp.Message)
		if resp.estado == false {//Muere el mercenario
			wg.Done();
			log.Printf("Murio el mercenario: %s", m.name)
		}
		else {
			fmt.Printf("Mercenario %s sobrevivio el piso", m.name)
		}

	case 3:
		aciertosMercenario := 0

		for i := 0; i < 5; i++ {
			var mercenarioNumero int
			fmt.Printf("Ronda %d: Elija un número entre 1 y 15: ", i+1)
			fmt.Scanf("%d", &mercenarioNumero)

			req := &pb.MercenaryDecisionRequest{
				Name:     m.name,
				Floor:    int(floor),
				Decision: int(mercenarioNumero),
			}
			resp, err := m.client.MercenaryDecision(context.Background(), req)
			if err != nil {
				log.Fatalf("Error al tomar decisión en el piso %d: %v", floor, err)
			}

			if resp.Acierto {
				aciertosMercenario++
			}
			fmt.Printf("Ronda %d: Decisión del mercenario %s en el piso %d: %s , aciertos: %d\n", i+1, m.name, floor, resp.Message, aciertosMercenario)
			if resp.estado == false {//Muere el mercenario
				wg.Done();
				log.Printf("Murio el mercenario: %s", m.name)
			}
			else {
				fmt.Printf("Mercenario %s sobrevivio el piso", m.name)
			}
		}
		fmt.Printf("Total de aciertos del mercenario %s en el piso %d: %d\n", m.name, floor, aciertosMercenario)

	default:
		log.Printf("Piso inválido: %d", floor)
	}
}

func (m *Mercenary) GetAcumulatedAmount() {
	req := &pb.GetAcumulatedAmount{}

	resp, err := m.client.GetAcumulatedAmount(context.Background(), req)
	if err != nil {
		log.Fatalf("Error al obtener el monto acumulado: %v", err)
	}

	fmt.Printf("El monto acumulado es: %f", resp.Monto)
}

func main() {
	Director_dir := "localhost:50050"
	var wg sync.WaitGroup

	mercenaries := []string{"mercenary1", "mercenary2", "mercenary3", "mercenary4", "mercenary5", "mercenary6", "mercenary7"}

	var Player string
	fmt.Printf("Ingrese el nombre de su mercenario: ")
	fmt.Scan(&Player)

	//Creacion de los 7 npcs
	for i := 0; i < 7; i++ {
		wg.Add(1)
		go func(i int) {
			m, err := NewMercenary(mercenaries[i], Director_dir)
			if err != nil {
				log.Fatalf("error creating Mercenary: %v", err)
			}
			m.Run(1)
			wg.Done()
		}(i)
	}
	//Creacion del player
	wg.Add(1)
	go func() { //Para el mercenario del jugador
		defer wg.Done()
		m, err := NewMercenary(Player, Director_dir)
		if err != nil {
			log.Fatalf("error creating Mercenary: %v", err)
		}
		m.Player_Ready(1) //Se envia a una funcion
	}()
	wg.Wait()

}
