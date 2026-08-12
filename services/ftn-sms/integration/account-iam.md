# FTN SMS Account/IAM Integration V1

The SMS service uses the existing FTN Control authentication and service-assignment model. It does not create a second identity system.

## Authorization flow

1. Resolve the authenticated FTN identity from the existing session/token context.
2. Require an active `sms` service assignment.
3. Evaluate the capability required by the operation (`sms.send`, `sms.schedule`, `sms.status.read`, `sms.resend`, or `sms.cancel`).
4. Apply SMS policy, Sender ID state, rate limits, and—when required—approval policy.
5. Only then enqueue or mutate the SMS record.

## Assignment rule

Only an active assignment is authorized. Pending, suspended, and revoked assignments are denied.

The central Control Panel remains the service-provisioning authority. Public signup and self-service creation of privileged SMS access are not permitted.

## Operation mapping

| Operation | Capability | Extra control |
|---|---|---|
| Send | `sms.send` | active Sender ID + rate limit |
| Schedule | `sms.schedule` | active Sender ID + schedule policy |
| Read status | `sms.status.read` | ownership/scope check |
| Resend | `sms.resend` | audit + retry policy |
| Cancel | `sms.cancel` | audit + message-state check |

## AI-generated messages

AI-generated SMS remains subject to the FTN approval policy. The SMS service must not treat an AI request as implicit authorization.

## Security boundary

The SMS service may consume authorization decisions from the existing control-plane/auth layer, but must never bypass them. Transport adapters receive only already-authorized delivery requests.
