

export const index = 5;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/stats/_page.svelte.js')).default;
export const imports = ["_app/immutable/nodes/5.C4XHxQk_.js","_app/immutable/chunks/BdL0vCbg.js","_app/immutable/chunks/CLkIolx5.js","_app/immutable/chunks/IHki7fMi.js"];
export const stylesheets = ["_app/immutable/assets/5.BA8oacvy.css"];
export const fonts = [];
