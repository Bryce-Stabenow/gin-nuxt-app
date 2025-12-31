<template>
  <PageContainer>
    <div class="flex items-center justify-center">
      <div class="bg-white rounded-xl shadow-2xl py-10 px-4 max-w-md w-full">
      <div v-if="isLoading" class="text-center">
        <div class="mb-4">
          <Icon name="svg-spinners:ring-resize" class="h-12 w-12 text-purple-600 mx-auto" />
        </div>
        <p class="text-gray-700 text-lg">{{ loadingMessage }}</p>
      </div>
      <div v-else-if="error" class="text-center">
        <div class="mb-4">
          <Icon name="heroicons:exclamation-circle" class="h-12 w-12 text-red-600 mx-auto" />
        </div>
        <p class="text-red-800 text-lg mb-4">{{ error }}</p>
        <NuxtLink
          to="/dashboard"
          class="inline-block px-6 py-3 bg-gradient-to-r from-purple-500 to-purple-700 text-white rounded-lg font-semibold no-underline transition-transform hover:-translate-y-0.5 hover:shadow-lg"
        >
          Go to Dashboard
        </NuxtLink>
      </div>
      <div v-else-if="success" class="text-center">
        <div class="mb-4">
          <Icon name="heroicons:check-circle" class="h-12 w-12 text-green-600 mx-auto" />
        </div>
        <p class="text-gray-700 text-lg mb-4">Successfully joined the list!</p>
        <p class="text-gray-500 text-sm mb-6">Redirecting to the list...</p>
      </div>
      </div>
    </div>
  </PageContainer>
</template>

<script setup lang="ts">
const route = useRoute()
const router = useRouter()
const { isAuthenticated, checkAuth } = useAuth()
const { shareList } = useLists()

const isLoading = ref(true)
const error = ref<string | null>(null)
const success = ref(false)
const loadingMessage = ref('Processing...')

onMounted(async () => {
  const listId = route.params.id as string

  // Check if user is authenticated
  const authenticated = await checkAuth()

  if (!authenticated) {
    // Redirect to signup with redirect parameter
    await router.push(`/signup?redirect=/lists/share/${listId}`)
    return
  }

  // User is authenticated, proceed with sharing
  try {
    loadingMessage.value = 'Adding you to the list...'
    await shareList(listId)
    success.value = true
    loadingMessage.value = 'Redirecting...'

    // Redirect to the list page after a short delay
    setTimeout(() => {
      router.push(`/lists/${listId}`)
    }, 1500)
  } catch (err: any) {
    if (err.statusCode === 401) {
      // Token expired or invalid, redirect to signin
      await router.push(`/signin?redirect=/lists/share/${listId}`)
    } else if (err.statusCode === 404) {
      error.value = 'List not found'
    } else if (err.statusCode === 400 && err.data?.error?.includes('owner')) {
      error.value = 'You are already the owner of this list'
    } else {
      error.value = err.data?.error || err.message || 'Failed to join list'
    }
    isLoading.value = false
  }
})
</script>

