// ============================================================================
// Notification Types
// ============================================================================

export type NotificationType = 'discord' | 'telegram' | 'slack' | 'email';

export interface Notification {
	id: string;
	teamId: string;
	type: NotificationType;
	name: string;
	config:
		| DiscordNotificationConfig
		| SlackNotificationConfig
		| TelegramNotificationConfig
		| EmailNotificationConfig;
	updatedAt: string;
	createdAt: string;
}

export interface DiscordNotificationConfig {
	webhookUrl: string;
}

export interface SlackNotificationConfig {
	webhookUrl: string;
}

export interface TelegramNotificationConfig {
	botToken: string;
	chatId: string;
}

export interface EmailNotificationConfig {
	emailAddress: string[];
}
