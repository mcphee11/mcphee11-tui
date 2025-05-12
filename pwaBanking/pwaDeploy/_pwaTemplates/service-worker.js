importScripts('https://storage.googleapis.com/workbox-cdn/releases/6.4.1/workbox-sw.js')

workbox.setConfig({
  debug: false,
})

workbox.routing.setDefaultHandler(new workbox.strategies.NetworkFirst())

//Google Fonts Stuff
const maxAgeSeconds = 60 * 60 * 24 * 365
const maxEntries = 30

// Cache the underlying font files with a cache-first strategy for 1 year.
workbox.routing.registerRoute(
  ({ url }) => url.origin === 'https://fonts.gstatic.com',
  new workbox.strategies.CacheFirst({
    cacheName: 'google-fonts-webfonts',
    plugins: [
      new workbox.cacheableResponse.CacheableResponsePlugin({
        statuses: [0, 200],
      }),
      new workbox.expiration.ExpirationPlugin({
        maxAgeSeconds,
        maxEntries,
      }),
    ],
  })
)

self.addEventListener('message', (event) => {
  if (event.data) {
    console.log(event.data)
    event.ports[0].postMessage({ reply: 'received from sw' })
    returnMessage = event.ports[0]
  }
})

workbox.precaching.precacheAndRoute([])
self.skipWaiting()
workbox.core.clientsClaim()

// --------------------- PUSH NOTIFICATIONS ------------------

self.addEventListener('push', (event) => {
  const data = event.data.json()
  const options = {
    body: data.body,
    icon: './LOGO', // Replace with your icon path
    data: {
      action: 'notification-click',
    },
  }
  event.waitUntil(self.registration.showNotification(data.title, options))
})

self.addEventListener('notificationclick', (event) => {
  event.notification.close()
  clients.matchAll({ type: 'window' }).then((clientList) => {
    clientList.forEach((client) => {
      client.postMessage({
        type: 'notification-clicked',
        message: 'message',
      })
    })
  })
})
