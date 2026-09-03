#!/bin/sh
set -eu

: "${FTN_SFU_API_KEY:?FTN_SFU_API_KEY is required}"
: "${FTN_SFU_API_SECRET:?FTN_SFU_API_SECRET is required}"

umask 077
sed \
  -e "s/__FTN_SFU_API_KEY__/${FTN_SFU_API_KEY}/g" \
  -e "s/__FTN_SFU_API_SECRET__/${FTN_SFU_API_SECRET}/g" \
  /etc/ftn-sfu/livekit.yaml.template > /run/ftn-sfu/livekit.yaml

exec /livekit-server --config /run/ftn-sfu/livekit.yaml
