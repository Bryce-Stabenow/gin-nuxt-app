<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-purple-500 to-purple-700 px-4">
    <div class="text-center">
      <h1 class="text-4xl font-bold text-white mb-8">Welcome to GrocerMe</h1>
      <p class="text-white text-lg mb-8">Your personal grocery management solution</p>
    </div>
  </div>
</template>

<script setup lang="ts">
const { checkAuth } = useAuth()

// Check authentication on page load
onMounted(async () => {
  await checkAuth()
})
</script>