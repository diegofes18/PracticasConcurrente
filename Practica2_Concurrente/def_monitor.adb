--AUTORS: DIEGO BERMEJO CABAÑAS I MARC CAÑELLAS GOMEZ
--VIDEOS EXPLICATIU: https://www.youtube.com/watch?v=6LDuejoOCYU

with Text_IO; use  Text_IO;
with Ada.Strings.Unbounded;    use Ada.Strings.Unbounded;
with Ada.Text_IO.Unbounded_IO; use Ada.Text_IO.Unbounded_IO;
with Ada.Task_Identification;  use Ada.Task_Identification;


package body def_monitor is
    protected body maitre is
        

        function disponible(tipo: in Integer) return Integer is
            trobat: Integer :=-1;
            I:Integer :=1;
            buscando : boolean:= True;
            begin 
                while (buscando) and (I<4) loop
                    if(arrCap(I)>0) and( (arrType(I)=1) or (arrType(I)=tipo) )then
                            trobat:=arrType(I);
                            buscando:=False;
                    end if;
                    I:=I+1;
                end loop;
                return trobat;
            end disponible;

        --Demanam taula de part d'un fumador
        entry demanaTaulaF(Nom: in Unbounded_String; Taula: out Integer) when disponible(2)>0 is
        --VARIABLES LOCALES
        I : integer:=1;
        sala:integer:=0;
        trobat:integer:=0;

        begin
            loop
                if arrCap(I)>0 and (arrType(I)=1 or arrType(I)=2) then
                    sala:=I;
                    trobat:=1;
                else
                    I:=I+1;
                end if;
                
                exit when trobat=1 or I=4;
            end loop;

            arrType(sala):=2;
            Taula:=sala;
            arrCap(sala):=arrCap(sala)-1;

            Put_Line("----------En "& Nom &" te taula al saló de fumadors " & sala'img & ", Disponibilidad:" & arrCap(I)'img);


        end demanaTaulaF;
        
        entry demanaTaulaNoF(Nom: in Unbounded_String;Taula: out Integer) when disponible(3)>0 is
            --VARIABLES LOCALES
        I : integer:=1;
        sala:integer:=0;
        trobat:integer:=0;

        begin
            loop
                if arrCap(I)>0 and (arrType(I)=1 or arrType(I)=3) then
                    sala:=I;
                    trobat:=1;
                else
                    I:=I+1;
                end if;
                
                exit when trobat=1 or I=4;
            end loop;

            arrType(sala):=3;
            Taula:=sala;
            arrCap(sala):=arrCap(sala)-1;

            Put_Line("----------En "& Nom &" te taula al saló de fumadors " & sala'img & ", Disponibilidad:" & arrCap(I)'img);
            
        end demanaTaulaNoF;

        
        --Sortim y lliberam la capacitat de la sala on es trobava el comensal
        procedure leave(Nom: in Unbounded_String; sala: in integer; fuma: in Boolean) is 
           begin
                arrCap(sala):=arrCap(sala)+1;

                --Si tothom ha partit el tipus sirà lliure
                if(arrCap(sala)=3) then
                    arrType(sala):=1;
                end if;

            if (fuma) then
                Put_Line("----------En "& Nom &" allibera una taula del saló "&sala'img&". Disponibilidad: "&arrCap(sala)'img& ". Tipus " &arrType(sala)'img);
            else
                Put_Line("++++++++++En "& Nom &" allibera una taula del saló "&sala'img&". Disponibilidad: "&arrCap(sala)'img& ". Tipus " &arrType(sala)'img);
            end if;
            
        end leave;

    end maitre;
end def_monitor;