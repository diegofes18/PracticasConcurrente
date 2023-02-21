// AUTORS: MARC CAÑELLAS GOMEZ I DIEGO BERMEJO CABAÑAS
// ENLLAÇ AL VIDEO EXPLICATIU:

// Classe per a controlar l'estat (posició) del director a cada ronda durant l'execució
public enum StateDirector {
    FORA, // encara no ha començat la ronda --> si hi ha festa no intervé
    ESPERANT, // ha començat una ronda pero no ha entrat a la sala per no molestar --> si hi ha festa hi entra
                                                                                // --> si no queda ningú també entra
    DINS // el director ha entrat a la sala --> no podrà entrar cap estudiant
}
