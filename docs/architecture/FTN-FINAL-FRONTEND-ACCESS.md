# FTN Final Frontend & Service Access

## Objective

Provide one consistent customer experience across Web, Android and PC while exposing only the services actually entitled to the account.

## Customer flow

`familytimenet.com → Service Catalog → Requirements → Eligibility/KYC → Activation → Entitlement → Web/Android/PC access`

## Service entitlement

The backend is the source of truth. A client must never infer entitlement from UI state.

Each entitlement is scoped by:

- account/customer
- service
- plan/tier
- device/resource where applicable
- tenant
- allowed actions
- expiration/status

The frontend dynamically renders only active services. Examples include FTN Internet, FTNDNS, Hosting, Cloud, Drive, CCTV Cloud, FiberMap, Media/TV, E-commerce, AI Assistant and Developer Platform.

## Frontend surfaces

### Web

- public service catalog and service-specific landing pages
- account dashboard
- active-service dashboard
- billing/payment
- service setup
- device management
- FiberMap where entitled
- Hosting/Developer console where entitled
- Drive/Media/TV interfaces where entitled
- support and issue tracking
- notification/status center

### Android

- unified FTN account
- entitlement-driven home screen
- service-specific modules
- secure device enrollment
- Drive/Media/TV/AI/CCTV modules according to entitlement
- push/deep-link enrollment
- support and incident updates

### PC/Desktop

- FTN Connect/Tunnel/agent surfaces where entitled
- hosting/developer tools
- Drive sync where entitled
- device/service diagnostics
- secure enrollment

## API contract

All clients use the same authentication and authorization contract. A UI permission is never a security boundary.

`Identity → Authentication → Entitlement → Authorization → API → Service`

Sensitive operations additionally require policy/approval and, where applicable, customer confirmation.

## Service domains

Service-specific domains/subdomains are allocated by the control plane and registered through the FTN DNS/provider adapter layer. Examples:

- `account.familytimenet.com`
- `billing.familytimenet.com`
- `pay.familytimenet.com`
- `ai.familytimenet.com`
- `drive.familytimenet.com`
- `hosting.familytimenet.com`
- `cloud.familytimenet.com`
- `dev.familytimenet.com`
- `media.familytimenet.com`
- `tv.familytimenet.com`
- `shop.familytimenet.com`
- `support.familytimenet.com`
- `status.familytimenet.com`
- `dns.familytimenet.com`

The control plane chooses subdomain/path routing based on tenancy, security, certificate and service requirements.

## Device-service requests

For authorized router/CCTV/device services, the request may include brand/model, MAC and serial number. The system verifies ownership/authorization before any device access. Online devices can expose approved telemetry through a driver/API/agent boundary. Firmware or configuration deployment requires explicit customer consent, signed artifacts and compatibility checks.

## Secure notification

Provisioning can immediately deliver a short-lived enrollment link to a registered client while the requested service is still being prepared. Credentials and long-lived secrets are never sent by SMS. Access links expire and are revocable.

## Public and paid services

The website may expose public/open services and product information without forcing unnecessary KYC. Core, paid or restricted services follow their own requirements and eligibility rules.

## Demo/catalog surface

The public website can provide isolated feature previews for FTN Drive, Media, TV Player, App Store, AI Assistant, Hosting, E-commerce and other services. Demo environments must not contain real customer data or production credentials.

## Final acceptance

Frontend completion means:

1. Web, Android and PC consume stable contracts.
2. Entitlement is enforced server-side.
3. Service-specific UI is dynamically available only after activation.
4. DNS/domain routing is generated from the service registry.
5. Authentication, mTLS where required, RBAC and audit are enforced.
6. Device access is issue-bound and authorization-bound.
7. Provisioning links are short-lived and revocable.
8. Production clients contain no embedded secrets.
9. Real infrastructure validation remains a release gate.
