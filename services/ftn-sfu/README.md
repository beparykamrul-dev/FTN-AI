# FTN Realtime SFU

FTN's realtime media endpoint is backed by LiveKit as the primary SFU and coturn as the TURN relay. The repository also registers mediasoup as a secondary provider-neutral SFU capability.

## Runtime

- SFU: LiveKit WebRTC server
- TURN: coturn
- Client transports: WebRTC/UDP with WebSocket/HTTP signaling support
- Credentials: generated into the deployment `.env`; never commit secrets
- Control policy: SFU selection remains behind FTN health/capacity/latency policy

## Standalone deployment

Copy `.env.example` to `.env`, replace the placeholder values with locally generated secrets, then run:

`docker compose --env-file .env up -d`

For the normal FTN installation, `deploy/one-click/bootstrap.sh` generates the required SFU/TURN secrets and starts this stack automatically unless `FTN_ENABLE_REALTIME_SFU=false` is explicitly set.

## Production requirements

Publish the configured WebRTC/TURN ports through the FTN edge firewall and provide the real public IP/DNS/certificate configuration for the deployment. Do not put provider credentials in Git.
