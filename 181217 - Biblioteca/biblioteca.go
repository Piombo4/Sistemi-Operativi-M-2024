package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Richiesta struct {
	id  int
	ack chan int
}

const MAXSTUD = 10
const MAXBUFF = 100
const TRIENNALE = 0
const MAGISTRALE = 1
const NONLAUREANDO = 0
const LAUREANDO = 1
const N = 6

var done = make(chan bool)
var termina = make(chan bool)
var terminaGestore = make(chan bool)

var consegnaDocumento = make(chan Richiesta, MAXBUFF)
var entraBiblioteca [2][2]chan Richiesta
var esciBiblioteca = make(chan int)
var ritiraDocumento = make(chan Richiesta)

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

/*
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
*/
func studente(myid int, tipo int) {
	var tt int
	tt = rand.Intn(5) + 1
	r := Richiesta{id: myid, ack: make(chan int)}
	laureato := rand.Intn(2)
	fmt.Printf("[STUDENTE %d] INIZIALIZZAZIONE \n", myid)
	consegnaDocumento <- r
	<-r.ack
	fmt.Printf("[STUDENTE %d] Documento consegnato in portineria \n", myid)
	time.Sleep(time.Duration(tt) * time.Second)
	entraBiblioteca[tipo][laureato] <- r
	<-r.ack
	fmt.Printf("[STUDENTE %d] Entrato in biblioteca \n", myid)
	time.Sleep(time.Duration(tt) * time.Second)
	esciBiblioteca <- tipo
	fmt.Printf("[STUDENTE %d] Uscito dalla biblioteca \n", myid)
	time.Sleep(time.Duration(tt) * time.Second)
	ritiraDocumento <- r
	<-r.ack
	fmt.Printf("[STUDENTE %d] Termino \n", myid)
	
	done <- true
}

func portineria() {
	var documenti = make([]bool, MAXSTUD)
	fmt.Printf("[PORTINERIA] INIZIALIZZAZIONE \n")
	for {

		select {
		case x := <-consegnaDocumento:
			documenti[x.id] = true
			x.ack <- 1
			fmt.Printf("[PORTINERIA] Studente %d ha consegnato il suo documento\n", x.id)
		case x := <-ritiraDocumento:
			if documenti[x.id] {
				documenti[x.id] = false
				x.ack <- 1
				fmt.Printf("[PORTINERIA] Studente %d ha ritirato il suo documento\n", x.id)
			} else {
				x.ack <- -1
				fmt.Printf("[PORTINERIA] Studente %d non presente\n", x.id)
			}

		case <-termina:
			fmt.Printf("[PORTINERIA] Termino\n")
			done <- true
			return
		default:
			continue
		}
	}

}
func biblioteca() {
	var currCap [2]int

	fmt.Printf("[BIBLIOTECA] INIZIALIZZAZIONE\n")
	for {
		select {
		//ENTRATA MAGISTRALE
		case x := <-whenR(currCap[MAGISTRALE]+currCap[TRIENNALE] < N && (currCap[MAGISTRALE] < currCap[TRIENNALE] || (len(entraBiblioteca[TRIENNALE][NONLAUREANDO]) == 0 && len(entraBiblioteca[TRIENNALE][LAUREANDO]) == 0)), entraBiblioteca[MAGISTRALE][LAUREANDO]):
			currCap[MAGISTRALE]++
			x.ack <- 1
			fmt.Printf("[BIBLIOTECA] Entrato studente %d alla magistrale e laureando!, tot: %d\n", x.id, currCap[TRIENNALE]+currCap[MAGISTRALE])
		case x := <-whenR(currCap[MAGISTRALE]+currCap[TRIENNALE] < N && (currCap[MAGISTRALE] < currCap[TRIENNALE] || (len(entraBiblioteca[TRIENNALE][NONLAUREANDO]) == 0 && len(entraBiblioteca[TRIENNALE][LAUREANDO]) == 0)) && len(entraBiblioteca[MAGISTRALE][LAUREANDO]) == 0, entraBiblioteca[MAGISTRALE][NONLAUREANDO]):
			currCap[MAGISTRALE]++
			x.ack <- 1
			fmt.Printf("[BIBLIOTECA] Entrato studente %d alla magistrale e non laureando!, tot: %d\n", x.id, currCap[TRIENNALE]+currCap[MAGISTRALE])
		//ENTRATA TRIENNALE
		case x := <-whenR(currCap[MAGISTRALE]+currCap[TRIENNALE] < N && (currCap[MAGISTRALE] >= currCap[TRIENNALE] || (len(entraBiblioteca[MAGISTRALE][NONLAUREANDO]) == 0 && len(entraBiblioteca[MAGISTRALE][LAUREANDO]) == 0)), entraBiblioteca[TRIENNALE][LAUREANDO]):
			currCap[TRIENNALE]++
			x.ack <- 1
			fmt.Printf("[BIBLIOTECA] Entrato studente %d alla triennale e laureando!, tot: %d\n", x.id, currCap[TRIENNALE]+currCap[MAGISTRALE])
		case x := <-whenR(currCap[MAGISTRALE]+currCap[TRIENNALE] < N && (currCap[MAGISTRALE] >= currCap[TRIENNALE] || (len(entraBiblioteca[MAGISTRALE][NONLAUREANDO]) == 0 && len(entraBiblioteca[MAGISTRALE][LAUREANDO]) == 0)) && len(entraBiblioteca[TRIENNALE][LAUREANDO]) == 0, entraBiblioteca[TRIENNALE][NONLAUREANDO]):
			currCap[TRIENNALE]++
			x.ack <- 1
			fmt.Printf("[BIBLIOTECA] Entrato studente %d alla triennale e non laureando!, tot: %d\n", x.id, currCap[TRIENNALE]+currCap[MAGISTRALE])

		case x := <-esciBiblioteca:
			switch x {
			case MAGISTRALE:
				currCap[MAGISTRALE]--
				fmt.Printf("[BIBLIOTECA] Uscito studente della magistrale!, tot: %d\n", currCap[TRIENNALE]+currCap[MAGISTRALE])
			case TRIENNALE:
				currCap[TRIENNALE]--
				fmt.Printf("[BIBLIOTECA] Uscito studente della magistrale!, tot: %d\n", currCap[TRIENNALE]+currCap[MAGISTRALE])
			}

		case <-terminaGestore:
			fmt.Printf("[BIBLIOTECA] Termino\n")
			done <- true
			return
		}

	}
}
func main() {
	rand.Seed(time.Now().Unix())
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			entraBiblioteca[i][j] = make(chan Richiesta,MAXBUFF)
		}
	}
	go biblioteca()
	go portineria()

	for i := 0; i < MAXSTUD; i++ {
		r := rand.Intn(2)
		go studente(i, r)
	}

	for i := 0; i < MAXSTUD; i++ {
		<-done
	}

	termina <- true
	<-done

	terminaGestore <- true
	<-done
	fmt.Printf("\n HO FINITO ")
}
