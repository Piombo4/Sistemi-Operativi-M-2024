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
const IN = 0
const OUT = 1
const N = 50
const MAXS = 2
const NC = 30

const SINGOLO = 0
const SCOLARESCA = 1
const ADDETTO = 2

var tipoUtente = [3]string{"Visitatore singolo", "Scolaresca", "Addetto"}

var done = make(chan bool)
var termina = make(chan bool)
var terminaMostra = make(chan bool)

var entraCorridoio [3][2]chan Richiesta
var esciCorridoio [3][2]chan Richiesta

/*
* entraCorridoio[][IN] -> entra nel corridoio
* esciCorridoio[][IN] -> esci dal corridoio ed entra nella sala
* entraCorridoio[][OUT] -> esci dalla sala ed entra nel corridoio
* esciCorridoio[][OUT] -> esci dal corridoio
* Nella soluzione ho ipotizzato che uno non potesse entrare nel corridoio se non ci fossero le condizioni per entrare in sala,
* altrimenti si verifica il caso in cui il corridoio è pieno ma le persone non possono entrare e le persone in sala non possono uscire
 */

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

/*
	func getTipo(tipo int) string {
		switch tipo {
		case SPAZZANEVE:
			return "Spazzaneve"
		case SPARGISALE:
			return "Spargisale"
		case CAMION:
			return "Camion rifornitore"
		}
		return ""
	}
*/
func visitatore(myid int, tipo int) {
	var tt int
	tt = rand.Intn(5) + 1
	r := Richiesta{id: myid, ack: make(chan int)}
	fmt.Printf("[Visitatore %d] INIZIALIZZAZIONE \n", myid)
	time.Sleep(time.Duration(tt) * time.Second)

	switch tipo {
	case SINGOLO:
		//ENTRATA
		entraCorridoio[SINGOLO][IN] <- r
		<-r.ack
		//fmt.Printf("[Visitatore singolo %d] Imbocco il corridoio in direzione IN... \n", myid)
		time.Sleep(time.Duration(tt) * time.Second)
		esciCorridoio[SINGOLO][IN] <- r
		<-r.ack
		//fmt.Printf("[Visitatore singolo %d] Sono entrato nella mostra... \n", myid)
		time.Sleep(time.Duration(tt) * time.Second)
		//USCITA
		entraCorridoio[SINGOLO][OUT] <- r
		<-r.ack
		//fmt.Printf("[Visitatore singolo %d] Esco dalla mostra, imbocco corridoio in direzione OUT \n", myid)
		time.Sleep(time.Duration(tt) * time.Second)
		esciCorridoio[SINGOLO][OUT] <- r
		<-r.ack
		//fmt.Printf("[Visitatore singolo %d] Sono uscito dal corridoio \n", myid)
	case SCOLARESCA:
		//ENTRATA
		entraCorridoio[SCOLARESCA][IN] <- r
		<-r.ack
		//fmt.Printf("[Scolaresca %d] Imbocchiamo il corridoio in direzione IN... \n", myid)
		time.Sleep(time.Duration(tt) * time.Second)
		esciCorridoio[SCOLARESCA][IN] <- r
		<-r.ack
		//fmt.Printf("[Scolaresca %d] Siamo entrati nella mostra... \n", myid)
		time.Sleep(time.Duration(tt) * time.Second)
		//USCITA
		entraCorridoio[SCOLARESCA][OUT] <- r
		<-r.ack
		//fmt.Printf("[Scolaresca %d] Usciamo dalla mostra, imbocchiamo corridoio in direzione OUT \n", myid)
		time.Sleep(time.Duration(tt) * time.Second)
		esciCorridoio[SCOLARESCA][OUT] <- r
		<-r.ack
		//fmt.Printf("[Scolaresca %d] Siamo usciti dal corridoio \n", myid)
	}
	fmt.Printf("["+tipoUtente[tipo]+" %d] Finito \n", myid)
	done <- true
}

