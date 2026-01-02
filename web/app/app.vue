<template>
  <NuxtLayout>
    <NuxtPage />
    <InstallPrompt />
  </NuxtLayout>
</template>

<script setup lang="ts">
useHead({
  link: [
    { rel: 'manifest', href: '/manifest.webmanifest' },
    { rel: 'apple-touch-icon', href: '/apple-touch-icon.png' }
  ],
  meta: [
    { name: 'theme-color', content: '#9333ea' },
    { name: 'apple-mobile-web-app-capable', content: 'yes' },
    { name: 'apple-mobile-web-app-status-bar-style', content: 'default' }
  ]
})
</script>
