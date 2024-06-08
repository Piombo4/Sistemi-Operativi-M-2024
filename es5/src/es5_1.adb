with Ada.Text_IO, Ada.Integer_Text_IO;
use Ada.Text_IO, Ada.Integer_Text_IO;

procedure es5_1 is
   
   type utente_ID is range 1..10;
   type servizio_ID is  (EVE, TUR);
   
   task type utente (ID: utente_ID; TIPO:servizio_ID);		-- task che rappresenta il generico cliente


   type ac is access utente;
   
   task type ufficio is
      entry acquisisci_sportello(servizio_ID) (ID:utente_ID);
      entry rilascia_sportello(servizio_ID) (ID:utente_ID);   
   end ufficio;
   
   U:ufficio;
   
   task body ufficio is
      MAX  : constant INTEGER := 5; --N sportelli
      al_lavoro: array(servizio_ID'Range) of Integer; 
   begin
      Put_Line ("SERVER iniziato!");
      for i in al_lavoro'Range loop
         al_lavoro(i):=0;
      end loop;
      delay 2.0;
      loop
         select
            when al_lavoro(EVE) + al_lavoro(TUR) < MAX  =>
               accept acquisisci_sportello(TUR) (ID : in utente_ID) do
                  Put_Line("acquisisco sportello TUR "& utente_ID'Image(ID) &" !");
                  al_lavoro(TUR) := al_lavoro(TUR) +1;
               end;
         or
            when al_lavoro(EVE) + al_lavoro(TUR) < MAX and acquisisci_sportello(TUR)'COUNT=0 =>
               accept acquisisci_sportello(EVE) (ID : in utente_ID) do
                  Put_Line("acquisisco sportello EVE "& utente_ID'Image(ID) &" !");
                  al_lavoro(EVE) := al_lavoro(EVE) +1;
               end;
         or           
            accept rilascia_sportello(EVE) (ID : in utente_ID) do
               Put_Line("rilascio sportello EVE "& utente_ID'Image(ID) &" !");
               al_lavoro(EVE) := al_lavoro(EVE) -1;
            end;
         or           
            accept rilascia_sportello(TUR) (ID : in utente_ID) do
               Put_Line("rilascio sportello TUR "& utente_ID'Image(ID) &" !");
               al_lavoro(TUR) := al_lavoro(TUR) -1;
            end;
         end select;
      end loop;
   end;
   
   task body utente is
   begin
      Put_Line ("utente " & utente_ID'Image (ID) &" iniziato!");
      U.acquisisci_sportello(TIPO) (ID);
      delay 1.0;
      U.rilascia_sportello(TIPO) (ID);
   end;
   
   New_utente: ac;
   tipo: servizio_ID;

begin -- equivale al main
   
   
   for I in utente_ID'Range loop  -- ciclo creazione task
      if I mod 2 = 0 then
         tipo := TUR;
      else
         tipo := EVE;
      end if;
         New_utente := new utente(I,tipo); -- creazione cliente I-simo
      end loop;
   end;
