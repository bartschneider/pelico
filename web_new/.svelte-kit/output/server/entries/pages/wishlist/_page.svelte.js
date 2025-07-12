import { c as create_ssr_component } from "../../../chunks/ssr.js";
/* empty css                                                     */
const Page = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  return `<div class="container-fluid mt-4"><div class="d-flex justify-content-between align-items-center mb-4" data-svelte-h="svelte-qggn8k"><h1>Wishlist</h1> <button class="btn btn-primary"><i class="fas fa-plus me-1"></i> Add Game to Wishlist</button></div> ${`<p data-svelte-h="svelte-16qtdg7">Loading wishlist...</p>`}</div>`;
});
export {
  Page as default
};
