package main

import (
	"fmt"
	"math/rand"
	"time"
)

const MAXVISM = 20
const MAXVISS = 20
const MAXVISSM = 20
const MAXOPM = 10
const MAXOPS = 10

const MAXBUFF = 100
const MAXPROC = 100
const NM = 5
const NS = 6

const MUSEO = 0 //OPERATORI MUSEO
const SALA = 1  //OPERATORI SALA MOSTRE

const OPM = 0 //OPERATORI MUSEO
const OPS = 1 //OPERATORI SALA MOSTRE

const VM = 0  //VISITATORI MUSEO
const VS = 1  //VISITATORI SALA MOSTRE
const VMS = 2 //VISITATORI MUSEO E SALA MOSTRE

type Richiesta struct {
	id  int
	ack chan int
}

var done = make(chan bool)
var termina = make(chan bool)
var terminaGestore = make(chan bool)

var entraMuseo = make(chan Richiesta, MAXBUFF)
var esceMuseo = make(chan Richiesta, MAXBUFF)
var entraSala = make(chan Richiesta, MAXBUFF)
var esceSala = make(chan Richiesta, MAXBUFF)
var entraSM = make(chan Richiesta, MAXBUFF)
var esceSM = make(chan Richiesta, MAXBUFF)

var entraMuseo_op = make(chan Richiesta, MAXBUFF)
var esceMuseo_op = make(chan Richiesta, MAXBUFF)
var entraSala_op = make(chan Richiesta, MAXBUFF)
var esceSala_op = make(chan Richiesta, MAXBUFF)

func when(b bool, c chan Richiesta) chan Richiesta {
	if !b {
		return nil
	}
	return c
}

func visitatore(myid int, tipo int) {
	var tt int
	tt = rand.Intn(5) + 1
	fmt.Printf("Inizializzazione visitatore %d - tipo %d \n", myid, tipo)
	r := Richiesta{id: myid, ack: make(chan int)}
	time.Sleep(time.Duration(tt) * time.Second)
	switch tipo {
	case VM:
		fmt.Printf("Visitatore %d entra nel museo\n", myid)
		entraMuseo <- r
		<-r.ack
		fmt.Printf("Visitatore %d gira per il museo\n", myid)
		time.Sleep(time.Duration(tt) * time.Second)
		fmt.Printf("Visitatore %d esce dal museo\n", myid)
		esceMuseo <- r
	case VS:
		fmt.Printf("Visitatore %d entra nella sala mostre\n", myid)
		entraSala <- r
		<-r.ack
		fmt.Printf("Visitatore %d gira per la sala mostre\n", myid)
		time.Sleep(time.Duration(tt) * time.Second)
		fmt.Printf("Visitatore %d esce dalla sala mostre\n", myid)
		esceSala <- r
	case VMS:
		fmt.Printf("Visitatore %d accede ad entrambe le aree\n", myid)
		entraSM <- r
		<-r.ack
		fmt.Printf("Visitatore %d gira per entrambe le aree\n", myid)
		time.Sleep(time.Duration(tt) * time.Second)
		fmt.Printf("Visitatore %d esce da entrambe le aree\n", myid)
		esceSM <- r
	}
	done <- true

}

