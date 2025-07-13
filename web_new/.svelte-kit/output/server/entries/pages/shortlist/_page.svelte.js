import { c as create_ssr_component, b as add_attribute, d as escape } from "../../../chunks/ssr.js";
/* empty css                                                     */
const Page = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let shortlist = [];
  let searchQuery = "";
  return `<div class="container-fluid mt-4"><div class="d-flex justify-content-between align-items-center mb-4"><h1 data-svelte-h="svelte-19nk407"><i class="fas fa-list me-2 text-warning"></i>
      Shortlist</h1> <div class="d-flex gap-2"><button class="btn btn-outline-secondary" ${"disabled"}><i class="${"fas fa-sync-alt " + escape("fa-spin", true)}"></i> Refresh</button> <button class="btn btn-primary" data-svelte-h="svelte-10j3eyc"><i class="fas fa-plus me-1"></i> Add to Shortlist</button></div></div>  <div class="row mb-4"><div class="col-md-6"><div class="input-group"><span class="input-group-text" data-svelte-h="svelte-1ag4q0b"><i class="fas fa-search"></i></span> <input type="text" class="form-control" placeholder="Search shortlist by title, platform, or reason..."${add_attribute("value", searchQuery, 0)}> ${``}</div></div> <div class="col-md-6"><p class="text-muted mt-2 mb-0">${`${escape(shortlist.length)} games to play next`}</p></div></div> ${`<div class="d-flex justify-content-center align-items-center" style="height: 300px;" data-svelte-h="svelte-1m07bg0"><div class="text-center"><div class="spinner-border text-primary mb-3" role="status"><span class="visually-hidden">Loading...</span></div> <p class="text-muted">Loading shortlist...</p></div></div>`}</div>  ${``}`;
});
export {
  Page as default
};
