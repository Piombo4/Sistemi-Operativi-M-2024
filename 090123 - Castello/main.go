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
const NS = 5
const NM = 4
const SALITA = 0
const DISCESA = 1
const AUTO = 0
const CAMPER = 1
const SPAZZANEVE = 2
const STANDARD = 0
const MAXI = 1

var done = make(chan bool)
var termina = make(chan bool)
var terminaGestore = make(chan bool)

var inizioSalita [3]chan Richiesta
var fineSalita [3]chan Richiesta
var inizioDiscesa [3]chan Richiesta
var fineDiscesa [3]chan Richiesta

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
	case SPAZZANEVE:
		return "Spazzaneve"
	case CAMPER:
		return "Camper"
	case AUTO:
		return "AUTO"
	}
	return ""
}
func turista(myid int, tipo int) {
	var tt int
	tt = rand.Intn(5) + 1
	r := Richiesta{id: myid, ack: make(chan int)}
	fmt.Printf("[TURISTA %d "+getTipo(tipo)+"] INIZIALIZZAZIONE \n", myid)
	time.Sleep(time.Duration(tt) * time.Second)

	inizioSalita[tipo] <- r
	<-r.ack
	fmt.Printf("[TURISTA %d "+getTipo(tipo)+"] Ho iniziato la salita... \n", myid)
	time.Sleep(time.Duration(tt) * time.Second)
	fineSalita[tipo] <- r
	<-r.ack
	fmt.Printf("[TURISTA %d "+getTipo(tipo)+"] Ho finito la salita... Visito il castello! \n", myid)
	time.Sleep(time.Duration(tt) * time.Second)

	inizioDiscesa[tipo] <- r
	<-r.ack
	fmt.Printf("[TURISTA %d "+getTipo(tipo)+"] Ho iniziato la discesa... \n", myid)
	time.Sleep(time.Duration(tt) * time.Second)
	fineDiscesa[tipo] <- r
	<-r.ack
	fmt.Printf("[TURISTA %d "+getTipo(tipo)+"] Ho finito la discesa! \n", myid)
	time.Sleep(time.Duration(tt) * time.Second)

	fmt.Printf("[TURISTA %d "+getTipo(tipo)+"] Termino\n", myid)
	done <- true
}

func spazzaneve() {
	var tt int
	tt = rand.Intn(5) + 1
	r := Richiesta{id: 1, ack: make(chan int)}
	fmt.Printf("[SPAZZANEVE] INIZIALIZZAZIONE \n")
	for {
		tt = rand.Intn(5) + 1
		time.Sleep(time.Duration(tt) * time.Second)
		inizioDiscesa[SPAZZANEVE] <- r
		<-r.ack
		fmt.Printf("[SPAZZANEVE] Ho iniziato la discesa... \n")
		time.Sleep(time.Duration(tt) * time.Second)
		fineDiscesa[SPAZZANEVE] <- r
		<-r.ack
		fmt.Printf("[SPAZZANEVE] Ho finito la discesa... Sono al barr \n")
		time.Sleep(time.Duration(tt) * time.Second)
		inizioSalita[SPAZZANEVE] <- r
		<-r.ack
		fmt.Printf("[SPAZZANEVE] Ho iniziato la salita... \n")
		time.Sleep(time.Duration(tt) * time.Second)
		fineSalita[SPAZZANEVE] <- r
		<-r.ack
		fmt.Printf("[SPAZZANEVE] Ho finito la salita... Sosto nel piazzale! \n")

		select {
		case <-termina:
			fmt.Printf("[SPAZZANEVE] Termino\n")
			done <- true
			return
		default:
			continue
		}
	}

}

