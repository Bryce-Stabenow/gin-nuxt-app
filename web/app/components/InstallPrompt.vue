<template>
  <div v-if="showPrompt" class="fixed bottom-4 left-4 right-4 bg-white shadow-lg rounded-lg p-4 flex items-center justify-between z-50 border border-gray-200">
    <div class="flex-1">
      <h3 class="font-semibold text-gray-900">Install GrocerMe</h3>
      <p class="text-sm text-gray-600">Install our app for a better experience</p>
    </div>
    <div class="flex gap-2 ml-4">
      <button 
        @click="dismiss" 
        class="px-3 py-1 text-sm text-gray-600 hover:text-gray-900 transition-colors"
      >
        Later
      </button>
      <button 
        @click="install" 
        class="px-4 py-2 text-sm bg-green-500 text-white rounded hover:bg-green-600 transition-colors"
      >
        Install
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
const showPrompt = ref(false)
let deferredPrompt: any = null

onMounted(() => {
  window.addEventListener('beforeinstallprompt', (e) => {
    e.preventDefault()
    deferredPrompt = e
    showPrompt.value = true
  })
  
  window.addEventListener('appinstalled', () => {
    showPrompt.value = false
    deferredPrompt = null
  })
})

const install = async () => {
  if (!deferredPrompt) return
  
  deferredPrompt.prompt()
  const { outcome } = await deferredPrompt.userChoice
  
  if (outcome === 'accepted') {
    showPrompt.value = false
  }
  
  deferredPrompt = null
}

const dismiss = () => {
  showPrompt.value = false
  deferredPrompt = null
}
</script>

