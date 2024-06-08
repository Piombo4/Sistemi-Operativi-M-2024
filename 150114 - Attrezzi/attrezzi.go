/*
* In questa soluzione quando un fornitore consegna una certa quantità
* di oggetti, se la capObj + qtaRifornimento > MAX => capObj = MAX.
* In questo modo non si sfora. Non so se il testo voleva che i fornitori che consegnano più di MAX vengano bloccati
 */
package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Richiesta struct {
	id  int
	ack chan int
}

const MAXBUFF = 100
const MAXC = 15
const MAX = 6
const K = 2
const MARTELLO = 0
const TENAGLIA = 1
const BADILE = 2
const NEGOZIANTE = 0
const PRIVATO = 1

var done = make(chan bool)
var terminaMagazzino = make(chan bool)
var acquistaPrivato [3]chan Richiesta
var acquistaNegoziante [3]chan Richiesta
var rifornimento [3]chan Richiesta

var finito = false
var tot = MAXC //per contare i clienti ancora attivi
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

func getTipo(tipo int) string {
	switch tipo {
	case NEGOZIANTE:
		return "negoziante"
	case PRIVATO:
		return "privato"
	}
	return ""
}
func getObj(tipo int) string {
	switch tipo {
	case MARTELLO:
		return "martello"
	case BADILE:
		return "badile"
	case TENAGLIA:
		return "tenaglia"
	}
	return ""
}
func cliente(myid int) {
	var tt int
	tt = rand.Intn(5) + 1
	tipo := rand.Intn(2)
	obj := rand.Intn(3)
	r := Richiesta{id: myid, ack: make(chan int)}
	fmt.Printf("[Cliente %d %s] INIZIALIZZAZIONE \n", myid, getTipo(tipo))
	time.Sleep(time.Duration(tt) * time.Second)

	switch tipo {
	case NEGOZIANTE:
		acquistaNegoziante[obj] <- r
		<-r.ack
		//fmt.Printf("[Cliente %d %s] Ho comprato %d %s \n", myid, getTipo(tipo), K, getObj(obj))
	case PRIVATO:
		acquistaPrivato[obj] <- r
		<-r.ack
		//fmt.Printf("[Cliente %d %s] Ho comprato un/una %s \n", myid, getTipo(tipo), getObj(obj))
	}
	time.Sleep(time.Duration(tt) * time.Second)
	tot--
	fmt.Printf("[Cliente %d %s] Termino, RESTANTI:%d \n", myid, getTipo(tipo), tot)
	done <- true
}

