<template>
  <div class="min-h-screen flex flex-col">
    <header class="bg-white shadow-sm border-b">
      <nav class="container mx-auto px-4 py-4">
        <div class="flex items-center justify-between">
          <NuxtLink to="/" class="text-2xl font-bold text-gray-900 hover:text-purple-600 transition-colors">
            GrocerMe
          </NuxtLink>
          <div class="flex gap-4 items-center">
            <template v-if="isAuthenticated">
              <NuxtLink to="/dashboard" class="text-gray-700 hover:text-purple-600 hover:underline transition-colors">
                Dashboard
              </NuxtLink>
              <button
                @click="handleSignOut"
                class="text-gray-700 hover:text-purple-600 hover:underline transition-colors"
              >
                Sign Out
              </button>
            </template>
            <template v-else>
              <NuxtLink to="/signin" class="text-gray-700 hover:text-purple-600 hover:underline transition-colors">
                Sign In
              </NuxtLink>
              <NuxtLink to="/signup" class="text-gray-700 hover:text-purple-600 hover:underline transition-colors">
                Sign Up
              </NuxtLink>
            </template>
          </div>
        </div>
      </nav>
    </header>
    <main class="flex-1">
      <slot />
    </main>
  </div>
</template>

<script setup lang="ts">
const { isAuthenticated, clearAuth } = useAuth()

const handleSignOut = async () => {
  // Clear auth state
  clearAuth()
  
  // Navigate to home page
  await navigateTo('/')
}
</script>

