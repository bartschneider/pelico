import { c as create_ssr_component, d as add_attribute, e as escape } from "../../chunks/ssr.js";
/* empty css                                                  */
import "bootstrap";
const Page = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let games = [];
  let searchQuery = "";
  return `<div class="container-fluid mt-4"><div class="d-flex justify-content-between align-items-center mb-4"><h1 data-svelte-h="svelte-8q22c4">Game Collection</h1> <div class="d-flex gap-2"><button class="btn btn-outline-secondary" ${"disabled"}><i class="${"fas fa-sync-alt " + escape("fa-spin", true)}"></i> Refresh</button> <button class="btn btn-primary" data-svelte-h="svelte-mzv9ad"><i class="fas fa-plus me-1"></i> Add Game</button></div></div>  <div class="row mb-4"><div class="col-md-6"><div class="input-group"><span class="input-group-text" data-svelte-h="svelte-1ag4q0b"><i class="fas fa-search"></i></span> <input type="text" class="form-control" placeholder="Search games by title, genre, or platform..."${add_attribute("value", searchQuery)}> ${``}</div></div> <div class="col-md-6"><p class="text-muted mt-2 mb-0">${`${escape(games.length)} games total`}</p></div></div> ${`<div class="d-flex justify-content-center align-items-center" style="height: 300px;" data-svelte-h="svelte-1g4zktf"><div class="text-center"><div class="spinner-border text-primary mb-3" role="status"><span class="visually-hidden">Loading...</span></div> <p class="text-muted">Loading games...</p></div></div>`}</div> ${``}`;
});
export {
  Page as default
};
