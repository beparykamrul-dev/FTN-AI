# FTN configuration package

Production configuration is loaded from environment variables. No credentials are stored in source control.

Required in production:

- `FTN_DATABASE_URL`
- `FTN_JWT_SECRET` with at least 32 characters

The package validates the effective configuration before it is published through `Init()`.
