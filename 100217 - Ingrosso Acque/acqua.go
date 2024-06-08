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

const MAXBUFF = 100

const Y = 2
const Z = 1
const X = 2
const K = 8

const MAX_P = 10
const MAX_V = 10

const CONTANTI = 0
const BANCOMAT = 1
const BONIFICO = 1

const PA = 3
const PV = 2

var done = make(chan bool)
var termina = make(chan bool)
var terminaGestore = make(chan bool)
var acquistaAcqua [2]chan Richiesta
var vendiAcqua [2]chan Richiesta

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

func acquirente(myid int, metodo int) {
	var tt int
	tt = rand.Intn(5) + 1
	r := Richiesta{id: myid, ack: make(chan int)}
	fmt.Printf("[ACQUIRENTE %d] INIZIALIZZAZIONE \n", myid)
	time.Sleep(time.Duration(tt) * time.Second)
	acquistaAcqua[metodo] <- r
	<-r.ack
	fmt.Printf("[ACQUIRENTE %d] Ho comprato %d acque tramite %d e consegnato %d vuoti \n", myid, Y, metodo, Z)
	done <- true
}

func fornitore(myid int, metodo int) {
	var tt int
	tt = rand.Intn(5) + 1
	r := Richiesta{id: myid, ack: make(chan int)}
	fmt.Printf("[FORNITORE %d] INIZIALIZZAZIONE \n", myid)
	for {

		time.Sleep(time.Duration(tt) * time.Second)
		vendiAcqua[metodo] <- r
		<-r.ack
		fmt.Printf("[FORNITORE %d] Ho venduto %d acque tramite %d e ritirato tutti i vuoti \n", myid, X, metodo)
		select {
		case <-termina:
			fmt.Printf("[FORNITORE %d] Termino\n", myid)
			done <- true
			return
		default:
			continue
		}
	}

}
func ditta() {

	var cassa int = X * PV
	var contoCorrente = X * PV
	var numPiene int = MAX_P
	var numVuote int = 0

	fmt.Printf("[DITTA] INIZIALIZZAZIONE\n")
	for {
		select {
		case x := <-whenR(numPiene >= Y && numVuote+Z <= MAX_V && (cassa < K || len(acquistaAcqua[BANCOMAT]) == 0), acquistaAcqua[CONTANTI]):
			cassa += PA * Y
			numVuote += Z
			numPiene -= Y
			x.ack <- 1
			fmt.Printf("[DITTA] Acquirente %d ha acquistato tramite contanti, tot cassa: %d, capacità piene: %d, capacità vuote: %d\n", x.id, cassa, numPiene, numVuote)

		case x := <-whenR(numPiene >= Y && numVuote+Z <= MAX_V && (cassa >= K || len(acquistaAcqua[CONTANTI]) == 0), acquistaAcqua[BANCOMAT]):
			contoCorrente += PA * Y
			numVuote += Z
			numPiene -= Y
			x.ack <- 1
			fmt.Printf("[DITTA] Acquirente %d ha acquistato tramite bancomat, capacità piene: %d, capacità vuote: %d\n", x.id, numPiene, numVuote)
		case x := <-whenR(numPiene+X <= MAX_P && contoCorrente >= PV*X && (cassa < K || len(vendiAcqua[CONTANTI]) == 0), vendiAcqua[BONIFICO]):
			contoCorrente -= PV * X
			numVuote = 0
			numPiene += X
			x.ack <- 1
			fmt.Printf("[DITTA] Fornitore %d ha venduto tramite bonifico, capacità piene: %d, capacità vuote: %d\n", x.id, numPiene, numVuote)

		case x := <-whenR(cassa >= PV*X && numPiene+X <= MAX_P && (cassa >= K || len(vendiAcqua[BONIFICO]) == 0), vendiAcqua[CONTANTI]):
			cassa -= PV * X
			numVuote = 0
			numPiene += X
			x.ack <- 1
			fmt.Printf("[DITTA] Fornitore %d ha venduto tramite contanti, tot cassa: %d, capacità piene: %d, capacità vuote: %d\n", x.id, cassa, numPiene, numVuote)

		case <-terminaGestore:
			fmt.Printf("[DITTA] Termino\n")
			done <- true
			return
		}

	}
}
func main() {
	rand.Seed(time.Now().Unix())
	nAcquirenti := 5
	nFornitori := 2
	go ditta()
	for i := 0; i < 2; i++ {
		acquistaAcqua[i] = make(chan Richiesta, MAXBUFF)
		vendiAcqua[i] = make(chan Richiesta, MAXBUFF)
	}

	for i := 0; i < nAcquirenti; i++ {
		var r = rand.Intn(2)
		go acquirente(i, r)
	}
	for i := 0; i < nFornitori; i++ {
		var r = rand.Intn(2)
		go fornitore(i, r)
	}
	for i := 0; i < nAcquirenti; i++ {
		<-done
	}
	for i := 0; i < nFornitori; i++ {
		termina <- true
		<-done
	}

	terminaGestore <- true
	<-done
	fmt.Printf("\n HO FINITO ")
}
