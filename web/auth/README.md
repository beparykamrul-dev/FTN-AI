# FTN Shared Authentication UI

Responsive, service-aware authentication surface for FTN personal, family, and Family Time Network services.

## UX contract

- One identity, multiple provisioned FTN services.
- No public service self-registration.
- Service cards are rendered from the Control Panel service registry/API.
- Layout adapts from phone to wide desktop/NOC displays.
- Authentication state is never stored in localStorage; the server-managed secure session cookie is authoritative.

The files in this directory are presentation primitives. Wire them to the existing FTN web runtime/API rather than creating a second authentication system.
