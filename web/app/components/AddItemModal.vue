<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="isOpen"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50"
        @click.self="close"
      >
        <div
          class="bg-white rounded-xl shadow-2xl p-8 max-w-md w-full mx-4 transform transition-all"
        >
          <div class="flex justify-between items-center mb-6">
            <h2 class="text-2xl font-bold text-gray-900">Add Item</h2>
          </div>

          <form @submit.prevent="handleSubmit" class="space-y-6">
            <div>
              <label
                for="item-name"
                class="block text-sm font-medium text-gray-700 mb-2"
              >
                Name*
              </label>
              <input
                id="item-name"
                v-model="form.name"
                type="text"
                required
                class="w-full px-4 py-2 border-2 border-gray-300 rounded-lg focus:outline-none focus:border-purple-500 transition-colors"
                placeholder="Enter item name"
                ref="nameInput"
              />
            </div>

            <div>
              <label
                for="item-quantity"
                class="block text-sm font-medium text-gray-700 mb-2"
              >
                Quantity
              </label>
              <input
                id="item-quantity"
                v-model.number="form.quantity"
                type="number"
                min="1"
                class="w-full px-4 py-2 border-2 border-gray-300 rounded-lg focus:outline-none focus:border-purple-500 transition-colors"
                placeholder="1"
              />
            </div>

            <div>
              <label
                for="item-details"
                class="block text-sm font-medium text-gray-700 mb-2"
              >
                Details
              </label>
              <textarea
                id="item-details"
                v-model="form.details"
                maxlength="512"
                rows="3"
                class="w-full px-4 py-2 border-2 border-gray-300 rounded-lg focus:outline-none focus:border-purple-500 transition-colors resize-none"
                placeholder="Add any additional details (optional)"
              />
              <div class="text-xs text-gray-500 mt-1 text-right">
                {{ (form.details || '').length }}/512
              </div>
            </div>

            <div v-if="error" class="text-red-600 text-sm">
              {{ error }}
            </div>

            <div class="flex gap-4 justify-center">
              <button
                type="button"
                @click="close"
                class="px-4 py-2 text-gray-700 border-2 border-gray-300 rounded-lg font-medium hover:bg-gray-50 transition-colors"
              >
                Cancel
              </button>
              <button
                type="submit"
                :disabled="isSubmitting"
                class="flex-1 px-4 py-2 bg-gradient-to-r from-purple-500 to-purple-700 text-white rounded-lg font-medium hover:shadow-lg transition-all disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <span v-if="isSubmitting">Adding...</span>
                <span v-else>Add Item</span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
interface Props {
  isOpen: boolean
}

interface Emits {
  (e: 'close'): void
  (e: 'item-added', item: any): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const { addListItem } = useLists()

const form = ref({
  name: '',
  quantity: 1,
  details: '',
})

const error = ref<string | null>(null)
const isSubmitting = ref(false)
const nameInput = ref<HTMLInputElement | null>(null)

const close = () => {
  emit('close')
}

const handleSubmit = async () => {
  if (!form.value.name.trim()) {
    error.value = 'Item name is required'
    return
  }

  isSubmitting.value = true
  error.value = null

  try {
    const listId = useRoute().params.id as string
    const updatedList = await addListItem(listId, {
      name: form.value.name.trim(),
      quantity: form.value.quantity || 1,
      details: form.value.details?.trim() || undefined,
    })
    
    emit('item-added', updatedList)
    
    // Reset form
    form.value = {
      name: '',
      quantity: 1,
      details: '',
    }
    close()
  } catch (err: any) {
    error.value = err.data?.error || err.message || 'Failed to add item'
  } finally {
    isSubmitting.value = false
  }
}

// Focus name input when modal opens
watch(() => props.isOpen, (isOpen) => {
  if (isOpen) {
    nextTick(() => {
      nameInput.value?.focus()
    })
  } else {
    // Reset form when closing
    form.value = {
      name: '',
      quantity: 1,
      details: '',
    }
    error.value = null
  }
})
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.1s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active > div,
.modal-leave-active > div {
  transition: transform 0.1s ease, opacity 0.1s ease;
}

.modal-enter-from > div,
.modal-leave-to > div {
  transform: scale(0.9);
  opacity: 0;
}
</style>

