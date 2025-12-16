import type { Monitor, Notification } from '../../../types';

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

export function decidedNotificationIcon(notification: Notification): string {
  switch (notification.type) {
    case 'discord':
      return 'ri:discord-fill';
    case 'telegram':
      return 'ri:telegram-fill';
    case 'email':
      return 'ri:mail-fill';
    default:
      return 'ri:notification-4-fill';
  }
}