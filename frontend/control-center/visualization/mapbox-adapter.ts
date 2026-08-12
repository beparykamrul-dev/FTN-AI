export type GeoNode = {
  id: string;
  name: string;
  status: "healthy" | "degraded" | "warning" | "critical" | "unknown" | "maintenance";
  latitude: number;
  longitude: number;
};

export type GeoLink = {
  id: string;
  from: string;
  to: string;
  status: GeoNode["status"];
  latencyMs: number | null;
};

export type MapFeature = {
  type: "Feature";
  geometry: {
    type: "Point" | "LineString";
    coordinates: number[] | number[][];
  };
  properties: Record<string, string | number | boolean | null>;
};

/**
 * Converts FTN-safe geographic state into provider-safe map features.
 * Private telemetry, credentials, raw logs and secrets must never enter this adapter.
 */
export function buildMapFeatures(nodes: GeoNode[], links: GeoLink[]): MapFeature[] {
  const byId = new Map(nodes.map((node) => [node.id, node]));
  const features: MapFeature[] = nodes.map((node) => ({
    type: "Feature",
    geometry: { type: "Point", coordinates: [node.longitude, node.latitude] },
    properties: { id: node.id, name: node.name, status: node.status },
  }));

  for (const link of links) {
    const from = byId.get(link.from);
    const to = byId.get(link.to);
    if (!from || !to) continue;

    features.push({
      type: "Feature",
      geometry: {
        type: "LineString",
        coordinates: [
          [from.longitude, from.latitude],
          [to.longitude, to.latitude],
        ],
      },
      properties: {
        id: link.id,
        status: link.status,
        latencyMs: link.latencyMs,
      },
    });
  }

  return features;
}
