// AUTORS: MARC CAÑELLAS GOMEZ I DIEGO BERMEJO CABAÑAS
// ENLLANÇ AL VIDEO EXPLICATIU:

// Classe per a controlar el funcionament del fil 'estudiant'
public class Student implements Runnable {
    private String nom; // atribut amb el nom de l'estudiant
    public Student(String name) {
        this.nom = name;
    }

    @Override
    public void run() {
        try {
            Thread.sleep((long) (Math.random() * 5000));

            // Si el professor es troba dins la sala no podrà entrar cap alumne
            if(SalaEstudi.estat==StateDirector.DINS) {
                SalaEstudi.potEntrarAlumne.acquire(); // espera fins que obté permís --> per part del director o per part d'un estudiant
                SalaEstudi.potEntrarAlumne.release(); // dona permís al següent estudiant que vol entrar
                Thread.sleep((long) (Math.random() * 3000));
            }

            // semàfor per a controlar que dos estudiants no modifiquin el comptador al mateix temps
            SalaEstudi.increment.acquire(); // mentre un fil modifica la variable cap la podrà modificar
            SalaEstudi.comptador++;
            System.out.println(this.nom + ": entra a la sala d'estudi, nombre d'estudiants: " + SalaEstudi.comptador);

            // si la sala d'estudi supera el nombre MAX hi haurà festa
            if (SalaEstudi.comptador >= SalaEstudi.MAX) {
                System.out.println(this.nom + ": FESTA!!!");

                // si el director es trobava esperant per entrar a la sala, donarem permís al director per entrar
                if (SalaEstudi.estat == StateDirector.ESPERANT) {
                    SalaEstudi.estat=StateDirector.DINS; // canvi de l'estat del director
                    System.out.println(this.nom+": ALERTA que vé el director!!!!!!!!!");
                    SalaEstudi.espera.release(); // donam permís al director
                }

            // si el comptador no supera el nombre MAX, l'estudiant estudiarà
            } else {
                    System.out.println(this.nom + " estudia");
            }

            SalaEstudi.increment.release(); // l'estudiant allibera el semàfor perquè no modificarà el comptador

            Thread.sleep((long) (Math.random() * 10000)); // sleep per a simular que l'estudiant estudia

            SalaEstudi.increment.acquire(); // quan deixa d'estudiar demanarà permís per a entrar a la secció crítica
            System.out.println(this.nom + ": ha sortit de la sala, nombre d'estudiants: " + (SalaEstudi.comptador-1));

            // si és el darrer estudiant que queda a la sala d'estudi
            if (SalaEstudi.comptador == 1) {
                if (SalaEstudi.estat == StateDirector.ESPERANT) {
                    System.out.println(this.nom + ": ADEU Senyor Director, pot entrar si vol, no hi ha ningú");
                    SalaEstudi.espera.release(); // donam permís al director per a que entri a la sala
                    SalaEstudi.espera.release(); // donam permís al director per a que pugui sortir de la sala

                } else if(SalaEstudi.estat==StateDirector.DINS) {
                    System.out.println(this.nom + ": ADEU Senyor Director es queda sol");
                    SalaEstudi.espera.release(); // donam permís al director per a que pugui sortir de la sala
                }
            }

            SalaEstudi.comptador--;
            SalaEstudi.increment.release(); // l'estudiant allibera el semàfor perquè no modificarà el comptador

        }

        catch(InterruptedException e){
            throw new RuntimeException(e);
        }

    }
}
