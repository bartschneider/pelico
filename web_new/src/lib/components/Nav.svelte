<script lang="ts">
	import { page } from '$app/stores';
</script>

<nav class="navbar navbar-expand-lg navbar-dark bg-dark">
	<div class="container">
		<a class="navbar-brand" href="/">
			<i class="fas fa-gamepad me-2"></i>Pelico
		</a>
		
		<button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav">
			<span class="navbar-toggler-icon"></span>
		</button>
		
		<div class="collapse navbar-collapse" id="navbarNav">
			<ul class="navbar-nav ms-auto">
				<li class="nav-item">
					<a class="nav-link" class:active={$page.url.pathname === '/'} href="/">
						<i class="fas fa-home me-1"></i>Dashboard
					</a>
				</li>
				<li class="nav-item">
					<a class="nav-link" class:active={$page.url.pathname === '/collection'} href="/collection">
						<i class="fas fa-gamepad me-1"></i>Collection
					</a>
				</li>
				<li class="nav-item">
					<a class="nav-link" class:active={$page.url.pathname === '/scanner'} href="/scanner">
						<i class="fas fa-search me-1"></i>ROM Scanner
					</a>
				</li>
				<li class="nav-item">
					<a class="nav-link" class:active={$page.url.pathname === '/sessions'} href="/sessions">
						<i class="fas fa-clock me-1"></i>Play Sessions
					</a>
				</li>
				<li class="nav-item">
					<a class="nav-link" class:active={$page.url.pathname === '/platforms'} href="/platforms">
						<i class="fas fa-desktop me-1"></i>Platforms
					</a>
				</li>
				<li class="nav-item">
					<a class="nav-link" class:active={$page.url.pathname === '/settings'} href="/settings">
						<i class="fas fa-cog me-1"></i>Settings
					</a>
				</li>
				<li class="nav-item">
					<a class="nav-link" class:active={$page.url.pathname === '/wishlist'} href="/wishlist">
						<i class="fas fa-heart me-1"></i>Wishlist
					</a>
				</li>
				<li class="nav-item">
					<a class="nav-link" class:active={$page.url.pathname === '/shortlist'} href="/shortlist">
						<i class="fas fa-list me-1"></i>Shortlist
					</a>
				</li>
				<li class="nav-item">
					<a class="nav-link" class:active={$page.url.pathname === '/stats'} href="/stats">
						<i class="fas fa-chart-bar me-1"></i>Stats
					</a>
				</li>
			</ul>
		</div>
	</div>
</nav>