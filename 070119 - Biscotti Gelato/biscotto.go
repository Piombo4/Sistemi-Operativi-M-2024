package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Richiesta struct {
	ack chan int
}

const MAXBUFF = 100
const MAX = 8
const N = 2
const M = 5
const TOT = 10

var currGelati = 0
var done = make(chan bool)
var termina = make(chan bool)

var depositoBiscotto = make(chan Richiesta)
var richiediRifornimento = make(chan Richiesta)
var prelievoBiscotti = make(chan Richiesta)

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
func mb() {
	var tt int
	tt = rand.Intn(5) + 1
	r := Richiesta{ack: make(chan int)}
	fmt.Printf("[MB] INIZIALIZZAZIONE \n")
	for currGelati < TOT {
		if currGelati == TOT {
			break
		}
		fmt.Printf("[MB] Produzione di un biscotto... \n")
		time.Sleep(time.Duration(tt) * time.Second)
		depositoBiscotto <- r
		<-r.ack
		fmt.Printf("[MB] Depositato un biscotto! Totale: %d\n", currGelati)
	}
	done <- true

}
func mg() {
	var tt int
	tt = rand.Intn(5) + 1
	r := Richiesta{ack: make(chan int)}
	serbatoio := M
	fmt.Printf("[MG] INIZIALIZZAZIONE \n")
	for currGelati < TOT {
		if currGelati == TOT {
			break
		}
		time.Sleep(time.Duration(tt) * time.Second)
		prelievoBiscotti <- r
		<-r.ack
		fmt.Printf("[MG] Prelevati 2 biscotti! \n")
		if serbatoio > 0 {
			serbatoio--
		} else {
			fmt.Printf("[MG] Serbatoio vuoto, chiamo l'operaio! \n")
			richiediRifornimento <- r
			<-r.ack
			serbatoio = M - 1 //ne usa uno per il gelato
			fmt.Printf("[MG] Rifornimento ottenuto! \n")
		}
		currGelati++
		fmt.Printf("[MG] Gelato fatto! Totale: %d\n", currGelati)
	}
	done <- true

}
func operaio() {
	var tt int
	tt = rand.Intn(5) + 1
	fmt.Printf("[OPERAIO] INIZIALIZZAZIONE \n")
	time.Sleep(time.Duration(tt) * time.Second)
	for {

		select {
		case x := <-richiediRifornimento:
			x.ack <- 1
			fmt.Printf("[OPERAIO] Rifornisco il serbatoio! \n")
		case <-termina:
			fmt.Printf("[OPERAIO] Termino\n")
			done <- true
			return
		default:
			continue
		}
	}
}
func alimentatore() {

	var currCap = 0

	fmt.Printf("[ALIMENTATORE] INIZIALIZZAZIONE \n")

	for {

		select {
		case x := <-whenR(currCap < MAX && (currCap >= MAX/2 || len(prelievoBiscotti) == 0), depositoBiscotto):
			currCap++
			x.ack <- 1
			fmt.Printf("[ALIMENTATORE] Biscotto depositato! \n")
		case x := <-whenR(currCap > 1 && (currCap < MAX/2 || len(depositoBiscotto) == 0), prelievoBiscotti):
			currCap -= N
			x.ack <- 1
			fmt.Printf("[ALIMENTATORE] Biscotto prelevato! \n")
		case <-termina:
			fmt.Printf("[ALIMENTATORE] Termino\n")
			done <- true
			return
		default:
			continue
		}
	}

}
func main() {
	rand.Seed(time.Now().Unix())

	go alimentatore()
	go operaio()
	go mg()
	go mb()

	<-done
	<-done
	termina <- true
	<-done

	termina <- true
	<-done

	fmt.Printf("\n HO FINITO ")
}