func fornitore(tipo int) {
	var tt int
	tt = rand.Intn(5) + 1
	r := Richiesta{id: -1, ack: make(chan int)}
	fmt.Printf("[Fornitore %d] INIZIALIZZAZIONE \n", tipo)
	var ris int
	for {
		qty := rand.Intn(MAX) + 1 //quanto rifornisce
		r.id = qty
		time.Sleep(time.Duration(tt) * time.Second)
		rifornimento[tipo] <- r
		ris = <-r.ack
		if ris < 0 {
			fmt.Printf("[Fornitore %d] Termino\n", tipo)
			done <- true
			return
		}
		fmt.Printf("[Fornitore %d] Ho aggiunto %d %s\n", tipo, qty, getObj(tipo))
	}

}
func canBuy() bool {
	return len(acquistaNegoziante[BADILE]) == 0 && len(acquistaNegoziante[MARTELLO]) == 0 && len(acquistaNegoziante[TENAGLIA]) == 0
}
func printTot(nObj []int) string {
	return "(M: " + strconv.Itoa(nObj[MARTELLO]) + " B: " + strconv.Itoa(nObj[BADILE]) + " T: " + strconv.Itoa(nObj[TENAGLIA]) + ")"
}
func magazzino() {

	var nObj [3]int
	for i := 0; i < 3; i++ {
		nObj[i] = MAX
	}

	fmt.Printf("[Magazzino] INIZIALIZZAZIONE\n")
	for {
		select {
		//NEGOZIANTI
		case x := <-whenR(nObj[MARTELLO] >= K, acquistaNegoziante[MARTELLO]):
			nObj[MARTELLO] -= K

			fmt.Printf("[Magazzino] Negoziante %d ha acquistato %d martelli, %s\n", x.id, K, printTot(nObj[:]))
			fmt.Printf("[Magazzino] Waiting:acquista[NEGOZ:[M,B,T],PRIV:[M,B,T]]: [%d,%d,%d],[%d,%d,%d], fornisci[M,B,T]: %d,%d,%d\n", len(acquistaNegoziante[MARTELLO]), len(acquistaNegoziante[BADILE]), len(acquistaNegoziante[TENAGLIA]), len(acquistaPrivato[MARTELLO]), len(acquistaPrivato[BADILE]), len(acquistaPrivato[TENAGLIA]), len(rifornimento[MARTELLO]), len(rifornimento[BADILE]), len(rifornimento[TENAGLIA]))

			x.ack <- 1
		case x := <-whenR(nObj[BADILE] >= K, acquistaNegoziante[BADILE]):
			nObj[BADILE] -= K

			fmt.Printf("[Magazzino] Negoziante %d ha acquistato %d badili, %s\n", x.id, K, printTot(nObj[:]))
			fmt.Printf("[Magazzino] Waiting:acquista[NEGOZ:[M,B,T],PRIV:[M,B,T]]: [%d,%d,%d],[%d,%d,%d], fornisci[M,B,T]: %d,%d,%d\n", len(acquistaNegoziante[MARTELLO]), len(acquistaNegoziante[BADILE]), len(acquistaNegoziante[TENAGLIA]), len(acquistaPrivato[MARTELLO]), len(acquistaPrivato[BADILE]), len(acquistaPrivato[TENAGLIA]), len(rifornimento[MARTELLO]), len(rifornimento[BADILE]), len(rifornimento[TENAGLIA]))

			x.ack <- 1
		case x := <-whenR(nObj[TENAGLIA] >= K, acquistaNegoziante[TENAGLIA]):
			nObj[TENAGLIA] -= K
			fmt.Printf("[Magazzino] Negoziante %d ha acquistato %d tenaglie, %s\n", x.id, K, printTot(nObj[:]))
			fmt.Printf("[Magazzino] Waiting:acquista[NEGOZ:[M,B,T],PRIV:[M,B,T]]: [%d,%d,%d],[%d,%d,%d], fornisci[M,B,T]: %d,%d,%d\n", len(acquistaNegoziante[MARTELLO]), len(acquistaNegoziante[BADILE]), len(acquistaNegoziante[TENAGLIA]), len(acquistaPrivato[MARTELLO]), len(acquistaPrivato[BADILE]), len(acquistaPrivato[TENAGLIA]), len(rifornimento[MARTELLO]), len(rifornimento[BADILE]), len(rifornimento[TENAGLIA]))

			x.ack <- 1
			//PRIVATI
		case x := <-whenR(canBuy() && nObj[MARTELLO] > 0, acquistaPrivato[MARTELLO]):
			nObj[MARTELLO]--
			fmt.Printf("[Magazzino] Privato %d ha acquistato un martello, %s\n", x.id, printTot(nObj[:]))
			fmt.Printf("[Magazzino] Waiting:acquista[NEGOZ:[M,B,T],PRIV:[M,B,T]]: [%d,%d,%d],[%d,%d,%d], fornisci[M,B,T]: %d,%d,%d\n", len(acquistaNegoziante[MARTELLO]), len(acquistaNegoziante[BADILE]), len(acquistaNegoziante[TENAGLIA]), len(acquistaPrivato[MARTELLO]), len(acquistaPrivato[BADILE]), len(acquistaPrivato[TENAGLIA]), len(rifornimento[MARTELLO]), len(rifornimento[BADILE]), len(rifornimento[TENAGLIA]))

			x.ack <- 1
		case x := <-whenR(canBuy() && nObj[BADILE] > 0, acquistaPrivato[BADILE]):
			nObj[BADILE]--
			fmt.Printf("[Magazzino] Privato %d ha acquistato un badile, %s\n", x.id, printTot(nObj[:]))
			fmt.Printf("[Magazzino] Waiting:acquista[NEGOZ:[M,B,T],PRIV:[M,B,T]]: [%d,%d,%d],[%d,%d,%d], fornisci[M,B,T]: %d,%d,%d\n", len(acquistaNegoziante[MARTELLO]), len(acquistaNegoziante[BADILE]), len(acquistaNegoziante[TENAGLIA]), len(acquistaPrivato[MARTELLO]), len(acquistaPrivato[BADILE]), len(acquistaPrivato[TENAGLIA]), len(rifornimento[MARTELLO]), len(rifornimento[BADILE]), len(rifornimento[TENAGLIA]))

			x.ack <- 1
		case x := <-whenR(canBuy() && nObj[TENAGLIA] > 0, acquistaPrivato[TENAGLIA]):
			nObj[TENAGLIA]--
			fmt.Printf("[Magazzino] Privato %d ha acquistato una tenaglia, %s\n", x.id, printTot(nObj[:]))
			fmt.Printf("[Magazzino] Waiting:acquista[NEGOZ:[M,B,T],PRIV:[M,B,T]]: [%d,%d,%d],[%d,%d,%d], fornisci[M,B,T]: %d,%d,%d\n", len(acquistaNegoziante[MARTELLO]), len(acquistaNegoziante[BADILE]), len(acquistaNegoziante[TENAGLIA]), len(acquistaPrivato[MARTELLO]), len(acquistaPrivato[BADILE]), len(acquistaPrivato[TENAGLIA]), len(rifornimento[MARTELLO]), len(rifornimento[BADILE]), len(rifornimento[TENAGLIA]))

			x.ack <- 1
			//FORNITORI
		case x := <-whenR(!finito && nObj[MARTELLO] < MAX, rifornimento[MARTELLO]):
			nObj[MARTELLO] = min(nObj[MARTELLO]+x.id, MAX)
			fmt.Printf("[Magazzino] Rifornitore ha aggiunto %d martelli, %s\n", x.id, printTot(nObj[:]))
			fmt.Printf("[Magazzino] Waiting:acquista[NEGOZ:[M,B,T],PRIV:[M,B,T]]: [%d,%d,%d],[%d,%d,%d], fornisci[M,B,T]: %d,%d,%d\n", len(acquistaNegoziante[MARTELLO]), len(acquistaNegoziante[BADILE]), len(acquistaNegoziante[TENAGLIA]), len(acquistaPrivato[MARTELLO]), len(acquistaPrivato[BADILE]), len(acquistaPrivato[TENAGLIA]), len(rifornimento[MARTELLO]), len(rifornimento[BADILE]), len(rifornimento[TENAGLIA]))

			x.ack <- 1
		case x := <-whenR(!finito && nObj[BADILE] < MAX && len(rifornimento[MARTELLO]) == 0 && len(rifornimento[TENAGLIA]) == 0, rifornimento[BADILE]):
			nObj[BADILE] = min(nObj[BADILE]+x.id, MAX)
			fmt.Printf("[Magazzino] Rifornitore ha aggiunto %d badili, %s\n", x.id, printTot(nObj[:]))
			fmt.Printf("[Magazzino] Waiting:acquista[NEGOZ:[M,B,T],PRIV:[M,B,T]]: [%d,%d,%d],[%d,%d,%d], fornisci[M,B,T]: %d,%d,%d\n", len(acquistaNegoziante[MARTELLO]), len(acquistaNegoziante[BADILE]), len(acquistaNegoziante[TENAGLIA]), len(acquistaPrivato[MARTELLO]), len(acquistaPrivato[BADILE]), len(acquistaPrivato[TENAGLIA]), len(rifornimento[MARTELLO]), len(rifornimento[BADILE]), len(rifornimento[TENAGLIA]))

			x.ack <- 1
		case x := <-whenR(!finito && nObj[TENAGLIA] < MAX && len(rifornimento[MARTELLO]) == 0, rifornimento[TENAGLIA]):
			nObj[TENAGLIA] = min(nObj[TENAGLIA]+x.id, MAX)

			fmt.Printf("[Magazzino] Rifornitore ha aggiunto %d tenaglie, %s\n", x.id, printTot(nObj[:]))
			fmt.Printf("[Magazzino] Waiting:acquista[NEGOZ:[M,B,T],PRIV:[M,B,T]]: [%d,%d,%d],[%d,%d,%d], fornisci[M,B,T]: %d,%d,%d\n", len(acquistaNegoziante[MARTELLO]), len(acquistaNegoziante[BADILE]), len(acquistaNegoziante[TENAGLIA]), len(acquistaPrivato[MARTELLO]), len(acquistaPrivato[BADILE]), len(acquistaPrivato[TENAGLIA]), len(rifornimento[MARTELLO]), len(rifornimento[BADILE]), len(rifornimento[TENAGLIA]))

			x.ack <- 1
		case x := <-whenR(finito, rifornimento[MARTELLO]):
			x.ack <- -1

		case x := <-whenR(finito, rifornimento[TENAGLIA]):
			x.ack <- -1

		case x := <-whenR(finito, rifornimento[BADILE]):
			x.ack <- -1
		case <-terminaMagazzino:
			fmt.Printf("[Magazzino] Termino\n")
			done <- true
			return
		}

	}
}
func main() {
	rand.Seed(time.Now().Unix())
	for i := 0; i < 3; i++ {
		acquistaNegoziante[i] = make(chan Richiesta, MAXBUFF)
		acquistaPrivato[i] = make(chan Richiesta, MAXBUFF)
		rifornimento[i] = make(chan Richiesta, MAXBUFF)
	}
	go magazzino()

	for i := 0; i < 3; i++ {
		go fornitore(i)
	}
	for i := 0; i < MAXC; i++ {
		go cliente(i)
	}
	for i := 0; i < MAXC; i++ {
		<-done
	}
	finito = true
	for i := 0; i < 3; i++ {

		<-done
	}

	terminaMagazzino <- true
	<-done
	fmt.Printf("\n HO FINITO ")
}
