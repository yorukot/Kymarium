import type { Monitor } from "../../types";

export function monitorTarget(monitor: Monitor): string {
	switch (monitor.type) {
		case 'http':
			if ('url' in monitor.config) {
				return monitor.config.url;
			}
			break;

		case 'ping':
			if ('host' in monitor.config) {
				return monitor.config.host;
			}
			break;
	}

	throw new Error(`Invalid monitor config for type: ${monitor.type}`);
}
