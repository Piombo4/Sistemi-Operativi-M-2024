//soluzione esercizio 3.1
package main

import (
	"fmt"
	"time"
)

const MAXPROC = 100
const N_EB = 1
const N_BT = 30

const BT = 0
const EB = 1
const FLEX = 2

type bici int //BT,EB

type req struct {
	id   int //identificatore cliente
	tipo int // BT, EB o FLEX
}

var richiesta = make(chan req)
var rilascio = make(chan bici)
var risorsa [MAXPROC]chan bici
var done = make(chan int)
var termina = make(chan int)

var liberaEB [N_EB]bool // stato allocazione ebike
var liberaBT [N_BT]bool // stato allocazione bici tradizionali

func client(r req) {
	var b bici
	if r.tipo == BT {
		fmt.Printf("[cliente %d ] chiedo una bici tradizionale...\n", r.id)
	} else if r.tipo == EB {
		fmt.Printf("[cliente %d ] chiedo una bici elettrica...\n", r.id)
	} else {
		fmt.Printf("[cliente %d ] richiesta flex...\n", r.id)
	}
	richiesta <- r //req vale BT, EB o FLEX
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

	var sospEB [MAXPROC]bool

	var nsospEB int = 0
	var sospBT [MAXPROC]bool
	var nsospBT int = 0
	for i := 0; i < N_EB; i++ {
		liberaEB[i] = true
	}
	for i := 0; i < N_BT; i++ {
		liberaBT[i] = true
	}
	for i := 0; i < MAXPROC; i++ {
		sospBT[i] = false
		sospEB[i] = false
	}

	for {
		time.Sleep(time.Second * 1)
		select {
		case b = <-rilascio:
			if b == EB { //restituzione bici elettrica
				if nsospEB == 0 {
					dispEB++
					fmt.Printf("[server]  restituita Bici elettrica.  \n")
				} else {
					for i := 0; i < MAXPROC; i++ {
						if sospEB[i] == true {
							risorsa[i] <- b
							nsospEB--
							sospEB[i] = false
							break
						}
					}
				}
			} else { //restituzione bici tradizionale
				if nsospBT == 0 {
					dispBT++
					fmt.Printf("[server]  restituita Bici tradizionale \n")
				} else {
					for i := 0; i < MAXPROC; i++ {
						if sospBT[i] == true {
							risorsa[i] <- b
							nsospBT--
							sospBT[i] = false
							break
						}
					}

				}
			}

		case r = <-richiesta:
			switch r.tipo {
			case FLEX: // richiesta  FLEX
				if dispEB > 0 {
					dispEB--
					b = EB
					fmt.Printf("[server]  allocata bici elettrica  a cliente %d \n", r.id)
					risorsa[r.id] <- b
				} else if dispBT > 0 {
					dispBT--
					b = BT
					fmt.Printf("[server]  allocata bici tradizionale a cliente FLEX %d \n", r.id)
					risorsa[r.id] <- b
				} else {
					//ATTESA
					nsospEB++
					sospEB[r.id] = true
				}
			case BT: // richiesta  BT
				if dispBT > 0 {
					dispBT--
					b = BT
					fmt.Printf("[server]  allocata bici tradizionale a cliente %d \n", r.id)
					risorsa[r.id] <- b
				} else {
					//ATTESA
					nsospBT++
					sospBT[r.id] = true
				}
			case EB:
				if dispEB > 0 {
					dispEB--
					b = EB
					fmt.Printf("[server]  allocata bici elettrica  a cliente %d \n", r.id)
					risorsa[r.id] <- b
				} else {
					//ATTESA
					nsospEB++
					sospEB[r.id] = true
				}
			}
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
	fmt.Printf("\n quanti clienti (max %d)? ", MAXPROC)
	fmt.Scanf("%d", &cli)
	fmt.Println("clienti:", cli)

	//inizializzazione canali
	for i := 0; i < MAXPROC; i++ {
		risorsa[i] = make(chan bici)
	}

	for i := 0; i < cli; i++ {
		r.id = i
		r.tipo = i % 3
		go client(r)
	}
	go server()

	for i := 0; i < cli; i++ {
		<-done
	}
	termina <- 1 //terminazione server
	<-done
}
