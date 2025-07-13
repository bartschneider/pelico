<script lang="ts">
	import { notifications, dismissNotification } from '$lib/stores/notifications';
	import { fly } from 'svelte/transition';

	function getToastClasses(type: string): string {
		const baseClasses = 'mb-4 max-w-sm w-full bg-white shadow-lg rounded-lg pointer-events-auto ring-1 ring-black ring-opacity-5 overflow-hidden';
		
		switch (type) {
			case 'success':
				return `${baseClasses} border-l-4 border-green-400`;
			case 'error':
				return `${baseClasses} border-l-4 border-red-400`;
			case 'warning':
				return `${baseClasses} border-l-4 border-yellow-400`;
			case 'info':
				return `${baseClasses} border-l-4 border-blue-400`;
			default:
				return `${baseClasses} border-l-4 border-gray-400`;
		}
	}

	function getIconClasses(type: string): string {
		const baseClasses = 'w-5 h-5';
		
		switch (type) {
			case 'success':
				return `${baseClasses} text-green-400`;
			case 'error':
				return `${baseClasses} text-red-400`;
			case 'warning':
				return `${baseClasses} text-yellow-400`;
			case 'info':
				return `${baseClasses} text-blue-400`;
			default:
				return `${baseClasses} text-gray-400`;
		}
	}

	function getIcon(type: string): string {
		switch (type) {
			case 'success':
				return 'fas fa-check-circle';
			case 'error':
				return 'fas fa-exclamation-circle';
			case 'warning':
				return 'fas fa-exclamation-triangle';
			case 'info':
				return 'fas fa-info-circle';
			default:
				return 'fas fa-bell';
		}
	}
</script>

<!-- Toast Container -->
<div class="fixed top-4 right-4 z-50 space-y-4">
	{#each $notifications as notification (notification.id)}
		<div
			class={getToastClasses(notification.type)}
			transition:fly={{ x: 300, duration: 300 }}
		>
			<div class="p-4">
				<div class="flex items-start">
					<div class="flex-shrink-0">
						<i class="{getIcon(notification.type)} {getIconClasses(notification.type)}"></i>
					</div>
					<div class="ml-3 w-0 flex-1 pt-0.5">
						<p class="text-sm font-medium text-gray-900">
							{notification.title}
						</p>
						<p class="mt-1 text-sm text-gray-500">
							{notification.message}
						</p>
					</div>
					{#if notification.dismissible}
						<div class="ml-4 flex-shrink-0 flex">
							<button
								class="bg-white rounded-md inline-flex text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
								on:click={() => dismissNotification(notification.id)}
							>
								<span class="sr-only">Close</span>
								<i class="fas fa-times w-5 h-5"></i>
							</button>
						</div>
					{/if}
				</div>
			</div>
		</div>
	{/each}
</div>