# FTN Daily Realtime AI Fabric

## Purpose
Provider-neutral realtime voice/video/AI capability for FTN Web, Android, PC and TV applications.

## Capability Registry
- WebRTC realtime audio/video
- realtime voice and multimodal AI
- live streaming and recording adapters
- realtime transcription adapter
- call-quality telemetry: packet loss, bitrate, latency and connection events
- room/presence/session abstraction
- REST/webhook integration boundary
- AI-agent transport abstraction
- adaptive media quality and network-aware routing

## FTN Architecture
```text
FTN Client (Web / Android / PC / TV)
        |
        v
FTN Realtime Gateway
        |
        +--> WebRTC / Realtime Adapter
        +--> Voice AI Adapter
        +--> Video/Media Adapter
        +--> Recording Adapter
        +--> Transcription Adapter
        |
        v
FTN Event + Telemetry Fabric
        |
        +--> OpenTelemetry
        +--> Prometheus
        +--> ClickHouse
        +--> AIOps / Health
        |
        v
FTN AI Agent / Assistant
```

## Provider Boundary
Daily is an optional provider adapter, not a hard dependency. FTN keeps its own room/session model, telemetry schema, authorization boundary and provider abstraction so another WebRTC/media backend can be substituted later.

## Open Source First
Use open-source WebRTC/realtime components and Daily's open-source ecosystem (including Pipecat) as implementation references where compatible. Official Daily credentials/API are an optional integration layer and are not required by the FTN core.

## Security
- TLS/mTLS at FTN service boundaries
- short-lived room/session tokens
- least-privilege provider credentials
- provider secrets isolated from client applications
- audit events without retaining unnecessary session data

## AI
Realtime AI assistants can consume approved realtime audio/video events through the FTN AI gateway. Actions remain policy-controlled and auditable.

## Notes
Daily currently provides WebRTC SDKs for JavaScript, React Native, iOS, Android, Python and Flutter, plus realtime analytics, recording, transcription and AI-oriented capabilities. FTN should expose these through the provider adapter rather than coupling the core architecture to Daily.
