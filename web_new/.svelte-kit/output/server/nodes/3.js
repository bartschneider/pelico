

export const index = 3;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/shortlist/_page.svelte.js')).default;
export const imports = ["_app/immutable/nodes/3.B0Ng4oh7.js","_app/immutable/chunks/-3cQnFvs.js","_app/immutable/chunks/CcNbzAFb.js","_app/immutable/chunks/IHki7fMi.js"];
export const stylesheets = ["_app/immutable/assets/GameCard.Diu31AHu.css"];
export const fonts = [];
