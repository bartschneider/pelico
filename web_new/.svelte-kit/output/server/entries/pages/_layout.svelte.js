import { c as create_ssr_component, b as subscribe, v as validate_component } from "../../chunks/ssr.js";
import { p as page } from "../../chunks/stores.js";
import "bootstrap";
const Nav = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let $page, $$unsubscribe_page;
  $$unsubscribe_page = subscribe(page, (value) => $page = value);
  $$unsubscribe_page();
  return `<nav class="navbar navbar-expand-lg navbar-dark bg-dark"><div class="container"><a class="navbar-brand" href="/" data-svelte-h="svelte-1r735sm">Pelico</a> <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav" data-svelte-h="svelte-i8a0eh"><span class="navbar-toggler-icon"></span></button> <div class="collapse navbar-collapse" id="navbarNav"><ul class="navbar-nav"><li class="nav-item"><a class="${["nav-link", $page.url.pathname === "/" ? "active" : ""].join(" ").trim()}" href="/" data-svelte-h="svelte-1noo0iw">Games</a></li> <li class="nav-item"><a class="${["nav-link", $page.url.pathname === "/wishlist" ? "active" : ""].join(" ").trim()}" href="/wishlist" data-svelte-h="svelte-1j6w9f8">Wishlist</a></li> <li class="nav-item"><a class="${["nav-link", $page.url.pathname === "/shortlist" ? "active" : ""].join(" ").trim()}" href="/shortlist" data-svelte-h="svelte-1lqo9qf">Shortlist</a></li> <li class="nav-item"><a class="${["nav-link", $page.url.pathname === "/stats" ? "active" : ""].join(" ").trim()}" href="/stats" data-svelte-h="svelte-110pkte">Stats</a></li></ul></div></div></nav>`;
});
const Layout = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  return `${validate_component(Nav, "Nav").$$render($$result, {}, {}, {})} <main class="container mt-4">${slots.default ? slots.default({}) : ``}</main>`;
});
export {
  Layout as default
};
