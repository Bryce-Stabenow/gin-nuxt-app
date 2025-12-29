<template>
  <div class="min-h-screen bg-gradient-to-br from-purple-500 to-purple-700 px-4 py-20">
    <div class="max-w-4xl mx-auto">
      <div class="bg-white rounded-xl shadow-2xl p-10 mb-6">
        <h1 class="text-3xl font-bold text-gray-900 mb-8 text-center">Dashboard</h1>
        <div v-if="isLoading" class="text-center text-gray-600 text-base py-5">
          Loading...
        </div>
        <div v-else-if="isAuthenticated && user" class="mt-5">
          <p class="text-gray-900 text-lg leading-relaxed text-center mb-6">
            Welcome, <strong class="text-purple-600 font-semibold">{{ user.email }}</strong>
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

      <!-- Lists Section -->
      <div v-if="isAuthenticated && !isLoading" class="bg-white rounded-xl shadow-2xl p-10">
        <div class="flex justify-between items-center mb-6">
          <h2 class="text-2xl font-bold text-gray-900">My Lists</h2>
          <NuxtLink
            to="/lists/new"
            class="inline-block px-6 py-3 bg-gradient-to-r from-purple-500 to-purple-700 text-white rounded-lg font-semibold no-underline transition-transform hover:-translate-y-0.5 hover:shadow-lg"
          >
            +
          </NuxtLink>
        </div>

        <div v-if="listsLoading" class="text-center text-gray-600 text-base py-5">
          Loading lists...
        </div>
        <div v-else-if="listsError" class="text-center text-red-800 py-5">
          <p class="mb-5 text-base">Error: {{ listsError }}</p>
        </div>
        <div v-else-if="lists.length === 0" class="text-center text-gray-600 py-10">
          <p class="mb-5 text-base">You don't have any lists yet.</p>
          <NuxtLink
            to="/lists/new"
            class="inline-block px-6 py-3 bg-gradient-to-r from-purple-500 to-purple-700 text-white rounded-lg font-semibold no-underline transition-transform hover:-translate-y-0.5 hover:shadow-lg"
          >
            Create Your First List
          </NuxtLink>
        </div>
        <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <NuxtLink
            v-for="list in lists"
            :key="list.id"
            :to="`/lists/${list.id}`"
            class="block p-5 border-2 border-gray-200 rounded-lg hover:border-purple-500 hover:shadow-md transition-all cursor-pointer no-underline"
          >
            <h3 class="text-xl font-semibold text-gray-900 mb-2">{{ list.name }}</h3>
            <p v-if="list.description" class="text-gray-600 text-sm mb-3 line-clamp-2">{{ list.description }}</p>
            <div class="flex justify-between items-center text-sm text-gray-500">
              <span>{{ list.items.length }} item{{ list.items.length !== 1 ? 's' : '' }}</span>
              <span>{{ new Date(list.created_at).toLocaleDateString() }}</span>
            </div>
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const { isAuthenticated, user, isLoading, checkAuth } = useAuth()
const { getLists } = useLists()

const lists = ref<any[]>([])
const listsLoading = ref(false)
const listsError = ref<string | null>(null)

// Check authentication on page load
onMounted(async () => {
  await checkAuth()
  
  if (isAuthenticated.value) {
    await loadLists()
  }
})

const loadLists = async () => {
  listsLoading.value = true
  listsError.value = null
  
  try {
    lists.value = await getLists()
  } catch (error: any) {
    listsError.value = error.data?.error || error.message || 'Failed to load lists'
  } finally {
    listsLoading.value = false
  }
}
</script>

