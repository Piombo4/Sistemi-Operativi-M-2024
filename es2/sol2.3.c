#include <pthread.h>
#include <semaphore.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <math.h>
#include <time.h>

#define NUM_THREADS 100
#define MaxP 20
#define MaxA 5

#define maxG 5

typedef struct{
int posti_liberi;
int quad_liberi;
int auto_libere;
int sospesi[maxG];
sem_t S[maxG]; //array di semafori privati
int ordine_accesso; // mi serve per verifica
pthread_mutex_t m;
} parco;

parco P;

void init(parco *p) {
	int i;
	p->posti_liberi=MaxP;
    p->auto_libere=MaxA;
    p->ordine_accesso=0;
    
    for(i=0; i<maxG; i++) //inizializzazione semafori e code
    {   sem_init(&(p->S[i]),0,0);
        p->sospesi[i]=0;
    }
    pthread_mutex_init(&p->m, NULL); 
}
int davanti(parco *p, int n) // verifica che vi siano processi in attesa con priorità più alta
{ int i;
    for (i=n; i<maxG; i++)
        if (p->sospesi[i]!=0)
            return 1;
    return 0;
}

int da_risvegliare(parco *p, int i) // restituisce il numero dei gruppi di indice i in attesa 
{   
    return p->sospesi[i];
}

void entra(parco *p, int num) {
 
  pthread_mutex_lock(&p->m);
  
  int myplace;  // mi serve solo per verificare 
  myplace= ++p->ordine_accesso;
  
  
  while(p->posti_liberi<num || p->auto_libere==0 || davanti(p,num))
        {   p->sospesi[num-1]++;
            pthread_mutex_unlock(&p->m);
            sem_wait(&(p->S[num-1])); // sospensione sul semaforo condizione
            pthread_mutex_lock(&p->m);
            p->sospesi[num-1]--;
           
        }
        p->posti_liberi-=num;
        p->auto_libere--;
   }
   printf("entra il gruppo %d-simo di %d persone\n\n", myplace, num);
   pthread_mutex_unlock (&p->m);
 }

void esci(parco *p, int num) 
{   int i,k;
    pthread_mutex_lock(&p->m);
    p->posti_liberi+=num;
    p->auto_libere++;
    
    for (i=maxG-1; i>=0; i--)   // risveglio tutti
        for(k=0; k< da_risvegliare(p,i); k++)
            sem_post(&p->S[i]); 
    pthread_mutex_unlock (&p->m);
}

void *visitatore(void *t) // gruppo visitatore-> t è il numero di componenti
{	int N;
   
	N = (int)t;
    entra(&P, N); // richiesta entrata di un gruppo di N persone
    sleep(rand()%3);
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

