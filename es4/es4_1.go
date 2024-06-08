package main

import (
	"fmt"
	"math/rand"
	"time"
)

const MAXPROC = 100
const MAXBUFF = 100
const MAXCAP = 10
const N int = 0
const S int = 1

var done = make(chan bool)
var termina = make(chan bool)
var entrata_auto_N = make(chan int, MAXBUFF)   // necessità di accodamento per priorità
var entrata_auto_S = make(chan int, MAXBUFF)   // necessità di accodamento per priorità
var entrata_pedone_N = make(chan int, MAXBUFF) // necessità di accodamento per priorità
var entrata_pedone_S = make(chan int, MAXBUFF) // necessità di accodamento per priorità
var uscita_auto_N = make(chan int)
var uscita_auto_S = make(chan int)
var uscita_pedone_N = make(chan int)
var uscita_pedone_S = make(chan int)
var ACK_A_N [MAXPROC]chan int //risposte client nord
var ACK_A_S [MAXPROC]chan int //risposte client nord
var ACK_P_N [MAXPROC]chan int //risposte client sud
var ACK_P_S [MAXPROC]chan int //risposte client sud
var r int

func when(b bool, c chan int) chan int {
	if !b {
		return nil
	}
	return c
}
func pedone(myid int, dir int) {

	var tt int
	tt = rand.Intn(5) + 1
	fmt.Printf("inizializzazione pedone  %d direzione %d in secondi %d \n", myid, dir, tt)
	time.Sleep(time.Duration(tt) * time.Second)
	if dir == N {

		entrata_pedone_N <- myid // send asincrona
		<-ACK_P_N[myid]          // attesa x sincronizzazione
		fmt.Printf("[pedone %d]  sul ponte in direzione  NORD\n", myid)
		tt = rand.Intn(5)
		time.Sleep(time.Duration(tt) * time.Second)
		uscita_pedone_N <- myid
	} else {
		entrata_pedone_N <- myid
		<-ACK_P_S[myid] // attesa x sincronizzazione
		fmt.Printf("[pedone %d]  sul ponte in direzione  SUD\n", myid)
		tt = rand.Intn(5)
		time.Sleep(time.Duration(tt) * time.Second)
		uscita_pedone_S <- myid
	}
	done <- true
}
func veicolo(myid int, dir int) {
	var tt int

	tt = rand.Intn(5) + 1
	fmt.Printf("inizializzazione veicolo  %d direzione %d in secondi %d \n", myid, dir, tt)
	time.Sleep(time.Duration(tt) * time.Second)

	if dir == N {

		entrata_auto_N <- myid // send asincrona
		<-ACK_A_N[myid]        // attesa x sincronizzazione
		fmt.Printf("[veicolo %d]  sul ponte in direzione  NORD\n", myid)
		tt = rand.Intn(5)
		time.Sleep(time.Duration(tt) * time.Second)
		uscita_auto_N <- myid
	} else {
		entrata_auto_S <- myid
		<-ACK_A_S[myid] // attesa x sincronizzazione
		fmt.Printf("[veicolo %d]  sul ponte in direzione  SUD\n", myid)
		tt = rand.Intn(5)
		time.Sleep(time.Duration(tt) * time.Second)
		uscita_auto_S <- myid
	}
	done <- true
}

func gestore() {
	var contN int = 0
	var contS int = 0

	var cont_A_N int = 0
	var cont_A_S int = 0

	for {
		select {
		case x := <-when((contS+contN < MAXCAP) && (cont_A_N == 0), entrata_pedone_S):
			contS++
			fmt.Printf("[ponte]  entrato pedone %d in direzione S!  \n", x)
			ACK_P_S[x] <- 1 // termine "call"

		case x := <-when((contS+contN+10 < MAXCAP) && (cont_A_N == 0) && (len(entrata_pedone_S) == 0) && (len(entrata_pedone_N) == 0), entrata_auto_S):
			cont_A_S++
			contS += 10
			fmt.Printf("[ponte]  entrato veicolo %d in direzione S!  \n", x)
			ACK_A_S[x] <- 1 // termine "call"

		case x := <-when((contS+contN < MAXCAP) && (cont_A_S == 0) && (len(entrata_pedone_S) == 0), entrata_pedone_N):
			contN++
			fmt.Printf("[ponte]  entrato pedone %d in direzione N!  \n", x)
			ACK_P_N[x] <- 1 // termine "call"

		case x := <-when((contS+contN+10 < MAXCAP) && (cont_A_S == 0) && (len(entrata_pedone_S) == 0) && (len(entrata_pedone_N) == 0) && (len(entrata_auto_S) == 0), entrata_auto_N):
			cont_A_N++
			contN += 10
			fmt.Printf("[ponte]  entrato veicolo %d in direzione N!  \n", x)
			ACK_A_N[x] <- 1 // termine "call"
		case x := <-uscita_pedone_S:
			contS--
			fmt.Printf("[ponte]  uscito pedone %d in direzione S!  \n", x)
		case x := <-uscita_auto_S:
			cont_A_S--
			contS -= 10
			fmt.Printf("[ponte]  uscito veicolo %d in direzione S!  \n", x)
		case x := <-uscita_pedone_N:
			contN--
			fmt.Printf("[ponte]  uscito pedone %d in direzione N!  \n", x)
		case x := <-uscita_auto_N:
			cont_A_N--
			contN -= 10
			fmt.Printf("[ponte]  uscito veicolo %d in direzione N!  \n", x)
		case <-termina: // quando tutti i processi hanno finito
			fmt.Println("FINE !!!!!!")
			done <- true
			return
		}

	}

}
func main() {

}
