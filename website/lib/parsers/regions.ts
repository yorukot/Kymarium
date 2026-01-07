import { REGION_MAP } from "@/lib/parsers/data/region";
import type { RegionRawData, Region } from "@/lib/schemas/region";

export function parseRegion(raw: RegionRawData): Region {
  const meta = REGION_MAP[raw.name];
  const displayName = meta
    ? `${meta.localVariant ?? meta.subdivisionName}, ${meta.countryCodeAlpha2}`
    : raw.name;

  return {
    id: raw.id,
    name: raw.name,
    displayName,
    flag: meta ? meta.countryCodeAlpha2.toLowerCase() : "un",
  };
}

export function parseRegions(rawList: RegionRawData[]): Region[] {
  return rawList.map(parseRegion);
}
