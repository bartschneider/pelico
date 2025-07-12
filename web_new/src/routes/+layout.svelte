<script>
  import 'bootstrap/dist/css/bootstrap.min.css';
  import "../app.css";
  import { onMount } from 'svelte';
  import 'bootstrap';

</script>

<svelte:head>
  <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css" rel="stylesheet">
</svelte:head>

<nav class="navbar navbar-expand-lg navbar-dark bg-dark">
  <div class="container-fluid">
    <a class="navbar-brand" href="/">
      <i class="fas fa-gamepad me-2"></i>Pelico
    </a>
    <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav">
      <span class="navbar-toggler-icon"></span>
    </button>
    <div class="collapse navbar-collapse" id="navbarNav">
      <ul class="navbar-nav me-auto">
        <li class="nav-item">
          <a class="nav-link" href="/">Collection</a>
        </li>
        <li class="nav-item">
          <a class="nav-link" href="/wishlist">Wishlist</a>
        </li>
        <li class="nav-item">
          <a class="nav-link" href="/shortlist">Shortlist</a>
        </li>
        <li class="nav-item">
          <a class="nav-link" href="/stats">Statistics</a>
        </li>
      </ul>
    </div>
  </div>
</nav>

<slot />