func operatore(myid int, tipo int) {
	var tt int
	tt = rand.Intn(5) + 1
	fmt.Printf("Inizializzazione operatore %d - tipo %d \n", myid, tipo)
	r := Richiesta{id: myid, ack: make(chan int)}
	time.Sleep(time.Duration(tt) * time.Second)
	for {
		switch tipo {
		case OPM:
			fmt.Printf("Operatore %d prova ad entrare nel museo\n", myid)
			entraMuseo_op <- r
			<-r.ack
			fmt.Printf("Operatore %d entrato nel museo\n", myid)
			time.Sleep(time.Duration(tt) * time.Second)
			fmt.Printf("Operatore %d vuole andare in pausa\n", myid)
			esceMuseo_op <- r
			<-r.ack
			fmt.Printf("Operatore %d in pausa\n", myid)
			time.Sleep(time.Duration(tt) * time.Second)
		case OPS:
			fmt.Printf("Operatore %d prova ad entrare nelle sale\n", myid)
			entraSala_op <- r
			<-r.ack
			fmt.Printf("Operatore %d entrato nelle sale\n", myid)
			time.Sleep(time.Duration(tt) * time.Second)
			fmt.Printf("Operatore %d vuole andare in pausa\n", myid)
			esceSala_op <- r
			<-r.ack
			fmt.Printf("Operatore %d in pausa\n", myid)
			time.Sleep(time.Duration(tt) * time.Second)
		}
		select {
		case <-termina:
			fmt.Printf("[OPERATORE %d] Termino\n", myid)
			done <- true
			return
		default:
			continue
		}

	}
}
func gestore() {
	var nVis [2]int //pedoni e auto in dir nord
	var nOp [2]int  //pedoni e auto in dir nord

	for {
		select {
		case x := <-when(nVis[MUSEO] < NM && nOp[MUSEO] > 0, entraMuseo):
			nVis[MUSEO]++
			fmt.Printf("[GESTORE] Autorizzo l'ingresso del visitatore %d al museo \n", x.id)
			x.ack <- 1
		case x := <-when(nVis[SALA] < NS && len(entraMuseo) == 0 && nOp[SALA] > 0, entraSala):
			nVis[SALA]++
			fmt.Printf("[GESTORE] Autorizzo l'ingresso del visitatore %d alla sala \n", x.id)
			x.ack <- 1
		case x := <-when(nVis[SALA] < NS && nVis[MUSEO] < NM && len(entraMuseo) == 0 && len(entraSala) == 0 && nOp[SALA] > 0 && nOp[MUSEO] > 0, entraSM):
			nVis[SALA]++
			nVis[MUSEO]++
			fmt.Printf("[GESTORE] Autorizzo l'ingresso del visitatore %d ad entrambe le aree \n", x.id)
			x.ack <- 1
		case x := <-esceMuseo:
			nVis[MUSEO]--
			fmt.Printf("[GESTORE] Visitatore %d esce dal museo \n", x.id)
		case x := <-esceSala:
			nVis[SALA]--
			fmt.Printf("[GESTORE] Visitatore %d esce dalla sala \n", x.id)
		case x := <-esceSM:
			nVis[SALA]--
			nVis[MUSEO]--
			fmt.Printf("[GESTORE] Visitatore %d esce da entrambe le aree \n", x.id)
		//OPERATORI
		case x := <-entraMuseo_op:
			nOp[MUSEO]++
			fmt.Printf("[GESTORE] Operatore %d entra nel museo \n", x.id)
			x.ack <- 1
		case x := <-entraSala_op:
			nOp[SALA]++
			fmt.Printf("[GESTORE] Operatore %d entra nella sala \n", x.id)
			x.ack <- 1
		case x := <-when(nVis[MUSEO] == 0 && len(entraMuseo) == 0 && len(entraSM) == 0 || nOp[MUSEO] > 1, esceMuseo_op):
			nOp[MUSEO]--
			fmt.Printf("[GESTORE] Operatore %d esce dal museo \n", x.id)
			x.ack <- 1
		case x := <-when(nVis[SALA] == 0 && len(entraSala) == 0 && len(entraSM) == 0 || nOp[SALA] > 1, esceSala_op):
			nOp[SALA]--
			fmt.Printf("[GESTORE] Operatore %d esce dalle sale \n", x.id)
			x.ack <- 1
		case <-terminaGestore:
			fmt.Printf("[GESTORE] Termino\n")
			done <- true
			return
		}
	}
}
func main() {

	nVisitatori := 6
	nOperatori := 3

	rand.Seed(time.Now().Unix())

	go gestore()

	for i := 0; i < nVisitatori; i++ {
		r := rand.Intn(3)
		go visitatore(i, r)
	}

	for i := 0; i < nOperatori; i++ {
		r := rand.Intn(2)
		go operatore(i, r)
	}

	for i := 0; i < nVisitatori; i++ {
		<-done
	}

	for i := 0; i < nOperatori; i++ {
		termina <- true
		<-done
	}

	terminaGestore <- true
	<-done
	fmt.Printf("\n HO FINITO ")

}
