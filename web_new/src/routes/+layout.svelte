<script>
	import '../app.css';
	import Nav from '$lib/components/Nav.svelte';
	import ToastNotification from '$lib/components/ui/ToastNotification.svelte';
</script>

<Nav />

<main class="container-fluid py-4">
	<slot />
</main>

<ToastNotification />

