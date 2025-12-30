import type { Notification, NotificationType } from '../../types';
import { apiRequest } from './utils';

export type NotificationsResponse = {
	message: string;
	data: Notification[];
};

export type NotificationCreateRequest = {
	type: NotificationType;
	name: string;
	config:
		| {
				webhook_url: string;
		  }
		| {
				bot_token: string;
				chat_id: string;
		  }
		| {
				email_address: string[];
		  };
};

export type NotificationCreateResponse = {
	message: string;
	data: Notification;
};

export type NotificationDeleteResponse = {
	message: string;
};

export type NotificationTestResponse = {
	message: string;
};

export type NotificationUpdateRequest = Partial<{
	type: NotificationType;
	name: string;
	config:
		| {
				webhook_url: string;
		  }
		| {
				bot_token: string;
				chat_id: string;
		  }
		| {
				email_address: string[];
		  };
}>;

export type NotificationUpdateResponse = {
	message: string;
	data: Notification;
};

export function getNotifications(teamID: string): Promise<NotificationsResponse> {
	return apiRequest<NotificationsResponse>(`/teams/${teamID}/notifications`, {
		defaultError: 'Failed to fetch notifications'
	});
}

export function createNotification(
	teamID: string,
	payload: NotificationCreateRequest
): Promise<NotificationCreateResponse> {
	return apiRequest<NotificationCreateResponse>(`/teams/${teamID}/notifications`, {
		method: 'POST',
		body: payload,
		defaultError: 'Failed to create notification'
	});
}

export function deleteNotification(
	teamID: string,
	notificationID: string
): Promise<NotificationDeleteResponse> {
	return apiRequest<NotificationDeleteResponse>(
		`/teams/${teamID}/notifications/${notificationID}`,
		{
			method: 'DELETE',
			defaultError: 'Failed to delete notification'
		}
	);
}

export function testNotification(
	teamID: string,
	notificationID: string
): Promise<NotificationTestResponse> {
	return apiRequest<NotificationTestResponse>(
		`/teams/${teamID}/notifications/${notificationID}/test`,
		{
			method: 'POST',
			defaultError: 'Failed to send test notification'
		}
	);
}

export function updateNotification(
	teamID: string,
	notificationID: string,
	payload: NotificationUpdateRequest
): Promise<NotificationUpdateResponse> {
	return apiRequest<NotificationUpdateResponse>(
		`/teams/${teamID}/notifications/${notificationID}`,
		{
			method: 'PATCH',
			body: payload,
			defaultError: 'Failed to update notification'
		}
	);
}
