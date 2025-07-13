import { c as create_ssr_component, d as add_attribute, e as escape } from "../../../chunks/ssr.js";
const Page = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let wishlist = [];
  let searchQuery = "";
  return `<div class="container-fluid mt-4"><div class="d-flex justify-content-between align-items-center mb-4"><h1 data-svelte-h="svelte-76ul9b"><i class="fas fa-heart me-2 text-danger"></i>
      Wishlist</h1> <div class="d-flex gap-2"><button class="btn btn-outline-secondary" ${"disabled"}><i class="${"fas fa-sync-alt " + escape("fa-spin", true)}"></i> Refresh</button> <button class="btn btn-primary" data-svelte-h="svelte-bczvwd"><i class="fas fa-plus me-1"></i> Add to Wishlist</button></div></div>  <div class="row mb-4"><div class="col-md-6"><div class="input-group"><span class="input-group-text" data-svelte-h="svelte-1ag4q0b"><i class="fas fa-search"></i></span> <input type="text" class="form-control" placeholder="Search wishlist by title or platform..."${add_attribute("value", searchQuery)}> ${``}</div></div> <div class="col-md-6"><p class="text-muted mt-2 mb-0">${`${escape(wishlist.length)} items on your wishlist`}</p></div></div> ${`<div class="d-flex justify-content-center align-items-center" style="height: 300px;" data-svelte-h="svelte-12ozwn5"><div class="text-center"><div class="spinner-border text-primary mb-3" role="status"><span class="visually-hidden">Loading...</span></div> <p class="text-muted">Loading wishlist...</p></div></div>`}</div>  ${``}`;
});
export {
  Page as default
};
