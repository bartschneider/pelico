import { c as create_ssr_component, a as subscribe, e as each, b as add_attribute, d as escape, v as validate_component } from "../../chunks/ssr.js";
import { p as page } from "../../chunks/stores.js";
import { w as writable } from "../../chunks/index.js";
const Nav = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let $page, $$unsubscribe_page;
  $$unsubscribe_page = subscribe(page, (value) => $page = value);
  $$unsubscribe_page();
  return `<nav class="navbar navbar-expand-lg navbar-dark bg-dark"><div class="container"><a class="navbar-brand" href="/" data-svelte-h="svelte-1majyhn"><i class="fas fa-gamepad me-2"></i>Pelico</a> <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav" data-svelte-h="svelte-z4v2rj"><span class="navbar-toggler-icon"></span></button> <div class="collapse navbar-collapse" id="navbarNav"><ul class="navbar-nav ms-auto"><li class="nav-item"><a class="${["nav-link", $page.url.pathname === "/" ? "active" : ""].join(" ").trim()}" href="/" data-svelte-h="svelte-s6c229"><i class="fas fa-home me-1"></i>Dashboard</a></li> <li class="nav-item"><a class="${["nav-link", $page.url.pathname === "/collection" ? "active" : ""].join(" ").trim()}" href="/collection" data-svelte-h="svelte-10wkdpv"><i class="fas fa-gamepad me-1"></i>Collection</a></li> <li class="nav-item"><a class="${["nav-link", $page.url.pathname === "/scanner" ? "active" : ""].join(" ").trim()}" href="/scanner" data-svelte-h="svelte-17o1deu"><i class="fas fa-search me-1"></i>ROM Scanner</a></li> <li class="nav-item"><a class="${["nav-link", $page.url.pathname === "/sessions" ? "active" : ""].join(" ").trim()}" href="/sessions" data-svelte-h="svelte-18l02bp"><i class="fas fa-clock me-1"></i>Play Sessions</a></li> <li class="nav-item"><a class="${["nav-link", $page.url.pathname === "/platforms" ? "active" : ""].join(" ").trim()}" href="/platforms" data-svelte-h="svelte-3oc58a"><i class="fas fa-desktop me-1"></i>Platforms</a></li> <li class="nav-item"><a class="${["nav-link", $page.url.pathname === "/settings" ? "active" : ""].join(" ").trim()}" href="/settings" data-svelte-h="svelte-18wo2t0"><i class="fas fa-cog me-1"></i>Settings</a></li> <li class="nav-item"><a class="${["nav-link", $page.url.pathname === "/wishlist" ? "active" : ""].join(" ").trim()}" href="/wishlist" data-svelte-h="svelte-1fnystd"><i class="fas fa-heart me-1"></i>Wishlist</a></li> <li class="nav-item"><a class="${["nav-link", $page.url.pathname === "/shortlist" ? "active" : ""].join(" ").trim()}" href="/shortlist" data-svelte-h="svelte-1eh155g"><i class="fas fa-list me-1"></i>Shortlist</a></li> <li class="nav-item"><a class="${["nav-link", $page.url.pathname === "/stats" ? "active" : ""].join(" ").trim()}" href="/stats" data-svelte-h="svelte-n51853"><i class="fas fa-chart-bar me-1"></i>Stats</a></li></ul></div></div></nav>`;
});
const notifications = writable([]);
function getToastClasses(type) {
  const baseClasses = "mb-4 max-w-sm w-full bg-white shadow-lg rounded-lg pointer-events-auto ring-1 ring-black ring-opacity-5 overflow-hidden";
  switch (type) {
    case "success":
      return `${baseClasses} border-l-4 border-green-400`;
    case "error":
      return `${baseClasses} border-l-4 border-red-400`;
    case "warning":
      return `${baseClasses} border-l-4 border-yellow-400`;
    case "info":
      return `${baseClasses} border-l-4 border-blue-400`;
    default:
      return `${baseClasses} border-l-4 border-gray-400`;
  }
}
function getIconClasses(type) {
  const baseClasses = "w-5 h-5";
  switch (type) {
    case "success":
      return `${baseClasses} text-green-400`;
    case "error":
      return `${baseClasses} text-red-400`;
    case "warning":
      return `${baseClasses} text-yellow-400`;
    case "info":
      return `${baseClasses} text-blue-400`;
    default:
      return `${baseClasses} text-gray-400`;
  }
}
function getIcon(type) {
  switch (type) {
    case "success":
      return "fas fa-check-circle";
    case "error":
      return "fas fa-exclamation-circle";
    case "warning":
      return "fas fa-exclamation-triangle";
    case "info":
      return "fas fa-info-circle";
    default:
      return "fas fa-bell";
  }
}
const ToastNotification = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let $notifications, $$unsubscribe_notifications;
  $$unsubscribe_notifications = subscribe(notifications, (value) => $notifications = value);
  $$unsubscribe_notifications();
  return ` <div class="fixed top-4 right-4 z-50 space-y-4">${each($notifications, (notification) => {
    return `<div${add_attribute("class", getToastClasses(notification.type), 0)}><div class="p-4"><div class="flex items-start"><div class="flex-shrink-0"><i class="${escape(getIcon(notification.type), true) + " " + escape(getIconClasses(notification.type), true)}"></i></div> <div class="ml-3 w-0 flex-1 pt-0.5"><p class="text-sm font-medium text-gray-900">${escape(notification.title)}</p> <p class="mt-1 text-sm text-gray-500">${escape(notification.message)} </p></div> ${notification.dismissible ? `<div class="ml-4 flex-shrink-0 flex"><button class="bg-white rounded-md inline-flex text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500" data-svelte-h="svelte-1mmalri"><span class="sr-only">Close</span> <i class="fas fa-times w-5 h-5"></i></button> </div>` : ``} </div></div> </div>`;
  })}</div>`;
});
const Layout = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  return `${validate_component(Nav, "Nav").$$render($$result, {}, {}, {})} <main class="container-fluid py-4">${slots.default ? slots.default({}) : ``}</main> ${validate_component(ToastNotification, "ToastNotification").$$render($$result, {}, {}, {})}`;
});
export {
  Layout as default
};
