

export const index = 4;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/stats/_page.svelte.js')).default;
export const imports = ["_app/immutable/nodes/4.CrKIvhsE.js","_app/immutable/chunks/iMvok6IP.js","_app/immutable/chunks/CHS-Rq5X.js","_app/immutable/chunks/IHki7fMi.js"];
export const stylesheets = ["_app/immutable/assets/4.BA8oacvy.css"];
export const fonts = [];
