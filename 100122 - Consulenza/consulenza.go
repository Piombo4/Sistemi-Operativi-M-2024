package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Richiesta struct {
	id    int
	nPers int
	ack   chan int
}

const MAXBUFF = 100
const SUPERBONUS = 0
const ALTRO = 1
const PROPRIETARIO_S = 0
const PROPRIETARIO_A = 1
const AMMINISTRATORE = 2

const NU = 5
const MAXS = 10

var done = make(chan bool)
var termina = make(chan bool)
var terminaGestore = make(chan bool)
var entraSala [3]chan Richiesta
var entraUfficio [2]chan Richiesta
var esciFiliale = make(chan int)

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
	case PROPRIETARIO_S:
		return "proprietario solo"
	case PROPRIETARIO_A:
		return "proprietario accompagnato"
	case AMMINISTRATORE:
		return "amministratore"
	}
	return ""
}
func getTipoFinanziamento(tipo int) string {
	switch tipo {
	case SUPERBONUS:
		return "superbonus"
	case ALTRO:
		return "altro"
	}
	return ""
}
func utente(myid int, tipo int) {
	var tt int
	tt = rand.Intn(5) + 1
	tipoFinanz := rand.Intn(2)
	var nPers int
	if tipo == PROPRIETARIO_S || tipo == AMMINISTRATORE {
		nPers = 1
	} else {
		nPers = 2
	}
	r := Richiesta{id: myid, nPers: nPers, ack: make(chan int)}
	fmt.Printf("[Utente %d %s - %s] INIZIALIZZAZIONE  \n", myid, getTipo(tipo), getTipoFinanziamento(tipoFinanz))
	time.Sleep(time.Duration(tt) * time.Second)

	entraSala[tipo] <- r
	<-r.ack
	//fmt.Printf("[Utente %d %s - %s] Sono nella sala d'aspetto \n", myid, getTipo(tipo), getTipoFinanziamento(tipoFinanz))
	time.Sleep(time.Duration(tt) * time.Second)

	entraUfficio[tipoFinanz] <- r
	<-r.ack
	//fmt.Printf("[Utente %d %s - %s] Sono nell'ufficio \n", myid, getTipo(tipo), getTipoFinanziamento(tipoFinanz))
	time.Sleep(time.Duration(tt) * time.Second)

	esciFiliale <- myid
	//fmt.Printf("[Utente %d %s - %s] Sono uscito dalla filiale \n", myid, getTipo(tipo), getTipoFinanziamento(tipoFinanz))
	time.Sleep(time.Duration(tt) * time.Second)

	fmt.Printf("[Utente %d %s - %s] Termino \n", myid, getTipo(tipo), getTipoFinanziamento(tipoFinanz))
	done <- true
}

func filiale() {

	var nUfficio = 0 //utenti in un ufficio
	var nSala = 0    //utenti in sala d'attesa

	fmt.Printf("[Filiale] INIZIALIZZAZIONE\n")
	for {
		select {
		//SALA D'ASPETTO
		case x := <-whenR(nSala < MAXS, entraSala[AMMINISTRATORE]):
			nSala += x.nPers
			x.ack <- 1
			fmt.Printf("[Filiale] Ammistratore %d entrato in sala d'aspetto, (NS: %d, NU: %d)\n", x.id,nSala,nUfficio)
		case x := <-whenR(nSala < MAXS && len(entraSala[AMMINISTRATORE]) == 0, entraSala[PROPRIETARIO_S]):
			nSala += x.nPers
			x.ack <- 1
			fmt.Printf("[Filiale] Proprietario solo %d entrato in sala d'aspetto, (NS: %d, NU: %d)\n", x.id,nSala,nUfficio)
		case x := <-whenR(nSala < MAXS-1 && len(entraSala[AMMINISTRATORE]) == 0 && len(entraSala[PROPRIETARIO_S]) == 0, entraSala[PROPRIETARIO_A]):
			nSala += x.nPers
			x.ack <- 1
			fmt.Printf("[Filiale] Proprietario accompagnato %d entrato in sala d'aspetto, (NS: %d, NU: %d)\n", x.id,nSala,nUfficio)

			//UFFICIO
		case x := <-whenR(nUfficio < NU, entraUfficio[SUPERBONUS]):
			nSala -= x.nPers
			nUfficio++
			x.ack <- 1
			fmt.Printf("[Filiale] Utente superbonus %d entrato in ufficio, (NS: %d, NU: %d)\n", x.id,nSala,nUfficio)
		case x := <-whenR(nUfficio < NU && len(entraUfficio[SUPERBONUS]) == 0, entraUfficio[ALTRO]):
			nSala -= x.nPers
			nUfficio++
			x.ack <- 1
			fmt.Printf("[Filiale] Utente altro %d entrato in ufficio, (NS: %d, NU: %d)\n", x.id,nSala,nUfficio)

			//USCITA
		case x := <-esciFiliale:
			nUfficio--
			fmt.Printf("[Filiale] Utente %d uscito dalla filiale, (NS: %d, NU: %d)\n", x,nSala,nUfficio)
		case <-terminaGestore:
			fmt.Printf("[Filiale] Termino\n")
			done <- true
			return
		}

	}
}
func main() {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 3; i++ {
		entraSala[i] = make(chan Richiesta, MAXBUFF)
	}
	for i := 0; i < 2; i++ {
		entraUfficio[i] = make(chan Richiesta, MAXBUFF)
	}
	go filiale()
	nUtenti := 15
	for i := 0; i < nUtenti; i++ {
		r := rand.Intn(3)
		go utente(i, r)
	}
	for i := 0; i < nUtenti; i++ {
		<-done
	}
	terminaGestore <- true
	<-done
	fmt.Printf("\n HO FINITO ")
}
