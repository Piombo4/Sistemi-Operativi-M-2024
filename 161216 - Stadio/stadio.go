package main

import (
	"fmt"
	"math/rand"
	"time"
)

const MAXBUFF = 100
const MAXPROC = 100

const LOCALI = 0
const OSPITI = 1

const NUM = 5

type Richiesta struct {
	id      int
	offerta float64
	ack     chan int
}

var done = make(chan bool)
var termina = make(chan bool)
var terminaStadio = make(chan bool)

var acquistaBiglietto = make(chan Richiesta, MAXBUFF)
var richiedi_controllo [2]chan Richiesta
var termina_controllo [2]chan int

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
func randFloats(min, max float64, n int) []float64 {
	res := make([]float64, n)
	for i := range res {
		res[i] = min + rand.Float64()*(max-min)
	}
	return res
}
func spettatore(myid int) {
	var tt int
	tt = rand.Intn(5) + 1
	r := Richiesta{id: myid, offerta: randFloats(1.0, 100.0, 1)[0], ack: make(chan int)}
	fmt.Printf("[SPETTATORE %d] INIZIALIZZAZIONE \n", myid)
	time.Sleep(time.Duration(tt) * time.Second)
	acquistaBiglietto <- r
	t := <-r.ack
	fmt.Printf("[SPETTATORE %d] Ha comprato un biglietto per la tribuna %d per %.2f€\n", myid, t, r.offerta)
	time.Sleep(time.Duration(tt) * time.Second)
	richiedi_controllo[t] <- r
	<-r.ack
	fmt.Printf("[SPETTATORE %d] Inizio controllo per la tribuna %d \n", myid, t)
	time.Sleep(time.Duration(tt) * time.Second)
	termina_controllo[t] <- myid
	done <- true
}
func biglietteria() {
	fmt.Printf("[BIGLIETTERIA] INIZIALIZZAZIONE \n")
	var tot float64 = 0
	for {
		select {
		case x := <-acquistaBiglietto:
			t := rand.Intn(2)
			tot += x.offerta
			fmt.Printf("[BIGLIETTERIA] Venduto biglietto a %d per tribuna %d - Totale: %.2f€\n", x.id, t, tot)
			x.ack <- t

		case <-termina:
			fmt.Printf("[BIGLIETTERIA] Termino\n")
			done <- true
			return
		default:
			continue
		}

	}
}
func stadio() {
	var nOp int
	var nSpet [2]int
	nOp = NUM
	fmt.Printf("[STADIO] INIZIALIZZAZIONE\n")
	for {
		select {
		case x := <-whenR(nOp > 0 && (nSpet[OSPITI] > nSpet[LOCALI] || len(richiedi_controllo[LOCALI]) == 0), richiedi_controllo[OSPITI]):
			fmt.Printf("[STADIO] Controllo dello spettatore %d\n", x.id)
			nOp--
			x.ack <- 1
		case x := <-whenR(nOp > 0 && (nSpet[LOCALI] >= nSpet[OSPITI] || len(richiedi_controllo[OSPITI]) == 0), richiedi_controllo[LOCALI]):
			fmt.Printf("[STADIO] Controllo dello spettatore %d\n", x.id)
			nOp--
			x.ack <- 1
		case x := <-termina_controllo[LOCALI]:
			nOp++
			nSpet[LOCALI]++
			fmt.Printf("[STADIO] Controllo dello spettatore %d terminato, totale tribuna LOCALI:%d \n", x, nSpet[LOCALI])
		case x := <-termina_controllo[OSPITI]:
			nOp++
			nSpet[OSPITI]++
			fmt.Printf("[STADIO] Controllo dello spettatore %d terminato, totale tribuna OSPITI:%d \n", x, nSpet[OSPITI])
		case <-terminaStadio:
			fmt.Printf("[STADIO] Termino\n")
			done <- true
			return
		}
	}
}
func main() {

	rand.Seed(time.Now().Unix())
	nSpettatori := 8
	go stadio()
	go biglietteria()

	for i := 0; i < 2; i++ {
		richiedi_controllo[i] = make(chan Richiesta, MAXBUFF)
		termina_controllo[i] = make(chan int, MAXBUFF)
	}

	for i := 0; i < nSpettatori; i++ {
		go spettatore(i)
	}
	for i := 0; i < nSpettatori; i++ {
		<-done
	}
	termina <- true
	<-done

	terminaStadio <- true
	<-done
	fmt.Printf("\n HO FINITO ")

}
