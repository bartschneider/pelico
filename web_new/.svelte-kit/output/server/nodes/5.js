

export const index = 5;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/wishlist/_page.svelte.js')).default;
export const imports = ["_app/immutable/nodes/5.CfpYEaOu.js","_app/immutable/chunks/BXdaxRue.js","_app/immutable/chunks/CO2FEBnq.js","_app/immutable/chunks/IHki7fMi.js"];
export const stylesheets = ["_app/immutable/assets/GameCard.Diu31AHu.css"];
export const fonts = [];
