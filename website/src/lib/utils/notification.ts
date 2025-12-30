import type { Notification, NotificationType } from '$lib/types';

export function decidedNotificationIcon(notification: Notification): string {
	switch (notification.type) {
		case 'discord':
			return 'ri:discord-fill';
		case 'slack':
			return 'ri:slack-fill';
		case 'telegram':
			return 'ri:telegram-fill';
		case 'email':
			return 'lucide:mail';
		default:
			return 'ri:notification-4-fill';
	}
}

export const notificationTypeMeta: Record<
	NotificationType,
	{
		label: string;
		icon: string;
		description: string;
	}
> = {
	telegram: {
		label: 'Telegram',
		icon: 'ri:telegram-fill',
		description: 'Send alerts via a Telegram bot'
	},
	discord: {
		label: 'Discord',
		icon: 'ri:discord-fill',
		description: 'Send alerts to a Discord webhook'
	},
	slack: {
		label: 'Slack',
		icon: 'ri:slack-fill',
		description: 'Send alerts to a Slack webhook'
	},
	email: {
		label: 'Email',
		icon: 'lucide:mail',
		description: 'Send alerts to inboxes'
	}
};
