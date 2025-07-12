import { c as create_ssr_component } from "../../chunks/ssr.js";
/* empty css                                                  */
import "bootstrap";
const Page = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  return `<div class="container-fluid mt-4"><div class="d-flex justify-content-between align-items-center mb-4"><h1 data-svelte-h="svelte-8q22c4">Game Collection</h1> <button class="btn btn-primary" data-svelte-h="svelte-1pqzu4l"><i class="fas fa-plus me-1"></i> Add Game</button></div> ${`<p data-svelte-h="svelte-uis8ph">Loading games...</p>`}</div> ${``}`;
});
export {
  Page as default
};
