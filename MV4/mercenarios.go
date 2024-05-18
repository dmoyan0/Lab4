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
		name:   name,
		client: client,
		conn:   conn,
		//floors:     []string{"Piso 1", "Piso 2", "Piso 3"},
		datanodeIP: randomDatanode,
	}, nil
}

// Ejecutar mercenario
func (m *Mercenary) Run() {
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
	for floor := 1; floor <= 3; floor++ {
		m.Decision(floor)
	}
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
				Name:  m.name,
				Ready: ready,
			}

			resp, err := m.client.MercenaryReady(context.Background(), req)

			if err != nil {
				log.Fatalf("Failed to confirm: %v", err)
			}
		}
	}
	fmt.Printf("Mercenario %s esta listo: %s\n", m.name, resp.Message)

	for floor := 1; floor <= 3; floor++ {
		m.PlayerDecision(floor)
	}
	//Meter lo anterior a una funcion para llamarla posteriormente
	//Una vez confirmada la preparacion se debe interactuar con los niveles del Director
	/*for ready {//cambiar ready
		switch req.Floor{
		case 1:
			armas := []int {1, 2, 3}
			rand.Seed(time.Now().UnixNano())
			randomIndex := rand.Intn(len(datanodes))
			arma := armas[randomIndex]
			//Cambiar por interfaz
		}
	}*/

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
			m.Run()
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
		m.Player() //Se envia a una funcion
	}()
	wg.Wait()

}
