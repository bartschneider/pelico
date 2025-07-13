<script lang="ts">
	import { page } from '$app/stores';
</script>

<nav class="navbar navbar-expand-lg navbar-dark bg-dark">
	<div class="container">
		<a class="navbar-brand" href="/">Pelico</a>
		
		<button 
			class="navbar-toggler" 
			type="button" 
			data-bs-toggle="collapse" 
			data-bs-target="#navbarNav"
		>
			<span class="navbar-toggler-icon"></span>
		</button>
		
		<div class="collapse navbar-collapse" id="navbarNav">
			<ul class="navbar-nav">
				<li class="nav-item">
					<a 
						class="nav-link" 
						class:active={$page.url.pathname === '/'} 
						href="/"
					>
						Games
					</a>
				</li>
				<li class="nav-item">
					<a 
						class="nav-link" 
						class:active={$page.url.pathname === '/wishlist'} 
						href="/wishlist"
					>
						Wishlist
					</a>
				</li>
				<li class="nav-item">
					<a 
						class="nav-link" 
						class:active={$page.url.pathname === '/shortlist'} 
						href="/shortlist"
					>
						Shortlist
					</a>
				</li>
				<li class="nav-item">
					<a 
						class="nav-link" 
						class:active={$page.url.pathname === '/stats'} 
						href="/stats"
					>
						Stats
					</a>
				</li>
			</ul>
		</div>
	</div>
</nav>