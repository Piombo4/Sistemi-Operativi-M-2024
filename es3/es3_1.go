package main

import (
	"fmt"
	"math/rand"
	"time"
)

type req struct {
	id      int
	reqType int
}
type ril struct {
	id       int
	bikeType int
}

const MAXPROC = 10 //massimo numero di processi
const MAXBT = 2
const MAXEB = 2

var richiesta = make(chan req)
var rilascio = make(chan ril)
var risorsa [MAXPROC]chan ril
var done = make(chan int)
var termina = make(chan int)

func gestore() {
	var BTdisp int = MAXBT
	var EBdisp int = MAXEB

	var BTlibere [MAXBT]bool
	var EBlibere [MAXEB]bool

	var sospesi [MAXPROC]bool
	var nSosp int = 0

	var i int
	var r req
	var b ril

	for i := 0; i < MAXBT; i++ {
		BTlibere[i] = true
	}
	for i := 0; i < MAXEB; i++ {
		EBlibere[i] = true
	}
	for i := 0; i < MAXPROC; i++ {
		sospesi[i] = false
	}

	for {
		select {
		case r = <-richiesta:
			if r.reqType == 1 && BTdisp > 0 {
				BTdisp--
				risorsa[r.id] <- ril{r.id, 0}
				fmt.Printf("[server]  BT allocata a cliente %d \n", r.id)
			} else if r.reqType == 2 && EBdisp > 0 {
				EBdisp--
				risorsa[r.id] <- ril{r.id, 1}
				fmt.Printf("[server]  EB allocata a cliente %d \n", r.id)
			} else if r.reqType == 3 && (EBdisp > 0 || BTdisp > 0) {
				if EBdisp > 0 {
					EBdisp--
					risorsa[r.id] <- ril{r.id, 1}
					fmt.Printf("[server]  EB allocata a cliente %d \n", r.id)
				} else {
					BTdisp--
					risorsa[r.id] <- ril{r.id, 0}
					fmt.Printf("[server]  BT allocata a cliente %d \n", r.id)
				}
			} else {
				nSosp++
				sospesi[r.id] = true
				fmt.Printf("[server]  il cliente %d attende..\n", i)

			}

		case b = <-rilascio:
			if b.bikeType == 0 {
				if nSosp == 0 {
					BTdisp++
					fmt.Printf("[server]  restituita BT\n")
				} else {
					for i = 0; i < MAXPROC && !sospesi[i]; i++ {
					}
					sospesi[i] = false
					nSosp--
					risorsa[i] <- ril{i, 0}

				}
			} else {
				if nSosp == 0 {
					EBdisp++
					fmt.Printf("[server]  restituita EB\n")
				} else {
					for i = 0; i < MAXPROC && !sospesi[i]; i++ {
					}
					sospesi[i] = false
					nSosp--
					risorsa[i] <- ril{i, 1}

				}
			}
		case <-termina: // quando tutti i processi clienti hanno finito, il server termina
			fmt.Println("FINE")
			done <- 1
			return
		}
	}

}

func client(id int, r int) {
	var res ril
	request := req{id, r}
	richiesta <- request

	res = <-risorsa[id]
	fmt.Printf("\n[client %d] uso della risorsa %d\n", res.id, res.bikeType)
	timeout := rand.Intn(3)
	time.Sleep(time.Duration(timeout) * time.Second)
	rilascio <- res
	done <- 0
}
func main() {

	rand.Seed(time.Now().UnixNano())

	//inizializzazione canali
	for i := 0; i < MAXPROC; i++ {
		risorsa[i] = make(chan ril)
	}

	for i := 0; i < MAXPROC; i++ {
		randomNumber := rand.Intn(3) + 1
		go client(i, randomNumber)
	}
	go gestore()

	//attesa della terminazione dei clienti:
	for i := 0; i < MAXPROC; i++ {
		<-done
	}
	termina <- 1 //terminazione server
	<-done
}
