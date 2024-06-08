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
const FFP2 = 0
const CHIRURGICHE = 1
const MISTO = 2

const NC = 10
const NF = 10
const LC = 3
const LF = 3
const LM = 2

const nAR = 5

var done = make(chan bool)
var terminaAR = make(chan bool)
var terminaF = make(chan bool)
var terminaGestore = make(chan bool)
var iniziaPrelievo [3]chan Richiesta
var finePrelievo [3]chan int
var iniziaRifornimento [2]chan Richiesta
var fineRifornimento [2]chan int

var tipoMask = [3]string{"FFP2", "CHIRURGICA", "MISTO"}

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

/*func getTipo(tipo int) string {
	switch tipo {
	case SPAZZANEVE:
		return "Spazzaneve"
	case SPARGISALE:
		return "Spargisale"
	case CAMION:
		return "Camion rifornitore"
	}
	return ""
}*/

func ar(myid int) {
	r := Richiesta{id: myid, ack: make(chan int)}
	fmt.Printf("[AR %d] INIZIALIZZAZIONE \n", myid)
	time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)
	for {
		tipo := rand.Intn(3)
		iniziaPrelievo[tipo] <- r
		<-r.ack
		fmt.Printf("[AR %d] Sto prelevando, impiegherò 3 secondi \n", myid)
		time.Sleep(time.Duration(3) * time.Second)
		finePrelievo[tipo] <- myid
		switch tipo {
		case FFP2:
			fmt.Printf("[AR %d] Ho prelevato %d FFP2\n", myid, LF)
		case CHIRURGICHE:
			fmt.Printf("[AR %d] Ho prelevato %d chirurgiche FFP2\n", myid, LC)
		case MISTO:
			fmt.Printf("[AR %d] Ho prelevato %d miste FFP2\n", myid, LM)
		}
		time.Sleep(time.Duration(5) * time.Second)
		select {
		case <-terminaAR:
			fmt.Printf("[AR %d] Termino\n", myid)
			done <- true
			return
		default:
			continue
		}
	}
}