func addetto(myid int) {
	var tt int
	tt = rand.Intn(5) + 1
	r := Richiesta{id: myid, ack: make(chan int)}
	fmt.Printf("[Addetto %d] INIZIALIZZAZIONE \n", myid)
	for {

		time.Sleep(time.Duration(tt) * time.Second)
		//ENTRATA
		entraCorridoio[ADDETTO][IN] <- r
		<-r.ack
		fmt.Printf("[Addetto %d] Imbocco il corridoio in direzione IN... \n", myid)
		time.Sleep(time.Duration(tt) * time.Second)
		esciCorridoio[ADDETTO][IN] <- r
		<-r.ack
		//fmt.Printf("[Addetto %d] Sono entrato nella mostra... \n", myid)
		time.Sleep(time.Duration(tt) * time.Second)
		//USCITA
		entraCorridoio[ADDETTO][OUT] <- r
		<-r.ack
		//fmt.Printf("[Addetto %d] Esco dalla mostra, imbocco corridoio in direzione OUT \n", myid)
		time.Sleep(time.Duration(tt) * time.Second)
		esciCorridoio[ADDETTO][OUT] <- r
		<-r.ack
		//fmt.Printf("[Addetto %d] Sono uscito dal corridoio \n", myid)
		select {
		case <-termina:
			fmt.Printf("[Addetto %d] Termino\n", myid)
			done <- true
			return
		default:
			continue
		}
	}

}
func debug(nSala int, nPersone []int, nScolaresche []int, nSorveglianti int) {
	fmt.Printf("[DEBUG] CORRIDOIO: (IN,OUT)-(%d,%d), nScolaresche(IN,OUT)-(%d,%d), SALA(V,A)-(%d,%d)\n", nPersone[IN], nPersone[OUT], nScolaresche[IN], nScolaresche[OUT], nSala, nSorveglianti)
}
func mostra() {

	nPersone := [2]int{0, 0}
	nScolaresche := [2]int{0, 0}
	var nSala = 0         //numero visitatori
	var nSorveglianti = 0 //numero sorveglianti in sala
	fmt.Printf("[Mostra] INIZIALIZZAZIONE\n")
	for {
		select {
		//USCITA DALLA SALA
		case x := <-whenR(NC-nPersone[OUT] >= 25 && nPersone[IN] == 0, entraCorridoio[SCOLARESCA][OUT]):
			nPersone[OUT] += 25
			nSala -= 25
			nScolaresche[OUT]++
			fmt.Printf("[Mostra] Scolaresca %d uscita dalla sala imbocca il corridoio in direzione OUT \n", x.id)
			debug(nSala, nPersone[:], nScolaresche[:], nSorveglianti)
			x.ack <- 1

		case x := <-whenR(nPersone[OUT]+nPersone[IN] < NC && nScolaresche[IN] == 0 && len(entraCorridoio[SCOLARESCA][OUT]) == 0, entraCorridoio[SINGOLO][OUT]):
			nPersone[OUT]++
			nSala--
			fmt.Printf("[Mostra] Visitatore singolo %d uscito dalla sala imbocca il corridoio in direzione OUT \n", x.id)
			debug(nSala, nPersone[:], nScolaresche[:], nSorveglianti)
			x.ack <- 1

		case x := <-whenR((nSorveglianti > 1 || nSala == 0) && nPersone[OUT]+nPersone[IN] < NC && nScolaresche[IN] == 0 && len(entraCorridoio[SCOLARESCA][OUT]) == 0 && len(entraCorridoio[SINGOLO][OUT]) == 0, entraCorridoio[ADDETTO][OUT]):
			nPersone[OUT]++
			nSorveglianti--
			fmt.Printf("[Mostra] Addetto uscito %d dalla sala imbocca il corridoio in direzione OUT \n", x.id)
			debug(nSala, nPersone[:], nScolaresche[:], nSorveglianti)
			x.ack <- 1

			//USCITA DAL CORRIDOIO
		case x := <-esciCorridoio[SCOLARESCA][OUT]:
			nPersone[OUT] -= 25
			nScolaresche[OUT]--
			fmt.Printf("[Mostra] Scolaresca %d uscita dal corridoio\n", x.id)
			debug(nSala, nPersone[:], nScolaresche[:], nSorveglianti)
			x.ack <- 1

		case x := <-esciCorridoio[SINGOLO][OUT]:
			nPersone[OUT]--
			fmt.Printf("[Mostra] Visitatore %d singolo uscito dal corridoio \n", x.id)
			debug(nSala, nPersone[:], nScolaresche[:], nSorveglianti)
			x.ack <- 1

		case x := <-esciCorridoio[ADDETTO][OUT]:
			nPersone[OUT]--
			fmt.Printf("[Mostra] Addetto %d uscito dal corridoio \n", x.id)
			debug(nSala, nPersone[:], nScolaresche[:], nSorveglianti)
			x.ack <- 1

			//ENTRATA NEL CORRIDOIO
		case x := <-whenR(nSala+nSorveglianti < N && nSorveglianti < MAXS && nPersone[OUT]+nPersone[IN] < NC && nScolaresche[OUT] == 0 && len(entraCorridoio[SCOLARESCA][OUT]) == 0 && len(entraCorridoio[SINGOLO][OUT]) == 0 && len(entraCorridoio[ADDETTO][OUT]) == 0, entraCorridoio[ADDETTO][IN]):
			nPersone[IN]++
			fmt.Printf("[Mostra] Addetto %d entrato nel corridoio verso la sala \n", x.id)
			debug(nSala, nPersone[:], nScolaresche[:], nSorveglianti)
			x.ack <- 1
		case x := <-whenR(nSorveglianti > 0 && nSala+nSorveglianti < N && nPersone[OUT]+nPersone[IN] < NC && nScolaresche[OUT] == 0 && len(entraCorridoio[SCOLARESCA][OUT]) == 0 && len(entraCorridoio[SINGOLO][OUT]) == 0 && len(entraCorridoio[ADDETTO][OUT]) == 0 && len(entraCorridoio[ADDETTO][IN]) == 0, entraCorridoio[SINGOLO][IN]):
			nPersone[IN]++
			fmt.Printf("[Mostra] Visitatore singolo %d entrato nel corridoio verso la sala \n", x.id)
			debug(nSala, nPersone[:], nScolaresche[:], nSorveglianti)
			x.ack <- 1
		case x := <-whenR(nSorveglianti > 0 && N-(nSala+nSorveglianti) >= 25 && NC-nPersone[IN] >= 25 && nPersone[OUT] == 0 && len(entraCorridoio[SCOLARESCA][OUT]) == 0 && len(entraCorridoio[SINGOLO][OUT]) == 0 && len(entraCorridoio[ADDETTO][OUT]) == 0 && len(entraCorridoio[ADDETTO][IN]) == 0 && len(entraCorridoio[SINGOLO][IN]) == 0, entraCorridoio[SCOLARESCA][IN]):
			nPersone[IN] += 25
			nScolaresche[IN]++
			fmt.Printf("[Mostra] Scolaresca %d entrato nel corridoio verso la sala \n", x.id)
			debug(nSala, nPersone[:], nScolaresche[:], nSorveglianti)
			x.ack <- 1

			//ENTRATA NELLA SALA
		case x := <-esciCorridoio[ADDETTO][IN]:
			nPersone[IN]--
			nSorveglianti++
			fmt.Printf("[Mostra] Addetto %d entrato nella sala \n", x.id)
			debug(nSala, nPersone[:], nScolaresche[:], nSorveglianti)
			x.ack <- 1
		case x := <-esciCorridoio[SINGOLO][IN]:
			nPersone[IN]--
			nSala++
			fmt.Printf("[Mostra] Visitatore singolo %d entrato nella sala \n", x.id)
			debug(nSala, nPersone[:], nScolaresche[:], nSorveglianti)
			x.ack <- 1
		case x := <-esciCorridoio[SCOLARESCA][IN]:
			nPersone[IN] -= 25
			nSala += 25
			nScolaresche[IN]--
			fmt.Printf("[Mostra] Scolaresca %d entrato nella sala \n", x.id)
			debug(nSala, nPersone[:], nScolaresche[:], nSorveglianti)
			x.ack <- 1

		case <-terminaMostra:
			fmt.Printf("[Mostra] Termino\n")
			done <- true
			return
		}

	}
}
func main() {
	rand.Seed(time.Now().Unix())
	nVisitatoriSingoli := 10
	nScolaresche := 3

	for i := 0; i < 3; i++ {
		for j := 0; j < 2; j++ {
			entraCorridoio[i][j] = make(chan Richiesta, MAXBUFF)
			esciCorridoio[i][j] = make(chan Richiesta, MAXBUFF)
		}
	}
	go mostra()
	for i := 0; i < MAXS; i++ {
		go addetto(i)
	}
	for i := 0; i < nVisitatoriSingoli; i++ {
		go visitatore(i, SINGOLO)
	}
	for i := 0; i < nScolaresche; i++ {
		go visitatore(i, SCOLARESCA)
	}

	for i := 0; i < nVisitatoriSingoli+nScolaresche; i++ {
		<-done
	}
	for i := 0; i < MAXS; i++ {
		termina <- true
		<-done
	}

	terminaMostra <- true
	<-done
	fmt.Printf("\n HO FINITO ")
}
