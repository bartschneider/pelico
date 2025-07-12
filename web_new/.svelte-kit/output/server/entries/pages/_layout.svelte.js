import { c as create_ssr_component } from "../../chunks/ssr.js";
import "bootstrap";
const Layout = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  return `${$$result.head += `<!-- HEAD_svelte-1v7xpgt_START --><link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css" rel="stylesheet"><!-- HEAD_svelte-1v7xpgt_END -->`, ""} <nav class="navbar navbar-expand-lg navbar-dark bg-dark" data-svelte-h="svelte-1628rer"><div class="container-fluid"><a class="navbar-brand" href="/"><i class="fas fa-gamepad me-2"></i>Pelico</a> <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav"><span class="navbar-toggler-icon"></span></button> <div class="collapse navbar-collapse" id="navbarNav"><ul class="navbar-nav me-auto"><li class="nav-item"><a class="nav-link" href="/">Collection</a></li> <li class="nav-item"><a class="nav-link" href="/wishlist">Wishlist</a></li> <li class="nav-item"><a class="nav-link" href="/shortlist">Shortlist</a></li> <li class="nav-item"><a class="nav-link" href="/stats">Statistics</a></li></ul></div></div></nav> ${slots.default ? slots.default({}) : ``}`;
});
export {
  Layout as default
};