func fornitore(tipo int) {
	r := Richiesta{id: tipo, ack: make(chan int)}
	fmt.Printf("[FORNITORE %d] INIZIALIZZAZIONE \n", tipo)
	time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)
	for {

		iniziaRifornimento[tipo] <- r
		<-r.ack
		fmt.Printf("[FORNITORE " + tipoMask[tipo] + "] Sto rifornendo , impiegherò 5 secondi \n")
		time.Sleep(time.Duration(5) * time.Second)
		fineRifornimento[tipo] <- tipo
		fmt.Printf("[FORNITORE " + tipoMask[tipo] + "] Ho rifornito! \n")
		time.Sleep(time.Duration(5) * time.Second)
		select {
		case <-terminaF:
			fmt.Printf("[FORNITORE " + tipoMask[tipo] + "] Termino\n")
			done <- true
			return
		default:
			continue
		}
	}
}
func magazzino() {
	var nScatole [2]int
	var prelevandoDaScaffale [2]int
	var rifornendoScaffale [2]int
	for i := 0; i < 2; i++ {
		nScatole[i] = 0
		prelevandoDaScaffale[i] = 0
		rifornendoScaffale[i] = 0
	}
	fmt.Printf("[MAGAZZINO] INIZIALIZZAZIONE\n")
	for {
		select {
		//AR MISTO
		case x := <-whenR(rifornendoScaffale[CHIRURGICHE] == 0 && rifornendoScaffale[FFP2]== 0 && nScatole[CHIRURGICHE] >= LM && nScatole[FFP2] >= LM, iniziaPrelievo[MISTO]):
			prelevandoDaScaffale[CHIRURGICHE]++
			prelevandoDaScaffale[FFP2]++
			nScatole[CHIRURGICHE] -= LM
			nScatole[FFP2] -= LM
			fmt.Printf("[MAGAZZINO] Iniziato il prelievo di un lotto misto da parte di AR %d,TOT:[C:%d,F:%d]\n", x.id, nScatole[CHIRURGICHE], nScatole[FFP2])
			fmt.Printf("[MAGAZZINO] Waiting: prelievo[F,C,M]:[%d,%d,%d] - rifornimento[F,C]:[%d,%d]\n", len(iniziaPrelievo[FFP2]), len(iniziaPrelievo[CHIRURGICHE]), len(iniziaPrelievo[MISTO]), len(iniziaRifornimento[FFP2]), len(iniziaPrelievo[CHIRURGICHE]))

			x.ack <- 1
		case x := <-finePrelievo[MISTO]:
			prelevandoDaScaffale[CHIRURGICHE]--
			prelevandoDaScaffale[FFP2]--
			fmt.Printf("[MAGAZZINO] Terminato il prelievo di un lotto misto da parte di AR %d,TOT:[C:%d,F:%d]\n", x, nScatole[CHIRURGICHE], nScatole[FFP2])
			fmt.Printf("[MAGAZZINO] Waiting: prelievo[F,C,M]:[%d,%d,%d] - rifornimento[F,C]:[%d,%d]\n", len(iniziaPrelievo[FFP2]), len(iniziaPrelievo[CHIRURGICHE]), len(iniziaPrelievo[MISTO]), len(iniziaRifornimento[FFP2]), len(iniziaPrelievo[CHIRURGICHE]))

			//AR FFP2
		case x := <-whenR(rifornendoScaffale[FFP2]==0 && nScatole[FFP2] >= LF && len(iniziaPrelievo[MISTO]) == 0, iniziaPrelievo[FFP2]):
			prelevandoDaScaffale[FFP2]++
			nScatole[FFP2] -= LF
			fmt.Printf("[MAGAZZINO] Iniziato il prelievo di un lotto di FFP2 da parte di AR %d,TOT:[C:%d,F:%d]\n", x.id, nScatole[CHIRURGICHE], nScatole[FFP2])
			fmt.Printf("[MAGAZZINO] Waiting: prelievo[F,C,M]:[%d,%d,%d] - rifornimento[F,C]:[%d,%d]\n", len(iniziaPrelievo[FFP2]), len(iniziaPrelievo[CHIRURGICHE]), len(iniziaPrelievo[MISTO]), len(iniziaRifornimento[FFP2]), len(iniziaPrelievo[CHIRURGICHE]))

			x.ack <- 1
		case x := <-finePrelievo[FFP2]:
			prelevandoDaScaffale[FFP2]--
			fmt.Printf("[MAGAZZINO] Terminato il prelievo di un lotto di ffp2 da parte di AR %d,TOT:[C:%d,F:%d]\n", x, nScatole[CHIRURGICHE], nScatole[FFP2])
			fmt.Printf("[MAGAZZINO] Waiting: prelievo[F,C,M]:[%d,%d,%d] - rifornimento[F,C]:[%d,%d]\n", len(iniziaPrelievo[FFP2]), len(iniziaPrelievo[CHIRURGICHE]), len(iniziaPrelievo[MISTO]), len(iniziaRifornimento[FFP2]), len(iniziaPrelievo[CHIRURGICHE]))

			//AR CHIRURGICHE
		case x := <-whenR(rifornendoScaffale[CHIRURGICHE]==0 && nScatole[CHIRURGICHE] >= LC && len(iniziaPrelievo[MISTO]) == 0 && len(iniziaPrelievo[FFP2]) == 0, iniziaPrelievo[CHIRURGICHE]):
			prelevandoDaScaffale[CHIRURGICHE]++
			nScatole[CHIRURGICHE] -= LC
			fmt.Printf("[MAGAZZINO] Iniziato il prelievo di un lotto di chirurgiche da parte di AR %d,TOT:[C:%d,F:%d]\n", x.id, nScatole[CHIRURGICHE], nScatole[FFP2])
			fmt.Printf("[MAGAZZINO] Waiting: prelievo[F,C,M]:[%d,%d,%d] - rifornimento[F,C]:[%d,%d]\n", len(iniziaPrelievo[FFP2]), len(iniziaPrelievo[CHIRURGICHE]), len(iniziaPrelievo[MISTO]), len(iniziaRifornimento[FFP2]), len(iniziaPrelievo[CHIRURGICHE]))

			x.ack <- 1
		case x := <-finePrelievo[CHIRURGICHE]:
			prelevandoDaScaffale[CHIRURGICHE]--
			fmt.Printf("[MAGAZZINO] Terminato il prelievo di un lotto di chirurgiche da parte di AR %d,TOT:[C:%d,F:%d]\n", x, nScatole[CHIRURGICHE], nScatole[FFP2])
			fmt.Printf("[MAGAZZINO] Waiting: prelievo[F,C,M]:[%d,%d,%d] - rifornimento[F,C]:[%d,%d]\n", len(iniziaPrelievo[FFP2]), len(iniziaPrelievo[CHIRURGICHE]), len(iniziaPrelievo[MISTO]), len(iniziaRifornimento[FFP2]), len(iniziaPrelievo[CHIRURGICHE]))

			//FORNITORI
		case x := <-whenR(prelevandoDaScaffale[CHIRURGICHE]==0 && (nScatole[CHIRURGICHE] < nScatole[FFP2] || len(iniziaRifornimento[FFP2]) == 0), iniziaRifornimento[CHIRURGICHE]):
			nScatole[CHIRURGICHE] = NC
			rifornendoScaffale[CHIRURGICHE]++
			fmt.Printf("[MAGAZZINO] Iniziato il rifornimento delle chirurgiche da parte di Fornitore %d,TOT:[C:%d,F:%d]\n", x.id, nScatole[CHIRURGICHE], nScatole[FFP2])
			fmt.Printf("[MAGAZZINO] Waiting: prelievo[F,C,M]:[%d,%d,%d] - rifornimento[F,C]:[%d,%d]\n", len(iniziaPrelievo[FFP2]), len(iniziaPrelievo[CHIRURGICHE]), len(iniziaPrelievo[MISTO]), len(iniziaRifornimento[FFP2]), len(iniziaPrelievo[CHIRURGICHE]))

			x.ack <- 1
		case x := <-whenR(prelevandoDaScaffale[FFP2]==0 && (nScatole[FFP2] <= nScatole[CHIRURGICHE] || len(iniziaRifornimento[CHIRURGICHE]) == 0), iniziaRifornimento[FFP2]):
			nScatole[FFP2] = NF
			rifornendoScaffale[FFP2]++
			fmt.Printf("[MAGAZZINO] Iniziato il rifornimento delle FFP2 da parte di Fornitore %d,TOT:[C:%d,F:%d]\n", x.id, nScatole[CHIRURGICHE], nScatole[FFP2])
			fmt.Printf("[MAGAZZINO] Waiting: prelievo[F,C,M]:[%d,%d,%d] - rifornimento[F,C]:[%d,%d]\n", len(iniziaPrelievo[FFP2]), len(iniziaPrelievo[CHIRURGICHE]), len(iniziaPrelievo[MISTO]), len(iniziaRifornimento[FFP2]), len(iniziaPrelievo[CHIRURGICHE]))

			x.ack <- 1
		case x := <-fineRifornimento[CHIRURGICHE]:
			rifornendoScaffale[CHIRURGICHE]--
			fmt.Printf("[MAGAZZINO] Terminato il rifornimento delle chirurgiche da parte di Fornitore %d,TOT:[C:%d,F:%d]\n", x, nScatole[CHIRURGICHE], nScatole[FFP2])
			fmt.Printf("[MAGAZZINO] Waiting: prelievo[F,C,M]:[%d,%d,%d] - rifornimento[F,C]:[%d,%d]\n", len(iniziaPrelievo[FFP2]), len(iniziaPrelievo[CHIRURGICHE]), len(iniziaPrelievo[MISTO]), len(iniziaRifornimento[FFP2]), len(iniziaPrelievo[CHIRURGICHE]))

		case x := <-fineRifornimento[FFP2]:
			rifornendoScaffale[FFP2]--
			fmt.Printf("[MAGAZZINO] Terminato il rifornimento delle FFP2 da parte di Fornitore %d,TOT:[C:%d,F:%d]\n", x, nScatole[CHIRURGICHE], nScatole[FFP2])
			fmt.Printf("[MAGAZZINO] Waiting: prelievo[F,C,M]:[%d,%d,%d] - rifornimento[F,C]:[%d,%d]\n", len(iniziaPrelievo[FFP2]), len(iniziaPrelievo[CHIRURGICHE]), len(iniziaPrelievo[MISTO]), len(iniziaRifornimento[FFP2]), len(iniziaPrelievo[CHIRURGICHE]))

		case <-terminaGestore:
			fmt.Printf("[MAGAZZINO] Termino\n")
			done <- true
			return
		}

	}
}
func main() {
	rand.Seed(time.Now().Unix())
	for i := 0; i < 3; i++ {
		iniziaPrelievo[i] = make(chan Richiesta, MAXBUFF)
		finePrelievo[i] = make(chan int)
	}
	for i := 0; i < 2; i++ {
		iniziaRifornimento[i] = make(chan Richiesta, MAXBUFF)
		fineRifornimento[i] = make(chan int)
	}
	go magazzino()
	go fornitore(FFP2)
	go fornitore(CHIRURGICHE)
	for i := 0; i < nAR; i++ {
		go ar(i)
	}

	time.Sleep(time.Duration(20) * time.Second)
	fmt.Printf("\n\n START ELIMINATION \n\n")
	for i := 0; i < nAR; i++ {
		terminaAR <- true
		<-done
	}
	fmt.Printf("\n\n AR DEAD \n\n")

	terminaF <- true
	<-done
	terminaF <- true
	<-done
	terminaGestore <- true
	<-done
	fmt.Printf("\n HO FINITO ")
}