func checkAutoId(id int, autoInPostiMaxi []int) int {
	for i := 0; i < NM; i++ {
		if autoInPostiMaxi[i] == id {
			return i
		}

	}
	return -1
}
func castello() {

	postiStandard := NS
	postiMaxi := NM
	var nAuto [2]int
	var nCamper [2]int
	var spazzaneveMovente = false
	var autoInPostiMaxi [NM]int //Array che uso per controllare quali auto si sono messe in un posto maxi
	for i := 0; i < NM; i++ {
		autoInPostiMaxi[i] = -1
	}
	fmt.Printf("[CASTELLO] INIZIALIZZAZIONE\n")
	for {
		select {
		//DISCESA SPAZZANEVE
		case x := <-whenR(nAuto[SALITA]+nCamper[SALITA]+nAuto[DISCESA]+nCamper[DISCESA] == 0, inizioDiscesa[SPAZZANEVE]):
			spazzaneveMovente = true
			fmt.Printf("[CASTELLO] Spazzaneve ha iniziato la discesa...\n")
			x.ack <- 1
		case x := <-fineDiscesa[SPAZZANEVE]:
			spazzaneveMovente = false
			fmt.Printf("[CASTELLO] Spazzaneve ha finito la discesa!\n")
			x.ack <- 1
			//DISCESA CAMPER
		case x := <-whenR(!spazzaneveMovente && nAuto[SALITA] == 0 && nCamper[SALITA] == 0 && len(inizioDiscesa[SPAZZANEVE]) == 0, inizioDiscesa[CAMPER]):
			postiMaxi++
			nCamper[DISCESA]++
			fmt.Printf("[CASTELLO] Camper ha iniziato la discesa...\n")
			x.ack <- 1
		case x := <-fineDiscesa[CAMPER]:
			nCamper[DISCESA]--
			fmt.Printf("[CASTELLO] Camper ha finito la discesa!\n")
			x.ack <- 1
			//DISCESA AUTO
		case x := <-whenR(!spazzaneveMovente && nCamper[SALITA] == 0 && len(inizioDiscesa[CAMPER]) == 0 && len(inizioDiscesa[SPAZZANEVE]) == 0, inizioDiscesa[AUTO]):
			var i = checkAutoId(x.id, autoInPostiMaxi[:])
			if i >= 0 {
				postiMaxi++
				autoInPostiMaxi[i] = -1

			} else {
				postiStandard++
			}
			nAuto[DISCESA]++
			fmt.Printf("[CASTELLO] Auto ha iniziato la discesa...\n")
			x.ack <- 1
		case x := <-fineDiscesa[AUTO]:
			nAuto[DISCESA]--
			fmt.Printf("[CASTELLO] Auto ha finito la discesa!\n")
			x.ack <- 1
			//SALITA CAMPER
		case x := <-whenR(postiMaxi > 0 && !spazzaneveMovente && nAuto[DISCESA] == 0 && nCamper[DISCESA] == 0 && len(inizioDiscesa[CAMPER]) == 0 && len(inizioDiscesa[AUTO]) == 0 && len(inizioDiscesa[SPAZZANEVE]) == 0, inizioSalita[CAMPER]):
			nCamper[SALITA]++
			postiMaxi--
			fmt.Printf("[CASTELLO] Camper ha iniziato la salita...\n")
			x.ack <- 1
		case x := <-fineSalita[CAMPER]:
			nCamper[SALITA]--
			fmt.Printf("[CASTELLO] Camper ha finito la salita!\n")
			x.ack <- 1
			//SALITA AUTO
		case x := <-whenR((postiStandard > 0 || postiMaxi > 0) && !spazzaneveMovente && len(inizioDiscesa[CAMPER]) == 0 && len(inizioDiscesa[AUTO]) == 0 && len(inizioDiscesa[SPAZZANEVE]) == 0 && len(inizioSalita[CAMPER]) == 0, inizioSalita[AUTO]):
			nAuto[SALITA]++
			//Se non ci sono posti standard va in un posto maxi.
			//A questo proposito va memorizzato l'id dell'auto che ha preso il posto maxi, così quando sarà necessario so qual è la macchina che lo deve liberare
			if postiStandard > 0 {
				postiStandard--
			} else {
				autoInPostiMaxi[NM-postiMaxi] = x.id
				postiMaxi--
			}
			fmt.Printf("[CASTELLO] Auto ha iniziato la salita...\n")
			x.ack <- 1
		case x := <-fineSalita[AUTO]:
			nAuto[SALITA]--
			fmt.Printf("[CASTELLO] Auto ha finito la salita!\n")
			x.ack <- 1
			//SALITA SPAZZANEVE
		case x := <-whenR(nAuto[SALITA]+nCamper[SALITA]+nAuto[DISCESA]+nCamper[DISCESA] == 0 && len(inizioDiscesa[CAMPER]) == 0 && len(inizioDiscesa[AUTO]) == 0 && len(inizioSalita[AUTO]) == 0 && len(inizioSalita[CAMPER]) == 0, inizioSalita[SPAZZANEVE]):
			spazzaneveMovente = true
			fmt.Printf("[CASTELLO] Spazzaneve ha iniziato la salita...\n")
			x.ack <- 1
		case x := <-fineSalita[SPAZZANEVE]:
			spazzaneveMovente = false
			fmt.Printf("[CASTELLO] Spazzaneve ha finito la salita!\n")
			x.ack <- 1
		case <-terminaGestore:
			fmt.Printf("[CASTELLO] Termino\n")
			done <- true
			return
		}

	}
}
func main() {
	rand.Seed(time.Now().Unix())
	for i := 0; i < 3; i++ {
		inizioSalita[i] = make(chan Richiesta, MAXBUFF)
		fineSalita[i] = make(chan Richiesta, MAXBUFF)
		fineDiscesa[i] = make(chan Richiesta, MAXBUFF)
		inizioDiscesa[i] = make(chan Richiesta, MAXBUFF)
	}
	go castello()
	go spazzaneve()
	nVeicoli := 30
	for i := 0; i < nVeicoli; i++ {
		tipo := rand.Intn(2)
		go turista(i, tipo)
	}
	for i := 0; i < nVeicoli; i++ {
		<-done
	}
	//SPAZZANEVE
	termina <- true
	<-done

	terminaGestore <- true
	<-done
	fmt.Printf("\n HO FINITO ")
}
