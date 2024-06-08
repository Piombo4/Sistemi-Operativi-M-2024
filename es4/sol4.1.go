// Ponte auto e pedoni project main.go

package main

import (
	"fmt"
	"math/rand"
	"time"
)

const MAXBUFF = 100
const MAXPROC = 100
const MAX = 35 // capacità
const N int = 0
const S int = 1
const PED int = 0
const AUT int = 1

type msg_out struct{ tipo, id int } //tipo parametri uscita
var done = make(chan bool)
var termina = make(chan bool)
var entrataN_A = make(chan int, MAXBUFF) // necessità di accodamento per priorità
var entrataN_P = make(chan int, MAXBUFF) // necessità di accodamento per priorità
var entrataS_A = make(chan int, MAXBUFF) // necessità di accodamento per priorità
var entrataS_P = make(chan int, MAXBUFF) // necessità di accodamento per priorità
var uscitaN = make(chan msg_out)
var uscitaS = make(chan msg_out)
var ACK [MAXPROC]chan int //risposte client

var r int

func when(b bool, c chan int) chan int {
	if !b {
		return nil
	}
	return c
}

func utente(myid int, dir int) {
	var m_out msg_out
	var tt int
	tt = rand.Intn(5) + 1
	tipo := rand.Intn(2) //0 pedone, 1 auto
	fmt.Printf("inizializzazione utente  %d direzione %d di tipo %d in secondi %d \n", myid, dir, tipo, tt)
	//inizializzazione messaggio uscita
	m_out.tipo = tipo //pedone o auto
	m_out.id = myid   //identificatore

	time.Sleep(time.Duration(tt) * time.Second)
	if dir == N { //NORD
		if tipo == PED { // pedone da nord
			entrataN_P <- myid // send asincrona
			<-ACK[myid]        // attesa x sincronizzazione
			fmt.Printf("[pedone %d]\t entrato da   NORD\n", myid)
			tt = rand.Intn(5)
			time.Sleep(time.Duration(tt) * time.Second)
			uscitaN <- m_out
			//fmt.Printf("[pedone %d]\t proveniente da  NORD è uscito\n", myid)
		} else { // auto da nord
			entrataN_A <- myid // send asincrona
			<-ACK[myid]        // attesa x sincronizzazione
			fmt.Printf("[auto %d]\t entrata da  NORD\n", myid)
			tt = rand.Intn(5)
			time.Sleep(time.Duration(tt) * time.Second)
			uscitaN <- m_out
			//fmt.Printf("[auto %d]\t proveniente da  NORD è uscita\n", myid)
		}

	} else { // SUD
		if tipo == PED { // pedone da sud
			entrataS_P <- myid // send asincrona
			<-ACK[myid]        // attesa x sincronizzazione
			fmt.Printf("[pedone %d]\t entrato da  SUD\n", myid)
			tt = rand.Intn(5)
			time.Sleep(time.Duration(tt) * time.Second)
			uscitaS <- m_out
			//fmt.Printf("[pedone %d]\t proveniente da  SUD è uscito\n", myid)
		} else { // auto da sud
			entrataS_A <- myid // send asincrona
			<-ACK[myid]        // attesa x sincronizzazione
			fmt.Printf("[auto %d]\t entrata  da  SUD\n", myid)
			tt = rand.Intn(5)
			time.Sleep(time.Duration(tt) * time.Second)
			uscitaS <- m_out
			//fmt.Printf("[auto %d]\t 30proveniente da  NORD è uscita\n", myid)
		}
	}
	done <- true
}

func server() {

	var contN [2]int //pedoni e auto in dir nord
	var contS [2]int //pedoni e auto in dir sud
	var tot int

	for {

		select {
		case x := <-when((tot < MAX) && (contN[AUT] == 0), entrataS_P): //pedone da sud
			contS[PED]++
			tot++
			ACK[x] <- 1 // termine "call"
		case x := <-when((tot < MAX) && (contS[AUT] == 0) && (len(entrataS_P) == 0), entrataN_P): //pedone da nord
			contN[PED]++
			tot++
			ACK[x] <- 1 // termine "call"
		case x := <-when((tot+10 <= MAX) && (contN[PED]+contN[AUT] == 0) && (len(entrataN_P)+len(entrataS_P) == 0), entrataS_A): //auto da sud
			contS[AUT]++
			tot = tot + 10 //un'auto vale 10 persone
			ACK[x] <- 1    // termine "call"
		case x := <-when((tot+10 <= MAX) && (contS[PED]+contS[AUT] == 0) && (len(entrataN_P)+len(entrataS_P)+len(entrataS_A) == 0), entrataN_A): //auto da nord
			contN[AUT]++
			tot = tot + 10 //un'auto vale 10 persone
			ACK[x] <- 1    // termine "call"
		case x := <-uscitaN:
			contN[x.tipo]--
			if x.tipo == PED {
				tot--
			} else {
				tot = tot - 10 //un'auto vale 10 persone
			}
		case x := <-uscitaS:
			contS[x.tipo]--
			if x.tipo == PED {
				tot--
			} else {
				tot = tot - 10 //un'auto vale 10 persone
			}
		case <-termina: // quando tutti i processi hanno finito
			fmt.Println("IL PONTE CHIUDE !")
			done <- true
			return
		}

	}
}

func main() {
	var VN int
	var VS int
	var i int

	fmt.Printf("\n quanti veicoli NORD (max %d)? ", MAXPROC/2)
	fmt.Scanf("%d", &VN)
	fmt.Printf("\n quanti veicoli SUD (max %d)? ", MAXPROC/2)
	fmt.Scanf("%d", &VS)

	//inizializzazione canali
	for i = 0; i < VN+VS; i++ {
		ACK[i] = make(chan int, MAXBUFF)
	}

	rand.Seed(time.Now().Unix())
	go server()

	for i = 0; i < VS; i++ {
		go utente(i, S)
	}
	for j := i; j < VN+VS; j++ {
		go utente(j, N)
	}

	for i := 0; i < VN+VS; i++ {
		<-done
	}
	termina <- true
	<-done
	fmt.Printf("\n HO FINITO ")
}
