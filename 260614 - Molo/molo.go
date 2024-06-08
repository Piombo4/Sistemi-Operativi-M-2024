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
const MAX = 6
const IN = 0
const OUT = 1

var done = make(chan bool)
var termina = make(chan bool)
var terminaPasserella = make(chan bool)
var percorriPasserella [2]chan Richiesta
var esciPasserella = make(chan int)
var setPasserella = make(chan bool) // è sottinteso che viene chiusa/aperta solo in direzione IN

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
	case IN:
		return "IN"
	case OUT:
		return "OUT"

	}
	return ""
}
func viaggiatore(myid int) {
	var tt int
	tt = rand.Intn(5) + 1
	r := Richiesta{id: myid, ack: make(chan int)}
	tipo := rand.Intn(2)
	fmt.Printf("[Viaggiatore %d] "+getTipo((tipo))+" \n", myid)
	time.Sleep(time.Duration(tt) * time.Second)
	percorriPasserella[tipo] <- r
	<-r.ack
	//fmt.Printf("[Viaggiatore %d] Sto percorrendo la passerella "+getTipo((tipo))+" \n", myid)
	time.Sleep(time.Duration(5) * time.Second)
	esciPasserella<-myid
	//fmt.Printf("[Viaggiatore %d] Sono uscito dalla passerella \n", myid)
	time.Sleep(time.Duration(tt) * time.Second)
	done <- true
}

func addetto() {

	fmt.Printf("[Addetto] INIZIALIZZAZIONE \n")
	for {
		time.Sleep(time.Duration(4) * time.Second)
		setPasserella <- false
		fmt.Printf("[Addetto] Ho chiuso la passerella \n")
		time.Sleep(time.Duration(10) * time.Second)
		setPasserella <- true
		fmt.Printf("[Addetto] Ho aperto la passerella \n")

		select {
		case <-termina:
			fmt.Printf("[Addetto] Termino\n")
			done <- true
			return
		default:
			continue
		}
	}

}
func passerella() {
	var nPers int
	var aperta = true
	fmt.Printf("[Passerella] INIZIALIZZAZIONE\n")
	for {
		select {
		//ENTRATA
		case x := <-whenR(nPers < MAX, percorriPasserella[OUT]):
			nPers++
			x.ack <- 1
			fmt.Printf("[Passerella] il viaggiatore %d OUT sulla passerella, tot: %d\n", x.id,nPers)
		case x := <-whenR(nPers < MAX && len(percorriPasserella[OUT]) == 0 && aperta, percorriPasserella[IN]):
			nPers++
			x.ack <- 1
			fmt.Printf("[Passerella] il viaggiatore %d IN è sulla passerella, tot: %d\n", x.id,nPers)
			//USCITA
		case x := <-esciPasserella:
			nPers--
			fmt.Printf("[Passerella] il viaggiatore %d non è più sulla passerella, tot: %d\n", x,nPers)

			//ADDETTO
		case x := <-setPasserella:
			aperta = x
			if x {
				fmt.Printf("[Passerella] La passerella in direzione IN è aperta\n")
			} else {
				fmt.Printf("[Passerella] La passerella in direzione IN è chiusa\n")
			}

		case <-terminaPasserella:
			fmt.Printf("[Passerella] Termino\n")
			done <- true
			return
		}

	}
}
func main() {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 2; i++ {
		percorriPasserella[i] = make(chan Richiesta, MAXBUFF)
	}
	go passerella()
	go addetto()
	nViaggiatori := 15
	for i := 0; i < nViaggiatori; i++ {
		go viaggiatore(i)
	}
	for i := 0; i < nViaggiatori; i++ {
		<-done
	}

	//Terminiamo l'addetto
	termina <- true
	<-done
	//Terminiamo la passerella
	terminaPasserella <- true
	<-done
	fmt.Printf("\n HO FINITO ")
}
