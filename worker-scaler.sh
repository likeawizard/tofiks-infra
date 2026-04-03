#!/bin/sh
# Adaptive worker scaler
# Monitors lichess-bot for active games and restarts server-worker
# with the appropriate thread count.
# No games:  2 threads
# Games:     1 thread

DEBOUNCE=30  # seconds of stable state before scaling
POLL=10      # seconds between checks
WORKER=openbench-server-worker
BOT=lichess-bot
COMPOSE_FILE=/compose/docker-compose.yml
ENV_FILE=/compose/.env

CURRENT_THREADS=""
STABLE_COUNT=0
LAST_STATE=""

get_game_count() {
    docker exec "$BOT" sh -c 'ls /proc/*/exe 2>/dev/null | xargs -I{} readlink {} 2>/dev/null | grep -c tofiks' 2>/dev/null || echo "0"
}

scale_worker() {
    THREADS=$1
    if [ "$CURRENT_THREADS" = "$THREADS" ]; then
        return
    fi
    echo "$(date): Scaling server-worker to $THREADS threads"
    SERVER_WORKER_THREADS=$THREADS docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" up -d --no-deps server-worker
    CURRENT_THREADS=$THREADS
    echo "$(date): Server-worker now running with $THREADS threads"
}

echo "$(date): Worker scaler started (debounce=${DEBOUNCE}s, poll=${POLL}s)"
sleep 20

# Force initial scale based on current state
GAMES=$(get_game_count)
echo "$(date): Initial game count: $GAMES"
if [ "$GAMES" -gt 0 ]; then
    LAST_STATE="active"
    scale_worker 1
else
    LAST_STATE="idle"
    scale_worker 2
fi
STABLE_COUNT=$DEBOUNCE

while true; do
    GAMES=$(get_game_count)

    if [ "$GAMES" -gt 0 ]; then
        STATE="active"
        WANT=1
    else
        STATE="idle"
        WANT=2
    fi

    if [ "$STATE" != "$LAST_STATE" ]; then
        STABLE_COUNT=0
        LAST_STATE="$STATE"
    else
        STABLE_COUNT=$((STABLE_COUNT + POLL))
    fi

    if [ "$STABLE_COUNT" -ge "$DEBOUNCE" ] && [ "$CURRENT_THREADS" != "$WANT" ]; then
        scale_worker "$WANT"
        STABLE_COUNT=0
    fi

    # Periodic heartbeat every ~5 minutes
    if [ $((STABLE_COUNT % 300)) -eq 0 ] && [ "$STABLE_COUNT" -gt 0 ]; then
        echo "$(date): Heartbeat - games=$GAMES threads=$CURRENT_THREADS"
    fi

    sleep "$POLL"
done
