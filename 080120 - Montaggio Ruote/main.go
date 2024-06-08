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
const A = 0
const B = 1
const MAXP = 10
const MAXC = 10

const TOT = 2

var currAuto = 0
var tipoModello = [2]string{"Modello A", " Modello B"}
var tipoNastro = [2]string{"Pneumatici", "Cerchi"}
var done = make(chan bool)
var termina = make(chan bool)
var terminaGestore = make(chan bool)

var depositaP [2]chan Richiesta
var depositaC [2]chan Richiesta

var prelevaP [2]chan Richiesta
var prelevaC [2]chan Richiesta

var monta [2]chan int

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

func nastroC(tipo int) {

	r := Richiesta{id: tipo, ack: make(chan int)}
	fmt.Printf("[Nastro Cerchi %s] INIZIALIZZAZIONE \n", tipoModello[tipo])
	for {
		time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)
		depositaC[tipo] <- r
		ris := <-r.ack
		if ris == -1 {
			fmt.Printf("[Nastro Pneumatici %s] Termino! \n", tipoModello[tipo])
			done <- true
			return
		} else {
			//fmt.Printf("[Nastro Cerchi %s] Ho depositato un cerchio \n", tipoModello[tipo])
		}
	}
}
func nastroP(tipo int) {

	r := Richiesta{id: tipo, ack: make(chan int)}
	fmt.Printf("[Nastro Pneumatici %s] INIZIALIZZAZIONE \n", tipoModello[tipo])
	for {
		time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)
		depositaP[tipo] <- r
		ris := <-r.ack
		if ris == -1 {
			fmt.Printf("[Nastro Pneumatici %s] Termino! \n", tipoModello[tipo])
			done <- true
			return
		} else {
			//fmt.Printf("[Nastro Pneumatici %s] Ho depositato uno pneumatico \n", tipoModello[tipo])
		}
	}
}
func robot(tipo int) {
	var tt int
	tt = rand.Intn(5) + 1
	r := Richiesta{id: tipo, ack: make(chan int)}
	fmt.Printf("[Robot %s] INIZIALIZZAZIONE \n", tipoModello[tipo])
	var counter = 0
	for {
		if currAuto == TOT {
			fmt.Printf("[Robot %s] Ho finito! \n", tipoModello[tipo])
			done <- true
			return
		}
		time.Sleep(time.Duration(tt) * time.Second)
		prelevaC[tipo] <- r
		<-r.ack
		//fmt.Printf("[%s] Ho prelevato un cerchio... \n", tipoModello[tipo])
		time.Sleep(time.Duration(tt) * time.Second)
		prelevaP[tipo] <- r
		<-r.ack
		//fmt.Printf("[%s] Ho prelevato uno pneumatico... \n", tipoModello[tipo])
		//fmt.Printf("[%s] Ruota montata! \n", tipoModello[tipo])
		counter++
		if counter % 4 == 0 {
			monta[tipo] <- tipo
			currAuto++
			counter=0
			fmt.Printf("[Robot %s] Ho montato un auto! TOT: %d\n", tipoModello[tipo], currAuto)
		}

	}

}
func debug(nCerchi int, nPneumatici int) {
	fmt.Printf("[DEBUG] nCerchi: %d - nPneumatici: %d\n", nCerchi, nPneumatici)
}
func deposito() {
	var nCerchi = 0
	var nPneumatici = 0
	var nMontaggi = [2]int{0, 0}
	var fine = false
	fmt.Printf("[Deposito] INIZIALIZZAZIONE\n")
	for {
		select {
		//DEPOSITO
		case x := <-whenR(!fine && nCerchi < MAXC && (nMontaggi[A] < nMontaggi[B] || (len(depositaC[B]) == 0 && len(depositaP[B]) == 0)), depositaC[A]):
			nCerchi++
			fmt.Printf("[Deposito] Nastro dei cerchi ha depositato un cerchio A\n")
			debug(nCerchi, nPneumatici)
			x.ack <- 1
		case x := <-whenR(!fine && nCerchi < MAXC && (nMontaggi[B] <= nMontaggi[A] || (len(depositaC[A]) == 0 && len(depositaP[A]) == 0)), depositaC[B]):
			nCerchi++
			fmt.Printf("[Deposito] Nastro dei cerchi ha depositato un cerchio B\n")

			debug(nCerchi, nPneumatici)
			x.ack <- 1
		case x := <-whenR(!fine && nPneumatici < MAXP && (nMontaggi[A] < nMontaggi[B] || (len(depositaC[B]) == 0 && len(depositaP[B]) == 0)), depositaP[A]):
			nPneumatici++
			fmt.Printf("[Deposito] Nastro dei cerchi ha depositato uno pneumatico A\n")

			debug(nCerchi, nPneumatici)
			x.ack <- 1
		case x := <-whenR(!fine && nPneumatici < MAXP && (nMontaggi[B] <= nMontaggi[A] || (len(depositaC[A]) == 0 && len(depositaP[A]) == 0)), depositaP[B]):
			nPneumatici++
			fmt.Printf("[Deposito] Nastro dei cerchi ha depositato uno pneumatico B\n")

			debug(nCerchi, nPneumatici)
			x.ack <- 1
			//PRELIEVO
		case x := <-whenR(nPneumatici > 0 && (nMontaggi[B] <= nMontaggi[A] || (len(prelevaC[A]) == 0 && len(prelevaP[A]) == 0)), prelevaP[B]):
			nPneumatici--
			fmt.Printf("[Deposito] Robot B ha prelevato uno pneumatico\n")
			debug(nCerchi, nPneumatici)
			x.ack <- 1
		case x := <-whenR(nCerchi > 0 && (nMontaggi[B] <= nMontaggi[A] || (len(prelevaC[A]) == 0 && len(prelevaP[A]) == 0)), prelevaC[B]):
			nCerchi--
			fmt.Printf("[Deposito] Robot B ha prelevato un cerchio\n")
			debug(nCerchi, nPneumatici)
			x.ack <- 1
		case x := <-whenR(nPneumatici > 0 && (nMontaggi[A] < nMontaggi[B] || (len(prelevaC[B]) == 0 && len(prelevaP[B]) == 0)), prelevaP[A]):
			nPneumatici--
			fmt.Printf("[Deposito] Robot A ha prelevato uno pneumatico\n")
			debug(nCerchi, nPneumatici)
			x.ack <- 1
		case x := <-whenR(nCerchi > 0 && (nMontaggi[A] < nMontaggi[B] || (len(prelevaC[B]) == 0 && len(prelevaP[B]) == 0)), prelevaC[A]):
			nCerchi--
			fmt.Printf("[Deposito] Robot A ha prelevato un cerchio\n")
			debug(nCerchi, nPneumatici)
			x.ack <- 1

		case x := <-monta[A]:
			nMontaggi[A]++
			fmt.Printf("[Deposito] Robot %s ha montato un'auto!, TOT: %d\n", tipoModello[x], nMontaggi[A])
		case x := <-monta[B]:
			nMontaggi[B]++
			fmt.Printf("[Deposito] Robot %s ha montato un'auto!, TOT: %d \n", tipoModello[x], nMontaggi[B])

		case x := <-whenR(fine, depositaC[A]):
			x.ack <- -1
			fmt.Printf("[Deposito] Termino nastro cerchi A\n")
		case x := <-whenR(fine, depositaC[B]):
			x.ack <- -1
			fmt.Printf("[Deposito] Termino nastro cerchi B\n")
		case x := <-whenR(fine, depositaP[A]):
			x.ack <- -1
			fmt.Printf("[Deposito] Termino nastro pneumatici A\n")
		case x := <-whenR(fine, depositaP[B]):
			x.ack <- -1
			fmt.Printf("[Deposito] Termino nastro pneumatici B\n")
		case <-termina:
			fine = true
		case <-terminaGestore:
			fmt.Printf("[Deposito] Termino\n")
			done <- true
			return
		}

	}
}
func main() {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 2; i++ {
		depositaC[i] = make(chan Richiesta, MAXBUFF)
		depositaP[i] = make(chan Richiesta, MAXBUFF)
		prelevaC[i] = make(chan Richiesta, MAXBUFF)
		prelevaP[i] = make(chan Richiesta, MAXBUFF)
		monta[i] = make(chan int, MAXBUFF)
	}

	go deposito()
	go nastroC(A)
	go nastroC(B)
	go nastroP(A)
	go nastroP(B)
	go robot(A)
	go robot(B)
	for i := 0; i < 2; i++ {
		<-done
	}
	termina <- true
	for i := 0; i < 4; i++ {
		<-done
	}
	terminaGestore <- true
	<-done
	fmt.Printf("\n HO FINITO ")
}
