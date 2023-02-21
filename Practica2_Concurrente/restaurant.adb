--AUTORS: DIEGO BERMEJO CABAÑAS I MARC CAÑELLAS GOMEZ
--VIDEOS EXPLICATIU: https://www.youtube.com/watch?v=6LDuejoOCYU

with Text_IO; use  Text_IO;
with Ada.Strings.Unbounded;    use Ada.Strings.Unbounded;
with Ada.Text_IO.Unbounded_IO; use Ada.Text_IO.Unbounded_IO;
with Ada.Task_Identification;  use Ada.Task_Identification;
with Ada.Float_Text_IO;        use Ada.Float_Text_IO;
with def_monitor; use def_monitor;
with Ada.Numerics.Discrete_Random;

procedure restaurant is

  THREADS : constant integer := 14;
  Names_file : constant String := "personas.txt";
  MIN_RANDOM: constant integer := 1; -- Minimo valor aleatorio
  MAX_RANDOM: constant integer := 6; -- Maximo valor aleatorio 

  -- Random number generator
   -- Use: Random(G)
   subtype RANDOM_RANGE is integer range MIN_RANDOM .. MAX_RANDOM;

   package R is new
      Ada.Numerics.Discrete_Random (RANDOM_RANGE);
   use R;
   G : Generator;

  -----
  -- Dada compartida
  -----
  monitor : maitre;

  -----
  -- Especificacio de la tasca
  -----
  task type comensal_task is
    entry Start (name: in Unbounded_String; fuma: in Boolean);
  end comensal_task;

  -----
  -- Cos de la tasca
  -----
  task body comensal_task is
    Taula: Integer;
    nom : Unbounded_String;
    fumador : Boolean;
    espera: integer;
     
  begin
    accept Start (name: in Unbounded_String; fuma: in Boolean) do
      nom := name;
      fumador := fuma;
    end Start;
      espera:=Random(G);
      if (fumador) then
        Put_Line("BON DIA som en "& nom& " i som fumador ");
        monitor.demanaTaulaF(nom,Taula);
        Put_Line("En " & Nom & " diu: Prendré el menú del dia. Som al saló "& Taula'img);
        delay(Duration(espera));
        Put_Line("En " &nom& " diu: Ja he dinat, el compte per favor");
        monitor.leave(nom,Taula,fumador);
        Put_Line("En " &nom&" S'EN VA");
        
      else
        Put_Line("          BON DIA som en  "& nom& " i som NO fumador ");
        monitor.demanaTaulaNoF(nom,Taula);
        Put_Line("     En " & Nom & " diu: Prendré el menú del dia. Som al saló "& Taula'img);
        delay(Duration(espera));
        Put_Line("          En " &nom& " diu: Ja he dinat, el compte per favor");
        monitor.leave(nom,Taula,fumador);
        Put_Line("          En " &nom&" S'EN VA");
        
      end if;
      
      

  end comensal_task;

  -----
  -- Array de tasques
  -----
  type Array_Noms is array (1..THREADS) of Unbounded_String;
  type Lect_esc is array (1..THREADS) of comensal_task;
  le         : Lect_esc;
  F          : File_Type;          --File
  Noms       : Array_Noms;         --Array de nombres
  Init       : Boolean;


begin
  Init       := TRUE;
    -- Leer archivo de nombres
   Open(F, In_File, Names_file);
   for I in Noms'Range loop
      Noms(I) := Ada.Strings.Unbounded.To_Unbounded_String(Get_Line(F));
   end loop;
   Close(F);

   Put_Line("++++++++++ El Maître està preparat");
   Put_Line("++++++++++ Hi ha 3 salons amb capacitat de 3 comensals cada un");
  -----
  -- Start les tasques
  -----
   for I in Noms'Range loop
     Init := not Init;
     le(I).Start(Noms(I),Init );
     delay(0.5);
   end loop;

end restaurant;
