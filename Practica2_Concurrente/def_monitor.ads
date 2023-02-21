--AUTORS: DIEGO BERMEJO CABAÑAS I MARC CAÑELLAS GOMEZ
--VIDEOS EXPLICATIU: https://www.youtube.com/watch?v=6LDuejoOCYU

with Ada.Strings.Unbounded;    use Ada.Strings.Unbounded;

package def_monitor is

   type capacidades is array (1..3) of integer;                           --Array amb les capacitats de les sales
   type tipos is array (1..3) of integer;                                  --TIPO 1: LIBRE ; TIPO 2: FUMADOR ; TIPO 3: NO FUMADOR
   
   protected type maitre is

      entry demanaTaulaF(Nom: in Unbounded_String; Taula: out Integer) ;             --Introduim un fumador al restaurant      
      entry demanaTaulaNoF(Nom: in Unbounded_String; Taula: out Integer);            --Introduim un nofumador al restaurant
      procedure leave(Nom: in Unbounded_String; sala: in integer; fuma : in Boolean);--El comensal surt del restaurant
      function disponible(tipo: in Integer) return Integer;                          --Mira si el comensal té puesto a una taula


   private

      arrCap : capacidades := (3,3,3);--Inicialitzam les capacitats de les sales a 3
      arrType : tipos := (1,1,1);     --Inicialitzam els tipus de les sales a lliure

     --VARIABLES PARA COMPROBAR EN EL ENTRY
      --salfum : integer := 0;   --Numero de sales de tipus fumadors que tenguin algun puesto lliure
      --salLibre :  integer := 3; --Numero de sales lliures
      --salnfum :  integer := 0; --Numero de sales de tipus no fumadors que tenguin algun puesto lliure
      
   end maitre;

end def_monitor;