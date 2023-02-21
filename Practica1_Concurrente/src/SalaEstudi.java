// AUTORS: MARC CAÑELLAS GOMEZ I DIEGO BERMEJO CABAÑAS
// ENLLANÇ AL VIDEO EXPLICATIU:

import java.util.concurrent.Semaphore;

public class SalaEstudi {

    static final int MAX = 6; // nombre d'alumnes per a que hi hagi festa
    static final int numEstu = 12; // nombre d'estudiants durant l'execució
    static volatile int comptador = 0; // compta el nombre d'estudiants a la sala d'estudi
    static volatile StateDirector estat = StateDirector.FORA; // variable que ens indica l'estat del director
    static Semaphore increment = new Semaphore(1); // semàfor per a controlar la secció crítica i protegir el comptador
    static Semaphore espera = new Semaphore(0); // semàfor que controla quan el director pot entrar a la sala d'estudi

    static Semaphore potEntrarAlumne = new Semaphore(0);


    public static void main(String[] args) throws InterruptedException{

        System.out.println("Nombre total d'estudiants: " + numEstu);
        System.out.println("Nombre màxim d'estudiants: " + MAX + "\n");

        // cream el fil del director
        Thread director = new Thread(new Director());

        // llista amb els noms que utilitzarem en una execució
        String[] nomEstu = {"Diego","Jose","Marc","Maria","Julia","Will","Pere",
                "Antonia", "Xisca", "Miquel", "Alba", "Mateu", "Pau","Ana","Marina", "Joan", "Carles", "Joana",
                "Mireia", "Edu", "Biel", "Lluis", "Carme"};

        // array que conté els fils de tots els estudiants
        Thread[] students = new Thread[numEstu];

        // iniciam els fils del director i dels estudiants i feim els joins corresponents
        director.start();
        for (int i = 0; i < students.length; i++) {
            students[i] = new Thread(new Student(nomEstu[i]));
            students[i].start();
        }

        director.join();
        for (Thread student : students) {
            student.join();
        }
    }

}