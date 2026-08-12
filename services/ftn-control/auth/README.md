# FTN Shared Sign In / Sign Up

The FTN Control Panel is the service provisioning authority.

## Flow

1. Administrator enables a service in the Control Panel.
2. The service appears in the shared FTN authentication UI.
3. Administrator provisions an eligible user/family/business account for that service.
4. The user signs in with the shared FTN identity.
5. Access is granted only when the identity has an active service assignment.

## Rules

- No arbitrary public service signup.
- Service visibility is driven by the registry.
- Authentication and authorization are separate.
- A disabled service immediately becomes unavailable for new access.
- Existing sessions must be revoked according to the central session policy when access is removed.
- Mail accounts are provisioned by FTN Control Panel; the mail service does not create accounts independently.
