// AUTORS: MARC CAÑELLAS GOMEZ I DIEGO BERMEJO CABAÑAS
// ENLLANÇ AL VIDEO EXPLICATIU:

// Classe per a controlar el funcionament del fil 'director'
public class Director implements Runnable{
    private StateDirector state; // estat del director mitjançant la classe StateDirector
    private int ronda; // comptador amb les rondes que fa el director en una execució
    private final int RONDES = 3; // nombre de rondes que fa el director en cada execució

    // Constructor de la classe
    public Director(){
        this.state = StateDirector.FORA; // el director comença amb l'estat 'FORA'
        this.ronda = 1;
    }

    @Override
    public void run() {
        try {
            for (int i = 0; i < RONDES; i++) {

                System.out.println("    El Sr. Director comença la ronda");

                // si comença la ronda i no hi ha ningú a la sala, el director entra i s'acaba la ronda
                if (SalaEstudi.comptador == 0) {
                    SalaEstudi.estat=StateDirector.DINS; // canvi de l'estat del director

                // si hi ha gent a la sala pero no supera el MAX, el director es queda esperant
                } else if (SalaEstudi.comptador > 0 && SalaEstudi.comptador < SalaEstudi.MAX) {
                    SalaEstudi.estat = StateDirector.ESPERANT; // canvi de l'esat del director

                    System.out.println("El Director està esperant per entrar. No molesta als que estudien");

                    // el director es quedarà esperant fins que el semàfor 'espera' li doni permís (sala buida o festa)
                    SalaEstudi.espera.acquire();

                    // si mentre el director espera hi ha festa, aquest ha d'entrar, aturar-la i acabar la ronda
                    if (SalaEstudi.comptador>= SalaEstudi.MAX){
                        SalaEstudi.estat=StateDirector.DINS; // canvi de l'estat del director

                        System.out.println("El Director està dins la sala d'estudi: S'HA ACABAT LA FESTA!");
                    }

                    SalaEstudi.espera.acquire(); // espera fins tenir permís per sortir de la sala
                    SalaEstudi.potEntrarAlumne.release(); // dona permis a un estudiant per entrar a la sala

                // si un cop el director inicia la ronda hi ha festa, entrarà a la sala, aturarà la festa i acabarà la ronda
                }else if(SalaEstudi.comptador >= SalaEstudi.MAX) {
                    SalaEstudi.estat = StateDirector.DINS; // canvi de l'estat del director

                    System.out.println("El Director està dins la sala d'estudi: S'HA ACABAT LA FESTA!");
                    SalaEstudi.espera.acquire(); // espera fins tenir permís per sortir de la sala
                    SalaEstudi.potEntrarAlumne.release(); // dona permis a un estudiant per entrar a la sala
                }


                System.out.println("    El Director veu que no hi ha ningú a la sala d'estudis"+
                        "\n El Director acaba la ronda "+(i+1)+ " de 3");

                SalaEstudi.estat=StateDirector.FORA; // canvi de l'estat del director

                Thread.sleep((long) (Math.random() * 6000));
            }

        }catch(InterruptedException e){
            throw new RuntimeException(e);
        }

    }
}
