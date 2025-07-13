import { writable } from 'svelte/store';

export interface Notification {
	id: string;
	type: 'success' | 'error' | 'warning' | 'info';
	title: string;
	message: string;
	duration?: number;
	dismissible?: boolean;
}

export const notifications = writable<Notification[]>([]);

let notificationId = 0;

export const addNotification = (notification: Omit<Notification, 'id'>) => {
	const id = String(++notificationId);
	const newNotification: Notification = {
		id,
		duration: 5000,
		dismissible: true,
		...notification
	};

	notifications.update(notifications => [...notifications, newNotification]);

	// Auto-dismiss if duration is set
	if (newNotification.duration && newNotification.duration > 0) {
		setTimeout(() => {
			dismissNotification(id);
		}, newNotification.duration);
	}

	return id;
};

export const dismissNotification = (id: string) => {
	notifications.update(notifications => 
		notifications.filter(notification => notification.id !== id)
	);
};

export const clearAllNotifications = () => {
	notifications.set([]);
};

// Helper functions for common notification types
export const showSuccess = (title: string, message: string) => 
	addNotification({ type: 'success', title, message });

export const showError = (title: string, message: string) => 
	addNotification({ type: 'error', title, message, duration: 8000 });

export const showWarning = (title: string, message: string) => 
	addNotification({ type: 'warning', title, message });

export const showInfo = (title: string, message: string) => 
	addNotification({ type: 'info', title, message });