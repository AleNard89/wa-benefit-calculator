#!/bin/bash
set -e
cd "$(dirname "$0")"

ENV_FILE=".env.prod"
COMPOSE="docker compose --env-file ${ENV_FILE} -f docker-compose.base.yml -f docker-compose.prod.yml"

# ------------------------------------------------------------------
# Controlli preliminari
# ------------------------------------------------------------------
check_env_file() {
    if [ ! -f "${ENV_FILE}" ]; then
        echo ""
        echo "ERRORE: File '${ENV_FILE}' non trovato."
        echo "  1. Copia .env.prod.example in .env.prod"
        echo "  2. Compila tutti i campi (in particolare JWT_SECRET, ENCRYPTION_KEY, password DB)"
        echo ""
        exit 1
    fi

    # Controlla che nessun placeholder sia rimasto
    if grep -q "CAMBIA_CON\|YOUR_SERVER\|your-api-key\|changeme" "${ENV_FILE}" 2>/dev/null; then
        echo ""
        echo "ATTENZIONE: Il file '${ENV_FILE}' contiene ancora dei placeholder."
        echo "  Aggiorna tutti i valori prima di avviare in produzione."
        echo ""
        exit 1
    fi
}

check_settings_file() {
    SETTINGS_FILE="vertical/api/src/settings.prod.json"
    if grep -q "YOUR_SERVER_IP_OR_DOMAIN" "${SETTINGS_FILE}" 2>/dev/null; then
        echo ""
        echo "ATTENZIONE: '${SETTINGS_FILE}' contiene ancora il placeholder 'YOUR_SERVER_IP_OR_DOMAIN'."
        echo "  Aggiorna host, baseUrl, allowedOrigins e refreshTokenCookieDomain con il vero IP/dominio."
        echo ""
        exit 1
    fi
}

# ------------------------------------------------------------------
# Help
# ------------------------------------------------------------------
show_help() {
    echo ""
    echo "Orbita — Production Tools"
    echo ""
    echo "Uso: ./prod.sh <comando> [opzioni]"
    echo ""
    echo "Comandi:"
    echo "  up              Avvia lo stack di produzione (build + up)"
    echo "  down            Ferma e rimuove i container"
    echo "  down --volumes  Ferma, rimuove container E volumi (ATTENZIONE: cancella i dati!)"
    echo "  stop            Ferma i container senza rimuoverli"
    echo "  restart         Ferma e riavvia tutto"
    echo "  update          Ricostruisce le immagini e riavvia (aggiornamento codice)"
    echo "  logs [servizio] Mostra i log (api, ui, postgres, redis, migrations)"
    echo "  ps              Mostra lo stato dei container"
    echo "  build [serv.]   Ricostruisce le immagini (senza riavviare)"
    echo "  db              Apre psql sul database"
    echo ""
    echo "Esempi:"
    echo "  ./prod.sh up             # Primo avvio (o dopo un 'down')"
    echo "  ./prod.sh update         # Dopo un 'git pull' con nuove modifiche"
    echo "  ./prod.sh logs api       # Log del backend"
    echo "  ./prod.sh logs ui        # Log di nginx"
    echo "  ./prod.sh down           # Ferma tutto (i dati restano nei volumi)"
    echo "  ./prod.sh down --volumes # Ferma tutto e CANCELLA i dati"
    echo ""
}

show_banner() {
    # Legge l'IP/host dal settings.prod.json se disponibile
    SERVER=$(grep -o '"host": *"[^"]*"' vertical/api/src/settings.prod.json 2>/dev/null | grep -o '"[^"]*"$' | tr -d '"' || echo "YOUR_SERVER_IP")
    echo ""
    echo "========================================"
    echo "   Orbita — PRODUZIONE"
    echo "========================================"
    echo ""
    echo "  App: http://${SERVER}"
    echo ""
    echo "========================================"
}

# ------------------------------------------------------------------
# Comandi
# ------------------------------------------------------------------
case "${1}" in
    up)
        check_env_file
        check_settings_file
        echo "Avviando stack di produzione..."
        $COMPOSE up -d --build
        show_banner
        ;;

    down)
        if [ "${2}" = "--volumes" ]; then
            echo "ATTENZIONE: Stai per cancellare tutti i dati (DB, media files)."
            read -p "Sei sicuro? (scrivi 'SI' per confermare): " CONFIRM
            if [ "${CONFIRM}" != "SI" ]; then
                echo "Operazione annullata."
                exit 0
            fi
            echo "Fermando e rimuovendo container + volumi..."
            $COMPOSE down -v
        else
            $COMPOSE down
        fi
        echo "Stack fermato."
        ;;

    stop)
        $COMPOSE stop
        echo "Container fermati (usa './prod.sh up' per riavviare)."
        ;;

    restart)
        check_env_file
        check_settings_file
        echo "Riavviando stack..."
        $COMPOSE down
        $COMPOSE up -d --build
        show_banner
        ;;

    update)
        check_env_file
        check_settings_file
        echo "Aggiornando stack (rebuild immagini + restart)..."
        $COMPOSE build
        $COMPOSE up -d
        echo "Aggiornamento completato."
        ;;

    logs)
        shift
        $COMPOSE logs -f "$@"
        ;;

    ps)
        $COMPOSE ps
        ;;

    build)
        check_env_file
        shift
        $COMPOSE build "$@"
        ;;

    db)
        check_env_file
        set -a; source "${ENV_FILE}"; set +a
        $COMPOSE exec postgres psql -U "${POSTGRES_USER}" -d "${POSTGRES_DB}"
        ;;

    *)
        show_help
        ;;
esac
