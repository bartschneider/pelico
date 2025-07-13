import type { Handle } from '@sveltejs/kit';

const API_BASE_URL = process.env.API_BASE_URL || 'http://pelico:8080';

export const handle: Handle = async ({ event, resolve }) => {
  // Proxy API requests to the backend service
  if (event.url.pathname.startsWith('/api/')) {
    const apiUrl = `${API_BASE_URL}${event.url.pathname}${event.url.search}`;
    
    try {
      const response = await fetch(apiUrl, {
        method: event.request.method,
        headers: {
          ...Object.fromEntries(event.request.headers),
          host: new URL(API_BASE_URL).host,
        },
        body: event.request.method !== 'GET' && event.request.method !== 'HEAD' 
          ? await event.request.text() 
          : undefined,
      });

      const responseHeaders = new Headers();
      response.headers.forEach((value, key) => {
        // Skip headers that shouldn't be forwarded
        if (!['content-encoding', 'content-length', 'transfer-encoding'].includes(key.toLowerCase())) {
          responseHeaders.set(key, value);
        }
      });

      return new Response(response.body, {
        status: response.status,
        statusText: response.statusText,
        headers: responseHeaders,
      });
    } catch (error) {
      console.error('API proxy error:', error);
      return new Response(JSON.stringify({ error: 'Backend service unavailable' }), {
        status: 503,
        headers: { 'Content-Type': 'application/json' },
      });
    }
  }

  return resolve(event);
};