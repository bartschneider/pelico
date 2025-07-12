

export const index = 4;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/stats/_page.svelte.js')).default;
export const imports = ["_app/immutable/nodes/4.DLFRky_h.js","_app/immutable/chunks/-3cQnFvs.js","_app/immutable/chunks/IHki7fMi.js"];
export const stylesheets = [];
export const fonts = [];
