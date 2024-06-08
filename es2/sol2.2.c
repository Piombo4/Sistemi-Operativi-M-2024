#include <pthread.h>
#include <semaphore.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <math.h>
#include <time.h>

#define NUM_THREADS 100
#define MaxP 15
#define MaxA 5


typedef struct{
    int posti_liberi;
    int auto_libere;
    int sospesi[NUM_THREADS];
    int Psosp; //numero totale processi sospesi
    sem_t S; //semaforo condizione
    int ordine_accesso; // quanti gruppi (entrati e non) finora
    sem_t m; //mutua esclusione
} parco;

parco P;

void init(parco *p) {
	int i;
	p->posti_liberi=MaxP;
    p->auto_libere=MaxA;
    p->ordine_accesso=0;
    sem_init(&(p->S),0,0);
    for(i=0; i<NUM_THREADS; i++)
        p->sospesi[i]=0;
    p->Psosp=0;
    sem_init(&(p->m), 0,1); 
}

int davanti(parco *p, int pos) // verifica che vi siano processi in attesa prima di pos-1
{ int i;
    for (i=pos-2; i>=0; i--)
        if (p->sospesi[i]!=0)
            return 1;
    return 0;
}

int da_risvegliare(parco *p) // restituisce il numero di processi sospesi
{   int i, ris=0;
    
    return p->Psosp;
}


void entra(parco *p, int num) {
 
  sem_wait(&p->m);
  
  int myplace;
  myplace= ++p->ordine_accesso;
  
  while(p->posti_liberi<num || p->auto_libere==0 || davanti(p,myplace))
        {   p->sospesi[myplace-1]++;
            p->Psosp++;
            sem_post(&p->m);
            sleep(1); // per evitare che un processo risvegliato torni in coda prima che tutti i risvegliati siano schedulati
            sem_wait(&(p->S)); // sospensione sul semaforo condizione
            sem_wait(&p->m);
            p->sospesi[myplace-1]--;
            p->Psosp--;
        }
        p->posti_liberi-=num;
        p->auto_libere--;
   printf("entra il gruppo %d-simo di %d persone (posti %d, auto %d)\n\n", myplace, num, p->posti_liberi, p->auto_libere);
   sem_post(&p->m);
 }

void esci(parco *p, int num) 
{   int i,k;
    sem_wait(&p->m);
    
    p->posti_liberi+=num;
    p->auto_libere++;
    k=da_risvegliare(p);
    for (i=0; i<k; i++)   
      sem_post(&(p->S)); // risveglio tutti
    
    sem_post (&p->m);
}


void *visitatore(void *t) // gruppo visitatore-> t è il numero di componenti
{	int N;
   
	N = (int)t;
    entra(&P, N); // richiesta entrata di un gruppo di N persone
    sleep(1);
    esci(&P,N);
}


int main (int argc, char *argv[])
{  pthread_t thread[NUM_THREADS];
   int rc, num;
   long t;
   float media, max;
   void *status;
   
   srand(time(NULL));
   init(&P);
  
   for(t=0; t<NUM_THREADS; t++) {
        num=(rand()%5)+1;
        rc = pthread_create(&thread[t], NULL, visitatore, (void *)num); 
        if (rc) {
            printf("ERRORE: %d\n", rc);
            exit(-1);   
        }
    }
	for(t=0; t<NUM_THREADS; t++) {
        rc = pthread_join(thread[t], &status);
        if (rc) 
            printf("ERRORE join thread %ld codice %d\n", t, rc);
   }
   return 0;
}

