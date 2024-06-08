package main

import (
	"fmt"
	"math/rand"
	"time"
)

const MAXBUFF = 100;

// ENUM

const (
	MARTELLO int = 0;
	TENAGLIA int = 1;
	BADILE int = 2;
)

const (
	NEGOZIANTE int = 0;
	PRIVATO int = 1;
)

const CLIENTS = 40

var tipoClienti = [2]string{"NEGOZIANTE", "PRIVATO"}
var tipoArticoli = [3]string{"Martello", "Tenaglia","Badile"}

const K = 2;
const MAX = 6;

var done = make(chan bool)
var termina = make(chan bool)

var FINITO bool = false

var acquista [2][3] chan int;
var ackAcquisto [2][3] chan int;
var fornisci[3] chan int;
var ackFornisci[3] chan int;


func when(b bool, c chan int) chan int {
	if !b {
		return nil
	}
	return c
}

func Magazzino() {

	var articoli[3] int = [3]int{0,0,0}
	
	for {
		time.Sleep(time.Duration(1))

		select {
		case <- when(articoli[MARTELLO] - K >= 0,acquista[NEGOZIANTE][MARTELLO]):

			articoli[MARTELLO] -= K
			ackAcquisto[NEGOZIANTE][MARTELLO] <- 1

			fmt.Printf("\n[Magazzino] comprati K martelli da negoziante")
			fmt.Printf("\n[Magazzino] STAT: articoli[M,T,B]: %d,%d,%d, WAITING: acquista[NEGOZ:[M,T,B],PRIV:[M,T,B]]: [%d,%d,%d],[%d,%d,%d], fornisci[M,T,B]: %d,%d,%d", articoli[MARTELLO], articoli[TENAGLIA], articoli[BADILE],len(acquista[NEGOZIANTE][MARTELLO]),len(acquista[NEGOZIANTE][TENAGLIA]),len(acquista[NEGOZIANTE][BADILE]),len(acquista[PRIVATO][MARTELLO]),len(acquista[PRIVATO][TENAGLIA]),len(acquista[PRIVATO][BADILE]),len(fornisci[MARTELLO]),len(fornisci[TENAGLIA]),len(fornisci[BADILE]))
			fmt.Printf("\n")
		
		case<- when(articoli[TENAGLIA] - K >= 0,acquista[NEGOZIANTE][TENAGLIA]):

			articoli[TENAGLIA] -= K
			ackAcquisto[NEGOZIANTE][TENAGLIA] <- 1

			fmt.Printf("\n[Magazzino] comprati K tenaglie da negoziante")
			fmt.Printf("\n[Magazzino] STAT: articoli[M,T,B]: %d,%d,%d, WAITING: acquista[NEGOZ:[M,T,B],PRIV:[M,T,B]]: [%d,%d,%d],[%d,%d,%d], fornisci[M,T,B]: %d,%d,%d", articoli[MARTELLO], articoli[TENAGLIA], articoli[BADILE],len(acquista[NEGOZIANTE][MARTELLO]),len(acquista[NEGOZIANTE][TENAGLIA]),len(acquista[NEGOZIANTE][BADILE]),len(acquista[PRIVATO][MARTELLO]),len(acquista[PRIVATO][TENAGLIA]),len(acquista[PRIVATO][BADILE]),len(fornisci[MARTELLO]),len(fornisci[TENAGLIA]),len(fornisci[BADILE]))
			fmt.Printf("\n")

		case <- when(articoli[BADILE] - K >= 0,acquista[NEGOZIANTE][BADILE]):

			articoli[BADILE] -= K
			ackAcquisto[NEGOZIANTE][BADILE] <- 1

			fmt.Printf("\n[Magazzino] comprati K badili da negoziante")
			fmt.Printf("\n[Magazzino] STAT: articoli[M,T,B]: %d,%d,%d, WAITING: acquista[NEGOZ:[M,T,B],PRIV:[M,T,B]]: [%d,%d,%d],[%d,%d,%d], fornisci[M,T,B]: %d,%d,%d", articoli[MARTELLO], articoli[TENAGLIA], articoli[BADILE],len(acquista[NEGOZIANTE][MARTELLO]),len(acquista[NEGOZIANTE][TENAGLIA]),len(acquista[NEGOZIANTE][BADILE]),len(acquista[PRIVATO][MARTELLO]),len(acquista[PRIVATO][TENAGLIA]),len(acquista[PRIVATO][BADILE]),len(fornisci[MARTELLO]),len(fornisci[TENAGLIA]),len(fornisci[BADILE]))
			fmt.Printf("\n")
		
		case <- when(articoli[MARTELLO] - 1 >= 0 && (len(acquista[NEGOZIANTE][MARTELLO]) +len(acquista[NEGOZIANTE][TENAGLIA]) + len(acquista[NEGOZIANTE][BADILE]) == 0),acquista[PRIVATO][MARTELLO]):

			articoli[MARTELLO]--
			ackAcquisto[PRIVATO][MARTELLO] <- 1

			fmt.Printf("\n[Magazzino] comprati un martello da privato")
			fmt.Printf("\n[Magazzino] STAT: articoli[M,T,B]: %d,%d,%d, WAITING: acquista[NEGOZ:[M,T,B],PRIV:[M,T,B]]: [%d,%d,%d],[%d,%d,%d], fornisci[M,T,B]: %d,%d,%d", articoli[MARTELLO], articoli[TENAGLIA], articoli[BADILE],len(acquista[NEGOZIANTE][MARTELLO]),len(acquista[NEGOZIANTE][TENAGLIA]),len(acquista[NEGOZIANTE][BADILE]),len(acquista[PRIVATO][MARTELLO]),len(acquista[PRIVATO][TENAGLIA]),len(acquista[PRIVATO][BADILE]),len(fornisci[MARTELLO]),len(fornisci[TENAGLIA]),len(fornisci[BADILE]))
			fmt.Printf("\n")
		
		case <- when(articoli[TENAGLIA] - 1 >= 0 && (len(acquista[NEGOZIANTE][MARTELLO]) +len(acquista[NEGOZIANTE][TENAGLIA]) + len(acquista[NEGOZIANTE][BADILE]) == 0),acquista[PRIVATO][TENAGLIA]):

			articoli[TENAGLIA]--
			ackAcquisto[PRIVATO][TENAGLIA] <- 1

			fmt.Printf("\n[Magazzino] comprato una tenaglia da privato")
			fmt.Printf("\n[Magazzino] STAT: articoli[M,T,B]: %d,%d,%d, WAITING: acquista[NEGOZ:[M,T,B],PRIV:[M,T,B]]: [%d,%d,%d],[%d,%d,%d], fornisci[M,T,B]: %d,%d,%d", articoli[MARTELLO], articoli[TENAGLIA], articoli[BADILE],len(acquista[NEGOZIANTE][MARTELLO]),len(acquista[NEGOZIANTE][TENAGLIA]),len(acquista[NEGOZIANTE][BADILE]),len(acquista[PRIVATO][MARTELLO]),len(acquista[PRIVATO][TENAGLIA]),len(acquista[PRIVATO][BADILE]),len(fornisci[MARTELLO]),len(fornisci[TENAGLIA]),len(fornisci[BADILE]))
			fmt.Printf("\n")

		case <- when(articoli[BADILE] - 1 >= 0 && (len(acquista[NEGOZIANTE][MARTELLO]) +len(acquista[NEGOZIANTE][TENAGLIA]) + len(acquista[NEGOZIANTE][BADILE]) == 0),acquista[PRIVATO][BADILE]):

			articoli[BADILE]--
			ackAcquisto[PRIVATO][BADILE] <- 1

			fmt.Printf("\n[Magazzino] comprati un badile da privato")
			fmt.Printf("\n[Magazzino] STAT: articoli[M,T,B]: %d,%d,%d, WAITING: acquista[NEGOZ:[M,T,B],PRIV:[M,T,B]]: [%d,%d,%d],[%d,%d,%d], fornisci[M,T,B]: %d,%d,%d", articoli[MARTELLO], articoli[TENAGLIA], articoli[BADILE],len(acquista[NEGOZIANTE][MARTELLO]),len(acquista[NEGOZIANTE][TENAGLIA]),len(acquista[NEGOZIANTE][BADILE]),len(acquista[PRIVATO][MARTELLO]),len(acquista[PRIVATO][TENAGLIA]),len(acquista[PRIVATO][BADILE]),len(fornisci[MARTELLO]),len(fornisci[TENAGLIA]),len(fornisci[BADILE]))
			fmt.Printf("\n")
		
		case x := <- when(!FINITO && articoli[MARTELLO] < MAX, fornisci[MARTELLO]):

			if articoli[MARTELLO] + x > MAX {
				articoli[MARTELLO] = MAX
			} else {
				articoli[MARTELLO] += x
			}

			ackFornisci[MARTELLO] <- 1

			fmt.Printf("\n[Magazzino] forniti %d martelli", x)
			fmt.Printf("\n[Magazzino] STAT: articoli[M,T,B]: %d,%d,%d, WAITING: acquista[NEGOZ:[M,T,B],PRIV:[M,T,B]]: [%d,%d,%d],[%d,%d,%d], fornisci[M,T,B]: %d,%d,%d", articoli[MARTELLO], articoli[TENAGLIA], articoli[BADILE],len(acquista[NEGOZIANTE][MARTELLO]),len(acquista[NEGOZIANTE][TENAGLIA]),len(acquista[NEGOZIANTE][BADILE]),len(acquista[PRIVATO][MARTELLO]),len(acquista[PRIVATO][TENAGLIA]),len(acquista[PRIVATO][BADILE]),len(fornisci[MARTELLO]),len(fornisci[TENAGLIA]),len(fornisci[BADILE]))
			fmt.Printf("\n")

		case x := <- when(!FINITO && articoli[TENAGLIA] < MAX && len(fornisci[MARTELLO]) == 0, fornisci[TENAGLIA]):

			if articoli[TENAGLIA] + x > MAX {
				articoli[TENAGLIA] = MAX
			} else {
				articoli[TENAGLIA] += x
			}

			ackFornisci[TENAGLIA] <- 1

			fmt.Printf("\n[Magazzino] forniti %d tenaglie", x)
			fmt.Printf("\n[Magazzino] STAT: articoli[M,T,B]: %d,%d,%d, WAITING: acquista[NEGOZ:[M,T,B],PRIV:[M,T,B]]: [%d,%d,%d],[%d,%d,%d], fornisci[M,T,B]: %d,%d,%d", articoli[MARTELLO], articoli[TENAGLIA], articoli[BADILE],len(acquista[NEGOZIANTE][MARTELLO]),len(acquista[NEGOZIANTE][TENAGLIA]),len(acquista[NEGOZIANTE][BADILE]),len(acquista[PRIVATO][MARTELLO]),len(acquista[PRIVATO][TENAGLIA]),len(acquista[PRIVATO][BADILE]),len(fornisci[MARTELLO]),len(fornisci[TENAGLIA]),len(fornisci[BADILE]))
			fmt.Printf("\n")

		case x := <- when(!FINITO && articoli[BADILE] < MAX && len(fornisci[TENAGLIA]) == 0 && len (fornisci[TENAGLIA]) == 0, fornisci[BADILE]):

			if articoli[BADILE] + x > MAX {
				articoli[BADILE] = MAX
			} else {
				articoli[BADILE] += x
			}

			ackFornisci[BADILE] <- 1

			fmt.Printf("\n[Magazzino] forniti %d badili", x)
			fmt.Printf("\n[Magazzino] STAT: articoli[M,T,B]: %d,%d,%d, WAITING: acquista[NEGOZ:[M,T,B],PRIV:[M,T,B]]: [%d,%d,%d],[%d,%d,%d], fornisci[M,T,B]: %d,%d,%d", articoli[MARTELLO], articoli[TENAGLIA], articoli[BADILE],len(acquista[NEGOZIANTE][MARTELLO]),len(acquista[NEGOZIANTE][TENAGLIA]),len(acquista[NEGOZIANTE][BADILE]),len(acquista[PRIVATO][MARTELLO]),len(acquista[PRIVATO][TENAGLIA]),len(acquista[PRIVATO][BADILE]),len(fornisci[MARTELLO]),len(fornisci[TENAGLIA]),len(fornisci[BADILE]))
			fmt.Printf("\n")
	
		case <- when(FINITO, fornisci[MARTELLO]):
			ackFornisci[MARTELLO] <- -1
		
		case <- when(FINITO, fornisci[TENAGLIA]):
			ackFornisci[TENAGLIA] <- -1

		case <- when(FINITO, fornisci[BADILE]):
			ackFornisci[BADILE] <- -1
	
		case <- termina:
			fmt.Printf("\n[Magazzino] termino!")
			done <- true
			return
	}
	}
}

