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
const PICCOLA = 0
const GRANDE = 1
const N = 10.0
const M1 = 7
const M2 = 5

var tipologia = [2]string{"Piccola", "Grande"}

var done = make(chan bool)
var termina = make(chan bool)
var terminaGestore = make(chan bool)

var inizioPrelievo [2]chan Richiesta
var finePrelievo = make(chan Richiesta,MAXBUFF)

var inizioManutenzione = make (chan Richiesta,MAXBUFF)
var fineManutenzione = make(chan Richiesta,MAXBUFF)

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

func cittadino(myid int) {
	var tt int
	tt = rand.Intn(2) + 1
	tipo := rand.Intn(2)
	r := Richiesta{id: myid, ack: make(chan int)}
	fmt.Printf("[Cittadino %d] INIZIALIZZAZIONE \n", myid)
	time.Sleep(time.Duration(tt) * time.Second)
	inizioPrelievo[tipo] <- r
	<-r.ack
	//fmt.Printf("[Cittadino %d] Ho iniziato il prelievo per una bottiglia %s\n", myid, tipologia[tipo])
	time.Sleep(time.Duration(tt) * time.Second)
	finePrelievo <- r
	<-r.ack
	//fmt.Printf("[Cittadino %d] Ho finito il prelievo per una bottiglia %s\n", myid, tipologia[tipo])
	fmt.Printf("[Cittadino %d] Termino!\n", myid)
	done <- true
}

func addetto() {
	r := Richiesta{id: 1, ack: make(chan int)}
	fmt.Printf("[Addetto] INIZIALIZZAZIONE \n")
	for {

		time.Sleep(time.Duration( rand.Intn(2) + 1) * time.Second)
		inizioManutenzione <- r
		<-r.ack
		//fmt.Printf("[Addetto] Ho iniziato la manutenzione \n")
		time.Sleep(time.Duration( 10) * time.Second)
		fineManutenzione <- r
		<-r.ack
		//fmt.Printf("[Addetto] Ho finito la manutenzione \n")
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

func debug(nLitri float64, nMonete[] int ){
	fmt.Printf("[DEBUG] nLitri: %f, nMonete[P,G]-[%d,%d]\n",nLitri,nMonete[PICCOLA],nMonete[GRANDE])
}

func erogatore() {

	
	var nLitri = 0.0
	var occupato = false
	var nMonete = [2]int{0,0}

	fmt.Printf("[Erogatore] INIZIALIZZAZIONE\n")
	for {
		select {
			//ADDETTO
		case x:=<-whenR(!occupato && (nLitri==0.0 || nMonete[PICCOLA]==M1 || nMonete[GRANDE]==M2 || (len(inizioPrelievo[PICCOLA])==0 && len(inizioPrelievo[GRANDE])==0)),inizioManutenzione):
			occupato = true
			nLitri = N
			nMonete[GRANDE] = 0
			nMonete[PICCOLA] = 0
			fmt.Printf("[Erogatore] Inizio della manutenzione!\n")
			x.ack<-1
		case x:= <-fineManutenzione:
			occupato = false
			fmt.Printf("[Erogatore] Fine della manutenzione!\n")
			debug(nLitri, nMonete[:] )
			x.ack<-1
			//PRELIEVO
		case x:= <-whenR(!occupato && nLitri >= 0.5 && nMonete[PICCOLA]<M1 &&(nMonete[GRANDE]<M2 || len(inizioManutenzione)==0),inizioPrelievo[PICCOLA]):
			occupato= true
			nLitri-=0.5
			nMonete[PICCOLA]++
			fmt.Printf("[Erogatore] Inizio erogazione al cittadino %d per bottiglia piccola!\n",x.id)
			debug(nLitri, nMonete[:] )
			x.ack<-1
		case x:= <-whenR(!occupato && nLitri >= 1.5 && nMonete[GRANDE]<M2 && len(inizioPrelievo[PICCOLA])==0 &&(nMonete[PICCOLA]<M1 || len(inizioManutenzione)==0),inizioPrelievo[GRANDE]):
			occupato= true
			nLitri-=1.5
			nMonete[GRANDE]++
			fmt.Printf("[Erogatore] Inizio erogazione al cittadino %d per bottiglia grande!\n",x.id)
			debug(nLitri, nMonete[:] )
			x.ack<-1
		case x:=<-finePrelievo:
			occupato = false
			fmt.Printf("[Erogatore] Fine erogazione al cittadino %d!\n",x.id)
			debug(nLitri, nMonete[:] )
			x.ack<-1
		
		case <-terminaGestore:
			fmt.Printf("[Erogatore] Termino\n")
			done <- true
			return
		}

	}
}
func main() {
	rand.Seed(time.Now().Unix())
	for i := 0; i < 2; i++ {
		inizioPrelievo[i] = make(chan Richiesta,MAXBUFF)
	}
	go erogatore()
	go addetto()
	nCittadini:= 15
	for i := 0; i < nCittadini; i++ {
		go cittadino(i)
	}
	for i := 0; i < nCittadini; i++ {
		<-done
	}
	termina <-true
	<-done

	terminaGestore <- true
	<-done
	fmt.Printf("\n HO FINITO ")
}
