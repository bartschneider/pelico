

export const index = 5;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/wishlist/_page.svelte.js')).default;
export const imports = ["_app/immutable/nodes/5.B2Yk0TSN.js","_app/immutable/chunks/iMvok6IP.js","_app/immutable/chunks/CHS-Rq5X.js","_app/immutable/chunks/IHki7fMi.js"];
export const stylesheets = [];
export const fonts = [];
