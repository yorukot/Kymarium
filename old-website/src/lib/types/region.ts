// ============================================================================
// Region Types
// ============================================================================

export interface Region {
	id: string;
	name: string;
	displayName: string;
}

export interface MonitorRegion {
	id: string;
	monitorId: string;
	regionId: string;
}
