# FTN Customer Access, Service Catalog & KYC

## Customer journey

`familytimenet.com -> Category -> Service -> Requirements -> Eligibility -> Transparent KYC/verification -> Order/Activation -> Service -> Monitoring`

FTN should expose a public service catalog for open services such as media, TV player, app store and AI assistant. Core regulated, paid, network or account-bound services must pass their defined eligibility and onboarding requirements before activation.

## Service categories

- Open/public services: media, TV player, app store, AI assistant and other services explicitly marked public.
- FTN account services: account, billing, notifications, support and customer portal.
- Network services: internet packages, VPN/connectivity and other services requiring eligibility checks.
- Commerce services: marketplace, checkout, subscriptions and business services.
- Corporate/government services: tenant-specific services with stronger identity, authorization and contractual requirements.

## Requirements engine

Every service has a versioned requirement policy containing:

- eligibility rules;
- required account attributes;
- age/guardian requirements where applicable;
- service-area/coverage requirements;
- package/contract requirements;
- KYC level, only where legally and operationally required;
- payment requirements;
- device/network prerequisites;
- approval requirements;
- terms and consent requirements.

The website should show the user what information is required and why before submission. Requirements must be minimal and service-specific.

## FTN Number verification

The FTN number can be used as an account identifier to retrieve an existing verified customer record where lawful and authorized. The platform should use data minimization and avoid repeatedly asking users for information already verified.

Verification status may be automated, but it must **not be hidden from the user**. The UI should clearly state that an existing record was verified, what information was reused, and provide correction/support paths where required.

The system must never silently assume that an FTN number proves identity. Stronger verification is required when the service or law requires it.

## Preventing incorrect KYC data

Do not make users manually type information that FTN already has lawfully verified when reuse is permitted. For fields that must be entered:

- use structured selectors instead of free-form text where possible;
- validate format and checksum locally;
- normalize names/addresses consistently;
- validate date ranges;
- provide clear field-level errors;
- use confirmation screens before submission;
- never silently alter submitted identity data;
- preserve correction/audit records;
- require re-verification when a material identity field changes.

No design can guarantee that a user can "never" provide incorrect information, so the system must detect, prevent, flag and correct errors rather than claim certainty.

## Privacy and security

- collect only required KYC attributes;
- encrypt sensitive data in transit and at rest;
- strict RBAC/least privilege;
- mTLS between verification services;
- immutable audit events;
- retention/deletion policy;
- tenant isolation;
- no sale or unrelated reuse of identity data;
- explicit consent where required;
- authorized access only;
- no covert surveillance.

## Public vs core service boundary

Public services can remain discoverable and usable without forcing unnecessary KYC. Core services should progressively request only the requirements necessary for activation.

`Open service -> minimal account -> eligible service -> required verification -> activation`

This keeps the public website lightweight while maintaining stronger controls for network, billing, commerce and regulated services.

## Implementation principle

The catalog, requirements engine, identity service, KYC/verification adapters, billing gateway and service provisioning system must communicate through versioned APIs/events. A service can therefore change its requirements without redesigning the whole FTN platform.
