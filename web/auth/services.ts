export type FTNService = {
  id: string;
  name: string;
  description?: string;
  enabled: boolean;
};

export function visibleServices(services: FTNService[]): FTNService[] {
  return services.filter((service) => service.enabled);
}

export function servicePath(service: FTNService): string {
  return `/services/${encodeURIComponent(service.id)}`;
}
