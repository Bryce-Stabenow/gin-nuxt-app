<template>
  <div class="min-h-screen bg-gradient-to-br from-purple-500 to-purple-700 px-4 py-20">
    <div class="max-w-4xl mx-auto">
      <div class="bg-white rounded-xl shadow-2xl p-10">
        <div v-if="isLoading" class="text-center text-gray-600 text-base py-10">
          Loading list...
        </div>
        <div v-else-if="error" class="text-center text-red-800 py-10">
          <p class="mb-5 text-base">{{ error }}</p>
          <NuxtLink
            to="/dashboard"
            class="inline-block px-6 py-3 bg-gradient-to-r from-purple-500 to-purple-700 text-white rounded-lg font-semibold no-underline transition-transform hover:-translate-y-0.5 hover:shadow-lg"
          >
            Back to Dashboard
          </NuxtLink>
        </div>
        <div v-else-if="list">
          <!-- Header -->
          <div class="mb-6">
            <div class="flex justify-between items-start mb-4">
              <div class="flex-1">
                <div class="flex items-center gap-2 mb-2">
                  <h1 v-if="!isEditingName" class="text-3xl font-bold text-gray-900">{{ list.name }}</h1>
                  <input
                    v-else
                    v-model="editingName"
                    @blur="saveName"
                    @keydown.enter="saveName"
                    @keydown.esc="cancelEditName"
                    class="text-3xl font-bold text-gray-900 bg-transparent border-b-2 border-purple-500 focus:outline-none focus:border-purple-700 w-full"
                    ref="nameInput"
                  />
                  <button
                    v-if="!isEditingName"
                    @click="startEditName"
                    class="p-1 text-gray-400 hover:text-purple-600 transition-colors"
                    title="Edit list name"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                  </button>
                </div>
                <p v-if="list.description" class="text-gray-600 text-base">{{ list.description }}</p>
              </div>
              <NuxtLink
                to="/dashboard"
                class="px-4 py-2 text-purple-600 border-2 border-purple-600 rounded-lg font-medium no-underline hover:bg-purple-50 transition-colors ml-4"
              >
                Back
              </NuxtLink>
            </div>
            <div class="flex gap-4 text-sm text-gray-500">
              <span>Created: {{ new Date(list.created_at).toLocaleDateString() }}</span>
              <span>Updated: {{ new Date(list.updated_at).toLocaleDateString() }}</span>
            </div>
          </div>

          <!-- Items Section -->
          <div class="border-t border-gray-200 pt-6">
            <h2 class="text-xl font-semibold text-gray-900 mb-4">
              Items ({{ list.items.length }})
            </h2>
            <div v-if="list.items.length === 0" class="text-center text-gray-500 py-10 border-2 border-dashed border-gray-200 rounded-lg">
              <p>No items in this list yet.</p>
              <p class="text-sm mt-2">Items can be added when the update feature is implemented.</p>
            </div>
            <div v-else class="space-y-3">
              <div
                v-for="(item, index) in list.items"
                :key="index"
                class="flex items-center gap-4 p-4 border-2 border-gray-200 rounded-lg hover:border-purple-300 transition-colors"
                :class="{ 'bg-gray-50': item.checked }"
              >
                <input
                  type="checkbox"
                  :checked="item.checked"
                  disabled
                  class="w-5 h-5 text-purple-600 border-gray-300 rounded focus:ring-purple-500"
                />
                <div class="flex-1">
                  <div class="flex items-center gap-2">
                    <span class="font-medium text-gray-900" :class="{ 'line-through text-gray-500': item.checked }">
                      {{ item.name }}
                    </span>
                  </div>
                  <div class="text-sm text-gray-500 mt-1">
                    <span v-if="item.quantity > 0">{{ item.quantity }}</span>
                    <span v-if="item.quantity > 0 && item.unit">{{ item.unit }}</span>
                    <span v-if="item.quantity > 0 && item.unit" class="mx-1">•</span>
                    <span>Added {{ new Date(item.added_at).toLocaleDateString() }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Shared With Section -->
          <div v-if="list.shared_with.length > 0" class="border-t border-gray-200 pt-6 mt-6">
            <h2 class="text-xl font-semibold text-gray-900 mb-4">
              Shared With ({{ list.shared_with.length }})
            </h2>
            <div class="flex flex-wrap gap-2">
              <span
                v-for="userId in list.shared_with"
                :key="userId"
                class="px-3 py-1 bg-purple-100 text-purple-700 rounded-full text-sm"
              >
                {{ userId }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const { getList, updateList } = useLists()

const list = ref<any>(null)
const isLoading = ref(true)
const error = ref<string | null>(null)
const isEditingName = ref(false)
const editingName = ref('')
const nameInput = ref<HTMLInputElement | null>(null)
const isSaving = ref(false)

onMounted(async () => {
  await loadList()
})

const loadList = async () => {
  isLoading.value = true
  error.value = null
  
  try {
    const listId = route.params.id as string
    list.value = await getList(listId)
  } catch (err: any) {
    if (err.statusCode === 404) {
      error.value = 'List not found'
    } else if (err.statusCode === 403) {
      error.value = 'You do not have access to this list'
    } else {
      error.value = err.data?.error || err.message || 'Failed to load list'
    }
  } finally {
    isLoading.value = false
  }
}

const startEditName = () => {
  if (!list.value) return
  editingName.value = list.value.name
  isEditingName.value = true
  nextTick(() => {
    nameInput.value?.focus()
    nameInput.value?.select()
  })
}

const cancelEditName = () => {
  isEditingName.value = false
  editingName.value = ''
}

const saveName = async () => {
  if (!list.value || isSaving.value) return
  
  const trimmedName = editingName.value.trim()
  if (!trimmedName) {
    cancelEditName()
    return
  }
  
  if (trimmedName === list.value.name) {
    cancelEditName()
    return
  }
  
  isSaving.value = true
  try {
    const listId = route.params.id as string
    const updatedList = await updateList(listId, { name: trimmedName })
    list.value = updatedList
    isEditingName.value = false
    editingName.value = ''
  } catch (err: any) {
    error.value = err.data?.error || err.message || 'Failed to update list name'
    // Keep editing mode on error so user can retry
  } finally {
    isSaving.value = false
  }
}
</script>

