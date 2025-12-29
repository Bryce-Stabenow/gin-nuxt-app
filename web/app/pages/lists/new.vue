<template>
  <div class="min-h-screen flex items-start justify-center bg-gradient-to-br from-purple-500 to-purple-700 px-4 py-10">
    <div class="bg-white rounded-xl shadow-2xl py-10 px-4 w-full max-w-md">
      <h1 class="text-3xl font-bold text-gray-900 mb-2">Create New List</h1>
      <p class="text-gray-600 text-sm mb-8">Create a new grocery list to get started</p>
      <form @submit.prevent="handleSubmit">
        <div class="mb-5">
          <label for="name" class="block text-gray-900 mb-2 font-medium text-sm">List Name *</label>
          <input
            type="text"
            id="name"
            v-model="name"
            required
            placeholder="e.g., Weekly Groceries"
            class="w-full px-3 py-3 border-2 border-gray-200 rounded-lg text-base transition-colors focus:outline-none focus:border-purple-500"
          />
        </div>
        <div class="mb-5">
          <label for="description" class="block text-gray-900 mb-2 font-medium text-sm">Description</label>
          <textarea
            id="description"
            v-model="description"
            rows="3"
            placeholder="Optional description for your list"
            class="w-full px-3 py-3 border-2 border-gray-200 rounded-lg text-base transition-colors focus:outline-none focus:border-purple-500 resize-none"
          ></textarea>
        </div>
        <button
          type="submit"
          :disabled="isSubmitting"
          class="w-full py-3.5 bg-gradient-to-r from-purple-500 to-purple-700 text-white rounded-lg text-base font-semibold cursor-pointer transition-transform hover:-translate-y-0.5 hover:shadow-lg active:translate-y-0 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {{ isSubmitting ? 'Creating...' : 'Create List' }}
        </button>
      </form>
      <div v-if="message" :class="[
        'mt-5 p-3 rounded-lg',
        messageType === 'success' ? 'bg-green-100 text-green-800 border border-green-200' : 'bg-red-100 text-red-800 border border-red-200'
      ]">
        {{ message }}
      </div>
      <div class="text-center mt-5">
        <NuxtLink to="/dashboard" class="text-purple-600 no-underline font-medium hover:underline text-sm">Back to Dashboard</NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const { createList } = useLists()

const name = ref('')
const description = ref('')
const message = ref('')
const messageType = ref<'success' | 'error'>('success')
const isSubmitting = ref(false)

const handleSubmit = async () => {
  message.value = ''
  isSubmitting.value = true
  
  try {
    const list = await createList(name.value, description.value || undefined)
    
    messageType.value = 'success'
    message.value = 'List created successfully!'
    
    // Redirect to the new list page
    await navigateTo(`/lists/${list.id}`)
  } catch (error: any) {
    messageType.value = 'error'
    message.value = 'Error: ' + (error.data?.error || error.message || 'Failed to create list')
    isSubmitting.value = false
  }
}
</script>

