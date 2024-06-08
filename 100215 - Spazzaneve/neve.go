/*	In questa soluzione i camion rifornitori riempiono il silos anche se è già pieno
*	in quanto il testo dice "indipendentemente dalla quantità di sale precedentemente contenuta nel silos"
*/
package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Richiesta struct {
	id   int
	tipo int
	ack  chan int
}

const MAXBUFF = 100
const MAXPROC = 100

const SPAZZANEVE = 0
const SPARGISALE = 1
const CAMION = 2

const N = 5
const K = 2

var done = make(chan bool)
var termina = make(chan bool)
var terminaGestore = make(chan bool)

var entraDMN = make(chan Richiesta,MAXBUFF)
var prendiSale = make(chan Richiesta,MAXBUFF)
var rifornisciSilos = make(chan Richiesta,MAXBUFF)
var esciDMN = make(chan Richiesta)

func when(b bool, c chan Richiesta) chan Richiesta {
	if !b {
		return nil
	}
	return c
}
func getTipo(tipo int) string {
	switch tipo {
	case SPAZZANEVE:
		return "Spazzaneve"
	case SPARGISALE:
		return "Spargisale"
	case CAMION:
		return "Camion rifornitore"
	}
	return ""
}
func spazzaneve(myid int, tipo int) {
	var tt int
	tt = rand.Intn(5) + 1
	r := Richiesta{id: myid, ack: make(chan int)}
	fmt.Printf("[SPAZZANEVE %d] INIZIALIZZAZIONE \n", myid)
	time.Sleep(time.Duration(tt) * time.Second)
	entraDMN <- r
	<-r.ack
	fmt.Printf("[SPAZZANEVE %d] Entrato nel DMN \n", myid)
	time.Sleep(time.Duration(tt) * time.Second)
	esciDMN <- r
	fmt.Printf("[SPAZZANEVE %d] Uscito dal DMN \n", myid)
	done <- true
}

func spargisale(myid int, tipo int) {
	var tt int
	tt = rand.Intn(5) + 1
	r := Richiesta{id: myid, tipo: tipo, ack: make(chan int)}
	fmt.Printf("[SPARGISALE %d] INIZIALIZZAZIONE \n", myid)
	time.Sleep(time.Duration(tt) * time.Second)
	prendiSale <- r
	<-r.ack
	fmt.Printf("[SPARGISALE %d] Rifornimento di sale \n", myid)
	time.Sleep(time.Duration(tt) * time.Second)
	esciDMN <- r
	fmt.Printf("[SPAZZANEVE %d] Uscito dal DMN \n", myid)
	done <- true
}

func rifornitore(myid int, tipo int) {
	var tt int
	tt = rand.Intn(5) + 1
	r := Richiesta{id: myid, tipo: tipo, ack: make(chan int)}
	fmt.Printf("[RIFORNITORE %d] INIZIALIZZAZIONE \n", myid)
	for {

		time.Sleep(time.Duration(tt) * time.Second)
		rifornisciSilos <- r
		<-r.ack
		fmt.Printf("[RIFORNITORE %d] Rifornimento del silos  \n", myid)
		time.Sleep(time.Duration(tt) * time.Second)
		esciDMN <- r
		fmt.Printf("[RIFORNITORE %d] Uscito dal DMN \n", myid)

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
	var currK = K
	var currN = 0

	fmt.Printf("[GESTORE] INIZIALIZZAZIONE\n")
	for {
		select {
		case x := <-when(currN < N, entraDMN):
			currN++
			x.ack <- 1
			fmt.Printf("[GESTORE] Spazzaneve %d entrato, capacità DMN: %d \n", x.id, currN)
		case x := <-when(currN < N && currK > 0 && len(entraDMN) == 0, prendiSale):
			currK--
			currN++
			x.ack <- 1
			fmt.Printf("[GESTORE] Spargisale %d rifornito, capacità Silos: %d \n", x.id, currK)
		case x := <-when(currN < N && len(entraDMN) == 0 && len(prendiSale) == 0, rifornisciSilos):
			currK = K
			currN++
			x.ack <- 1
			fmt.Printf("[GESTORE] Camion rifornitore %d entrato, capacità Silos: %d \n", x.id, currK)
		case x := <-esciDMN:
			currN--
			fmt.Printf("[GESTORE] "+getTipo(x.tipo)+" uscito, capacità DMN: %d \n", currN)
		case <-terminaGestore:
			fmt.Printf("[GESTORE] Termino\n")
			done <- true
			return
		}

	}
}
func main() {
	rand.Seed(time.Now().Unix())
	nSpazzaneve := 4
	nSpargisale := 3
	nCamion := 2
	go gestore()

	for i := 0; i < nSpazzaneve; i++ {
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
	}

	terminaGestore <- true
	<-done
	fmt.Printf("\n HO FINITO ")
}
