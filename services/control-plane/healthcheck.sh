#!/bin/sh
set -eu
wget -q -O - http://127.0.0.1:8080/healthz >/dev/null
