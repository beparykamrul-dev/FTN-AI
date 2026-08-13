# FTN Social

Optional FTN-native social/communications application family, inspired by the capabilities users expect from modern messaging platforms without depending on Telegram as a runtime dependency.

## Initial product boundary

- account and profile integration through FTN Identity
- private and group messaging
- realtime delivery through FTN WebSocket/Socket
- media/file references through FTN storage APIs
- notifications
- moderation and abuse controls
- device/session management
- message search
- optional voice/video integration through approved FTN media services
- Android/web clients delivered through the FTN App Store

## Architecture

```text
FTN Social Client
      ↓
FTN API / Identity
      ↓
FTN Social Service
      ├── PostgreSQL
      ├── FTN WebSocket
      ├── FTN Notification
      └── FTN Storage
```

The application is optional: it can be enabled and deployed from the FTN App Store when required.

## Security boundary

Authentication, authorization, rate limits, moderation, audit and data-retention policies remain under FTN Control Plane. Secrets and private message content are not exported to external telemetry providers.
