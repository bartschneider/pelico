

export const index = 3;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/shortlist/_page.svelte.js')).default;
export const imports = ["_app/immutable/nodes/3.B8ORtNN5.js","_app/immutable/chunks/DAdVehBS.js","_app/immutable/chunks/B-DzudNv.js","_app/immutable/chunks/IHki7fMi.js"];
export const stylesheets = ["_app/immutable/assets/GameCard.DapLlfny.css"];
export const fonts = [];
