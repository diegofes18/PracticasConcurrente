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
class Blancaneus implements Runnable {

    @Override
    public void run() {
        System.out.println("BON DIA som na Blancaneus");
        while (BlancaneusSemafor.aDormir < 7) {
            try {
                BlancaneusSemafor.mutex.acquire();
                if (BlancaneusSemafor.perMenjar == 0) {
                    BlancaneusSemafor.mutex.release();
                    System.out.println("Blancaneus se'n va a fer una passeig");
                    Thread.sleep((long)(Math.random() * 50000));
                } else {
                    BlancaneusSemafor.mutex.release();
                    System.out.println("Blancaneus cuina per un nan");
                    Thread.sleep((long)(Math.random() * 5000));
                    System.out.println("Blancaneus té el menjat cuit");
                    BlancaneusSemafor.torn.release();
                }
            } catch (InterruptedException ex) {
                System.out.println(" Blancaneus EXEPCIÓ!" + ex.getMessage());
            }
        }
        System.out.println("Blancaneus s'en va a DORMIR");
    }

}
