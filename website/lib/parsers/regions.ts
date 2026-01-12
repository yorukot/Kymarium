import { REGION_MAP } from "@/lib/parsers/data/region";
import type { RegionRawData, Region } from "@/lib/schemas/region";
import { COUNTRY_CODE_TO_NAME } from "./data/country";

export function parseRegion(raw: RegionRawData): Region {
  const meta = REGION_MAP[raw.name];
  const city = meta ? meta.localVariant ?? meta.subdivisionName : raw.name;
  const country = meta ? COUNTRY_CODE_TO_NAME[meta.countryCodeAlpha2] : "";
  const displayName = country ? `${city}, ${country}` : city;

  return {
    id: raw.id,
    name: raw.name,
    city,
    country,
    displayName,
    flag: meta ? meta.countryCodeAlpha2.toLowerCase() : "un",
  };
}

export function parseRegions(rawList: RegionRawData[]): Region[] {
  return rawList.map(parseRegion);
}
