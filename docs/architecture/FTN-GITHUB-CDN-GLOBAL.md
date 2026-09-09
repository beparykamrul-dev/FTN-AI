# FTN GitHub-Origin Global CDN

## Purpose

This profile adds a production-safe content-delivery plane to FTN-AI using the public GitHub repository as a versioned origin and an authorized CDN for global edge delivery.

## Traffic planes

```text
                    GLOBAL INTERNET
                           |
             +-------------+-------------+
             |                           |
       FTN BGP / IX / Transit      FTN CDN Plane
             |                           |
        FTN Core Routers          GitHub Origin
             |                           |
        FTN POP / Backbone         CDN Edge Network
             |                           |
        ISP/customer traffic       Public content traffic
```

The two planes are intentionally independent. CDN requests are content-delivery traffic and must not be counted as FTN IP transit capacity.

## GitHub-origin delivery

Repository: `beparykamrul-dev/FTN-AI`

For versioned public files, the jsDelivr GitHub endpoint has the form:

`https://cdn.jsdelivr.net/gh/beparykamrul-dev/FTN-AI@<immutable-ref>/<path>`

Use a release tag or commit SHA for production. Do not use an unversioned `main` URL for production assets.

## Suitable content

- public JavaScript/CSS assets
- public documentation assets
- versioned web application assets
- public release metadata and other distributable files

Never publish credentials, private configuration, customer data, or privileged API responses through this plane.

## FTN integration

The provider-neutral edge contracts already exist in FTN-AI. This profile binds GitHub-origin content delivery to the existing provider/capability model without replacing FTN DNS, routing, or control-plane contracts.

The implementation therefore supports:

1. GitHub as immutable source/version registry.
2. CDN as global content-delivery layer.
3. FTN Core/BGP as the separate network-transit layer.
4. Provider health and delivery metadata in FTN observability.
5. Approval-first policy for privileged changes.

## Operations

GitHub Actions validates the CDN profile on changes. Credentials remain outside Git and production content should use immutable refs. CDN provider configuration is additive and does not modify existing routing state.

## Boundary

GitHub and a public CDN do not provide FTN with global BGP transit. Actual Internet transit requires authorized upstream/IX/peering connectivity and remains under the FTN Core network plane.
