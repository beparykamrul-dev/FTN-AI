# FTN Codec Fabric

FTN Codec Fabric is the provider-neutral capability boundary for compression, chunking, deduplication, media encoding and media processing.

## Design rules

- Codec/encoder implementations are replaceable workers, not control-plane dependencies.
- Hardware acceleration is preferred when the target node advertises the required device capability.
- Large-file transfer is lossless by default: chunk -> hash -> optional deduplication -> compression -> encrypted transport -> verify -> reassemble.
- Originals are preserved unless an explicit lifecycle policy says otherwise.
- Output is verified before it becomes the published artifact.
- CPU/RAM/storage/GPU/network/path health are placement inputs, so work can move between eligible FTN nodes.

## Registered external workers

- `HWEncoderX`: H.265 hardware-accelerated media worker. FTN registers it as an isolated implementation for VAAPI/NVENC/QSV-capable nodes.
- `jumpcutter`: media-cut/video-processing worker. FTN registers it as an isolated media implementation.

The repository records integration metadata and contracts; it does not vendor third-party source into the FTN control plane.

## Service bindings

The same capability layer can serve FTN API, proxy, tunnel, capture, storage, media and IP-stream workloads. A service only requests the capability it needs; the scheduler chooses a compatible implementation and node.
