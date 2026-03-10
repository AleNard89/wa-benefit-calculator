#!/bin/bash
cd "$(dirname "$0")"

COMPOSE="docker compose --env-file .env.dev.local -f docker-compose.base.yml -f docker-compose.dev.yml"

show_help() {
    echo ""
    echo "Benefit Calculator - Dev Tools"
    echo ""
    echo "Uso: ./dev.sh <comando> [opzioni]"
    echo ""
    echo "Comandi:"
    echo "  up [servizio]   Avvia lo stack (o un singolo servizio, con rebuild)"
    echo "  down            Ferma e rimuove i container"
    echo "  down --volumes  Ferma, rimuove container E cancella i dati (DB, Redis)"
    echo "  stop            Ferma i container senza rimuoverli"
    echo "  restart [serv.] Riavvia lo stack (o un singolo servizio)"
    echo "  logs [servizio] Mostra i log (api, ui, postgres, redis, migrations)"
    echo "  ps              Mostra lo stato dei container"
    echo "  build [serv.]   Ricostruisce le immagini (o un singolo servizio)"
    echo "  db              Apre psql sul database"
    echo ""
    echo "Esempi:"
    echo "  ./dev.sh up              # Avvia tutto"
    echo "  ./dev.sh up api          # Ricostruisce e riavvia solo le API"
    echo "  ./dev.sh restart api     # Riavvia solo le API"
    echo "  ./dev.sh logs api        # Log delle API"
    echo "  ./dev.sh down            # Ferma tutto"
    echo "  ./dev.sh down --volumes  # Ferma tutto + cancella dati"
    echo "  ./dev.sh db              # Apre psql"
    echo ""
}

show_banner() {
    echo ""
    echo "========================================"
    echo "   Benefit Calculator - DEV"
    echo "========================================"
    echo ""
    echo "  Frontend:   http://localhost:5174"
    echo "  API:        http://localhost:8082/api/health"
    echo "  PostgreSQL: localhost:5435"
    echo "  Redis:      localhost:6381"
    echo ""
    echo "  Login su http://localhost:5174"
    echo "    Email:    admin@example.com"
    echo "    Password: Admin123!"
    echo ""
    echo "========================================"
}

case "${1}" in
    up)
        shift
        if [ -n "${1}" ]; then
            echo "Ricostruendo e riavviando: ${1}..."
            $COMPOSE up -d --build "${1}"
        else
            $COMPOSE up -d --build
            show_banner
        fi
        ;;
    down)
        if [ "${2}" = "--volumes" ]; then
            echo "Fermando e rimuovendo container + volumi (DB, Redis)..."
            $COMPOSE down -v
        else
            $COMPOSE down
        fi
        echo "Stack fermato."
        ;;
    stop)
        $COMPOSE stop
        echo "Container fermati (usa './dev.sh up' per riavviare)."
        ;;
    restart)
        shift
        if [ -n "${1}" ]; then
            echo "Riavviando: ${1}..."
            $COMPOSE restart "${1}"
        else
            $COMPOSE down
            $COMPOSE up -d --build
            show_banner
        fi
        ;;
    logs)
        shift
        $COMPOSE logs -f "$@"
        ;;
    ps)
        $COMPOSE ps
        ;;
    build)
        shift
        $COMPOSE build "$@"
        ;;
    db)
        $COMPOSE exec postgres psql -U benefit_user -d benefit_calculator
        ;;
    *)
        show_help
        ;;
esac
