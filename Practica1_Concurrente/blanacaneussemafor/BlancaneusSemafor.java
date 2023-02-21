/*
 * Blancaneus i els set nans viuen en una casa on nomes hi ha quatre cadires, 
 * que els nans utilitzen per menjar.
 * Quan un nan torna de treballar de la mina comprova si hi ha una cadira lliure
 * per seure. Si hi ha una cadira lliure, llavors indica la Blancaneus que ja 
 * està assegut, i espera pacientment el seu torn a ser servit. 
 * Quan ha estat servit, Blancaneus li indica que pot començar a menjar. 
 * El nan menja i quan acaba, deixa la cadira lliure i torna a la mina. 
 * Per la seva banda, Blancaneus quan no té cap nan pendent de servir, se'n va passejar
 *
 */
package pkg290_blanacaneussemafor;

import java.util.concurrent.Semaphore;

/**
 *
 * @author miquelmascarooliver
 */
public class BlancaneusSemafor {

    /**
     * @param args the command line arguments
     */
    static volatile int perMenjar = 0;
    static volatile int aDormir = 0;

    static Semaphore cadires = new Semaphore(4); //Comptador per cadires
    static Semaphore mutex = new Semaphore(1); //Protegeix perMenjar i aDormir
    static Semaphore torn = new Semaphore(0); //Torns de menjar

    public static void main(String[] args) throws InterruptedException {

        Thread b = new Thread(new Blancaneus());
        b.start();
        String[] nomNan = {"Mudet", "Esternuts", "Vergonyós", "Dormilega",
            "Feliç", "Rondinaire", "Savi"};
        Thread[] nan = new Thread[nomNan.length];
        for (int i = 0; i < nomNan.length; i++) {
            nan[i] = new Thread(new Nan(nomNan[i]));
            nan[i].start();
        }
        b.join();
        for (int i = 0; i < nomNan.length; i++) {
            nan[i].join();
        }
    }

}
