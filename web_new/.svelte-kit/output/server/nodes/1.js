

export const index = 1;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/fallbacks/error.svelte.js')).default;
export const imports = ["_app/immutable/nodes/1.Ck7naN6J.js","_app/immutable/chunks/-3cQnFvs.js","_app/immutable/chunks/IHki7fMi.js","_app/immutable/chunks/5I0edPBj.js"];
export const stylesheets = [];
export const fonts = [];
