<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-purple-500 to-purple-700 px-4 py-20">
    <div class="bg-white rounded-xl shadow-2xl p-10 w-full max-w-2xl">
      <h1 class="text-3xl font-bold text-gray-900 mb-8 text-center">Dashboard</h1>
      <div v-if="isLoading" class="text-center text-gray-600 text-base py-5">
        Loading...
      </div>
      <div v-else-if="isAuthenticated && user" class="mt-5">
        <p class="text-gray-900 text-lg leading-relaxed text-center">
          Currently, <strong class="text-purple-600 font-semibold">{{ user.email }}</strong> is the email address of the user.
        </p>
      </div>
      <div v-else class="text-center text-red-800 py-5">
        <p class="mb-5 text-base">You are not authenticated. Please sign in.</p>
        <NuxtLink
          to="/signin"
          class="inline-block px-6 py-3 bg-gradient-to-r from-purple-500 to-purple-700 text-white rounded-lg font-semibold no-underline transition-transform hover:-translate-y-0.5 hover:shadow-lg"
        >
          Sign In
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const { isAuthenticated, user, isLoading, checkAuth } = useAuth()

// Check authentication on page load
onMounted(async () => {
  await checkAuth()
})
</script>

