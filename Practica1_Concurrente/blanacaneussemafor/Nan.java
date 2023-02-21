/*
 * 
 */
package pkg290_blanacaneussemafor;

import java.util.logging.Level;
import java.util.logging.Logger;

/**
 *
 * @author miquelmascarooliver
 */
public class Nan implements Runnable {

    private String nom;
    private final int QUANTITAT = 2;

    public Nan(String nom) {
        this.nom = nom;
    }

    @Override
    public void run() {
        System.out.println("BON DIA som en " + this.nom);
        try {
            for (int i = 0; i < QUANTITAT; i++) {

                System.out.println(this.nom + " treballa a la mina");
                Thread.sleep((long) (Math.random() * 100000));
                System.out.println(this.nom + " espera per una cadira");
                BlancaneusSemafor.cadires.acquire();
                System.out.println(this.nom + " ja seu");

                BlancaneusSemafor.mutex.acquire();
                BlancaneusSemafor.perMenjar++;
                BlancaneusSemafor.mutex.release();

                System.out.println(this.nom + " espera torn per menjar");
                BlancaneusSemafor.torn.acquire();

                BlancaneusSemafor.mutex.acquire();
                BlancaneusSemafor.perMenjar--;
                BlancaneusSemafor.mutex.release();

                System.out.println("----------------> " + this.nom + " menja!!!");
                Thread.sleep((long) (Math.random() * 1000));
                BlancaneusSemafor.cadires.release();

            }
            BlancaneusSemafor.mutex.acquire();
            BlancaneusSemafor.aDormir++;
            BlancaneusSemafor.mutex.release();
            System.out.println(this.nom + " se'n va a DORMIR " + BlancaneusSemafor.aDormir + "/7");
        } catch (InterruptedException ex) {
            System.out.println(this.nom + " EXEPCIÓ!" + ex.getMessage());
        }
    }

}
