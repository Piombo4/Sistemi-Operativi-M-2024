
package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Richiesta struct {
	id   int
	ack  chan int
}

const MAXBUFF = 100


var done = make(chan bool)
var termina = make(chan bool)
var terminaGestore = make(chan bool)

func sleepRandTime(timeLimit int) {
	if timeLimit > 0 {
		time.Sleep(time.Duration(rand.Intn(timeLimit)+1) * time.Second)
	}
}

func whenR(b bool, c chan Richiesta) chan Richiesta {
	if !b {
		return nil
	}
	return c
}
func when(b bool, c chan int) chan int {
	if !b {
		return nil
	}
	return c
}

/*func getTipo(tipo int) string {
	switch tipo {
	case SPAZZANEVE:
		return "Spazzaneve"
	case SPARGISALE:
		return "Spargisale"
	case CAMION:
		return "Camion rifornitore"
	}
	return ""
}*/
func client(myid int) {
	var tt int
	tt = rand.Intn(5) + 1
	r := Richiesta{id: myid, ack: make(chan int)}
	fmt.Printf("[CLIENT %d] INIZIALIZZAZIONE \n", myid)
	time.Sleep(time.Duration(tt) * time.Second)
	done <- true
}


func persistentClient(myid int) {
	var tt int
	tt = rand.Intn(5) + 1
	r := Richiesta{id: myid, ack: make(chan int)}
	fmt.Printf("[persistentClient %d] INIZIALIZZAZIONE \n", myid)
	for {

		time.Sleep(time.Duration(tt) * time.Second)
		
		select {
		case <-termina:
			fmt.Printf("[RIFORNITORE %d] Termino\n", myid)
			done <- true
			return
		default:
			continue
		}
	}

}
func gestore() {
	

	fmt.Printf("[GESTORE] INIZIALIZZAZIONE\n")
	for {
		select {
		
		case <-terminaGestore:
			fmt.Printf("[GESTORE] Termino\n")
			done <- true
			return
		}

	}
}
func main() {
	rand.Seed(time.Now().Unix())
	
	go gestore()

	/*for i := 0; i < nSpazzaneve; i++ {
		go spazzaneve(i, SPAZZANEVE)
	}
	for i := 0; i < nSpargisale; i++ {
		go spargisale(i, SPARGISALE)
	}
	for i := 0; i < nCamion; i++ {
		go rifornitore(i, CAMION)
	}
	for i := 0; i < nSpazzaneve+nSpargisale; i++ {
		<-done
	}
	for i := 0; i < nCamion; i++ {
		termina <- true
		<-done
	}*/

	terminaGestore <- true
	<-done
	fmt.Printf("\n HO FINITO ")
}
