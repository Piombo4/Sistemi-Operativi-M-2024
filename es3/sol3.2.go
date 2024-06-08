// esercizio 3.2 - soluzione noleggio bici con guardie logiche
package main

import (
	"fmt"
	"math/rand"
	"time"
)

const MAXPROC = 100
const N_EB = 3
const N_BT = 10
const BT = 0
const EB = 1
const FLEX = 2
const DIMBUF = 300

type bici int //BT,EB

type req struct {
	id   int //identificatore cliente
	tipo int // BT, EB o FLEX
}

var richiestaBT = make(chan req, DIMBUF)
var richiestaEB = make(chan req, DIMBUF)
var richiestaFLEX = make(chan req, DIMBUF)
var rilascio = make(chan bici, DIMBUF)
var risorsa [MAXPROC]chan bici
var done = make(chan int)
var termina = make(chan int)

func when(b bool, c chan req) chan req {
	if !b {
		return nil
	}
	return c
}

func client(r req) {
	var b bici
	if r.tipo == BT {
		fmt.Printf("[cliente %d ] chiedo una bici tradizionale...\n", r.id)
		richiestaBT <- r //req vale BT
	} else if r.tipo == EB {
		fmt.Printf("[cliente %d ] chiedo una bici elettrica...\n", r.id)
		richiestaEB <- r //req vale EB
	} else {
		fmt.Printf("[cliente %d ] richiesta flex...\n", r.id)
		richiestaFLEX <- r //req vale BT, EB o FLEX
	}
	b = <-risorsa[r.id]
	if b == BT {
		fmt.Printf("[cliente %d ] ho ottenuto una bici tradizionale\n", r.id)
	} else {
		fmt.Printf("[cliente %d ] ho ottenuto una bici elettrica\n", r.id)
	}
	time.Sleep(time.Second * 2)

	rilascio <- b
	done <- r.id
}

func server() {

	var dispEB int = N_EB
	var dispBT int = N_BT
	var b bici
	var r req

	for {
		time.Sleep(time.Second * 1)
		select {
		case b = <-rilascio:
			switch b {
			case EB:
				{ //restituzione bici elettrica
					dispEB++
					fmt.Printf("[server]  restituita Bici elettrica.  \n")
				}
			case BT:
				{
					dispBT++
					fmt.Printf("[server]  restituita Bici tradizionale \n")
				}
			}
		case r = <-when(dispBT > 0, richiestaBT):
			dispBT--
			b = BT
			fmt.Printf("[server]  assegnata bici tradizionale a cliente %d \n", r.id)
			risorsa[r.id] <- b
		case r = <-when(dispEB > 0, richiestaEB):
			dispEB--
			b = EB
			fmt.Printf("[server]  assegnata e-bike a cliente %d \n", r.id)
			risorsa[r.id] <- b
		case r = <-when(dispEB > 0, richiestaFLEX):
			dispEB--
			b = EB
			fmt.Printf("[server]  assegnata e-bike a cliente %d \n", r.id)
			risorsa[r.id] <- b
		case r = <-when(dispEB == 0 && dispBT > 0, richiestaFLEX):
			dispBT--
			b = BT
			fmt.Printf("[server]  assegnata bici tradizionale a cliente flex %d \n", r.id)
			risorsa[r.id] <- b
		case r = <-when(dispEB == 0 && dispBT == 0, richiestaFLEX):
			fmt.Printf("[server]  il cliente flex %d si accoda per una ebike..\n", r.id)
			richiestaEB <- r // metto in attesa il proc. r per una EB (send asincrona!!)
		case <-termina: // quando tutti i processi clienti hanno finito
			fmt.Println("FINE !")
			done <- 1
			return

		}
	}
}

func main() {
	var cli int
	var r req

	rand.Seed(time.Now().Unix())
	fmt.Printf("\n quanti clienti (max %d)? ", MAXPROC)
	fmt.Scanf("%d", &cli)
	fmt.Println("clienti:", cli)

	//inizializzazione canali
	for i := 0; i < MAXPROC; i++ {
		risorsa[i] = make(chan bici, DIMBUF)
	}

	for i := 0; i < cli; i++ {
		r.id = i
		r.tipo = rand.Intn(3)
		go client(r)
	}
	go server()

	for i := 0; i < cli; i++ {
		<-done
	}
	termina <- 1 //terminazione server
	<-done
}
