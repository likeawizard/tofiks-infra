#!/bin/sh
set -e

if [ -z "$LICHESS_TOKEN" ]; then
  echo "ERROR: LICHESS_TOKEN environment variable is not set"
  exit 1
fi

# Inject token into config
sed "s|\${LICHESS_TOKEN}|${LICHESS_TOKEN}|g" config.yml.template > config.yml

exec python lichess-bot.py