func Cliente(id int, tipo int, articolo int) {

	fmt.Printf("\n[Cliente %d] iniziato cliente %s che acquista %s", id, tipoClienti[tipo], tipoArticoli[articolo])

	acquista[tipo][articolo] <- 1
	<- ackAcquisto[tipo][articolo]

	fmt.Printf("\n[Cliente %d] cliente %s ha acquistato %s", id, tipoClienti[tipo], tipoArticoli[articolo])
	done <- true
	fmt.Printf("\n[Cliente %d] termino!", id)
}

func Fornitore(articolo int) {

	fmt.Printf("\n[Fornitore %s] partito", tipoArticoli[articolo])

	var ris int
	var x int

	for {
		time.Sleep(time.Duration((rand.Intn(5) + 1)) * time.Second)

		x = (rand.Intn(3) + 1)
		fornisci[articolo] <- x
		ris = <- ackFornisci[articolo]

		if ( ris == -1 ) {
			fmt.Printf("\n[Fornitore %s] termino", tipoArticoli[articolo])
			done <- true
			return
		}
		

		fmt.Printf("\n[Fornitore %s] consegnati %d articoli", tipoArticoli[articolo],x)
	}
}

func main() {
	rand.Seed(time.Now().Unix())


	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			acquista[i][j] = make(chan int, MAXBUFF)
			ackAcquisto[i][j] = make(chan int)
		}
	}

	for i := 0; i < 3; i++ {
		fornisci[i] = make(chan int, MAXBUFF)
		ackFornisci[i] = make(chan int)
	}

	go Magazzino()

	go Fornitore(MARTELLO)
	go Fornitore(TENAGLIA)
	go Fornitore(BADILE)

	for i := 0; i < CLIENTS; i++ {
		go Cliente(i,rand.Intn(2),rand.Intn(3))
	}

	for i := 0; i < CLIENTS; i++ {
		<-done
	}

	FINITO = true

	for i := 0; i < 3; i++ {
		<-done
	}

	termina <- true
	<-done
	fmt.Printf("[main] APPLICAZIONE TERMINATA \n")
}