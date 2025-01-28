// Service Worker

// Cache name
const CACHE_NAME = 'site-cache';

// Add any file you want to display offline here:
const urlsToCache = [
  'offline.html',
  'sw.js',
];

// Install event
self.addEventListener('install', (event) => {
  // Perform installation steps
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then((cache) => {
        console.log('Opened cache');
        return cache.addAll(urlsToCache);
      })
  );
});

// Fetch event
self.addEventListener('fetch', (event) => {
  event.respondWith(
    caches.match(event.request)
      .then((response) => {
        // Cache hit - return response
        if (response) {
          return response;
        }
        // Make a network request
        return fetch(event.request);
      })
  );
});

// Activate event
self.addEventListener('activate', (event) => {
  const cacheWhitelist = ['site-cache'];
  event.waitUntil(
    caches.keys()
      .then((cacheNames) => {
        return Promise.all(
          cacheNames.map((cacheName) => {
            if (!cacheWhitelist.includes(cacheName)) {
              return caches.delete(cacheName);
            }
          })
        );
      })
  );
});
