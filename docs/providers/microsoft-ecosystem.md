# Microsoft Ecosystem Reference

FTN-AI provider-neutral reference profile for Microsoft ecosystem sources.

## Sources
- Microsoft365DSC
- MicrosoftDocs
- OfficeDev
- AzureAD
- MicrosoftLearning
- Azure
- microsoftarchive
- microsoft

## Capability mapping
- Configuration-as-code and drift detection
- Microsoft Graph / identity integration patterns
- Azure infrastructure and service integration
- Microsoft 365 application/platform integration
- Documentation and SDK reference ingestion
- DevOps / learning / operational runbooks
- Windows and Microsoft platform compatibility
- Security, identity, policy and compliance reference patterns

## FTN boundary
These sources are references and optional adapters, not hard dependencies. FTN core remains provider-neutral. Public/open-source components are evaluated by license and compatibility before reuse. Official APIs and credentials can be attached later through the same adapter interface.

## Runtime model
Provider source -> adapter -> capability normalizer -> central FTN database -> metrics/telemetry -> AI agent -> unified control plane.
