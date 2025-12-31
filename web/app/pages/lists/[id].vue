<template>
  <PageContainer>
    <div class="max-w-4xl mx-auto">
      <div class="bg-white rounded-xl shadow-2xl py-10 px-4">
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
        <div v-else-if="list" class="relative">
          <!-- Header -->
          <div class="sticky top-0 bg-white z-10 py-2">
            <div class="flex items-center justify-between">
              <h1
                v-if="!isEditingName"
                class="text-3xl font-bold text-gray-900"
              >
                {{ list.name }}
              </h1>
              <input
                v-else
                v-model="editingName"
                @blur="saveName"
                @keydown.enter="saveName"
                @keydown.esc="cancelEditName"
                class="text-3xl font-bold text-gray-900 bg-transparent border-b-2 border-purple-500 focus:outline-none focus:border-purple-700 w-full"
                ref="nameInput"
              />
              <div class="flex items-center gap-2">
                <button
                  v-if="!isEditingName"
                  @click="handleShareList"
                  class="p-1 text-gray-400 hover:text-purple-600 transition-colors"
                  title="Share list"
                >
                  <Icon name="heroicons:share" class="h-5 w-5" />
                </button>
                <button
                  v-if="!isEditingName"
                  @click="startEditName"
                  class="p-1 text-gray-400 hover:text-purple-600 transition-colors"
                  title="Edit list name"
                >
                  <Icon name="heroicons:pencil" class="h-5 w-5" />
                </button>
                <button
                  v-if="!isEditingName && isListOwner"
                  @click="handleDeleteList"
                  :disabled="isDeletingList"
                  class="p-1 text-gray-400 hover:text-red-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                  title="Delete list"
                >
                  <Icon
                    v-if="isDeletingList"
                    name="svg-spinners:ring-resize"
                    class="h-5 w-5"
                  />
                  <Icon v-else name="heroicons:trash" class="h-5 w-5" />
                </button>
              </div>
            </div>
          </div>

          <!-- Share Notification -->
          <div
            v-if="shareNotification"
            class="mb-4 p-3 bg-green-100 border border-green-300 text-green-800 rounded-lg text-sm"
          >
            {{ shareNotification }}
          </div>

          <div class="mb-6">
            <div
              v-if="list.description"
              class="flex justify-between items-start mb-4"
            >
              <div class="flex-1">
                <p class="text-gray-600 text-base">
                  {{ list.description }}
                </p>
              </div>
            </div>
            <div class="flex gap-4 text-sm text-gray-500">
              <span
                >Created:
                {{ new Date(list.created_at).toLocaleDateString() }}</span
              >
              <span
                >Updated:
                {{ new Date(list.updated_at).toLocaleDateString() }}</span
              >
            </div>
          </div>

          <!-- Items Section -->
          <div class="border-t border-gray-200 pt-6">
            <div class="flex justify-between items-center mb-4">
              <h2
                v-if="!isSearchOpen"
                class="text-xl font-semibold text-gray-900"
              >
                Items ({{ list.items.length }})
              </h2>
              <input
                v-else
                v-model="searchQuery"
                @keydown.esc="closeSearch"
                class="text-xl font-semibold text-gray-900 bg-transparent border-b-2 border-purple-500 focus:outline-none focus:border-purple-700 flex-1 mr-2"
                placeholder="Search items..."
                ref="searchInput"
              />
              <div class="flex items-center gap-2">
                <button
                  v-if="isSearchOpen"
                  @click="closeSearch"
                  class="p-1 text-gray-400 hover:text-purple-600 transition-colors"
                  title="Close search"
                >
                  <Icon name="heroicons:x-mark" class="h-5 w-5" />
                </button>
                <button
                  v-if="!isSearchOpen"
                  @click="openSearch"
                  class="p-1 text-gray-400 hover:text-purple-600 transition-colors"
                  title="Search items"
                >
                  <Icon name="heroicons:magnifying-glass" class="h-5 w-5" />
                </button>
              </div>
            </div>
            <div
              v-if="list.items.length === 0 && !showAddForm"
              class="text-center text-gray-500 py-10 border-2 border-dashed border-gray-200 rounded-lg"
            >
              <p>No items in this list yet.</p>
              <p class="text-sm mt-2 text-purple-600">
                Click the button below to get started.
              </p>
            </div>
            <div v-else class="space-y-5">
              <ListItem
                v-for="(sortedItem, displayIndex) in sortedItems"
                :key="sortedItem.originalIndex"
                :item="sortedItem.item"
                :original-index="sortedItem.originalIndex"
                @toggle="toggleItemChecked"
                @change="handleItemCheckedChange"
                @click="openEditModal"
              />
            </div>
            <!-- Inline Add Item Form -->
            <div
              v-if="showAddForm"
              class="mt-6 p-4 border-2 border-purple-300 rounded-lg bg-purple-50"
            >
              <form @submit.prevent="handleAddItem" class="space-y-4">
                <div>
                  <label
                    for="add-item-name"
                    class="block text-sm font-medium text-gray-700 mb-2"
                  >
                    Name*
                  </label>
                  <input
                    id="add-item-name"
                    v-model="addForm.name"
                    type="text"
                    required
                    class="w-full px-4 py-2 border-2 border-gray-300 rounded-lg focus:outline-none focus:border-purple-500 transition-colors"
                    placeholder="Enter item name"
                    ref="addNameInput"
                  />
                </div>

                <div>
                  <label
                    for="add-item-quantity"
                    class="block text-sm font-medium text-gray-700 mb-2"
                  >
                    Quantity
                  </label>
                  <input
                    id="add-item-quantity"
                    v-model.number="addForm.quantity"
                    type="number"
                    min="1"
                    class="w-full px-4 py-2 border-2 border-gray-300 rounded-lg focus:outline-none focus:border-purple-500 transition-colors"
                    placeholder="1"
                  />
                </div>

                <div>
                  <label
                    for="add-item-details"
                    class="block text-sm font-medium text-gray-700 mb-2"
                  >
                    Details
                  </label>
                  <textarea
                    id="add-item-details"
                    v-model="addForm.details"
                    maxlength="512"
                    rows="3"
                    class="w-full px-4 py-2 border-2 border-gray-300 rounded-lg focus:outline-none focus:border-purple-500 transition-colors resize-none"
                    placeholder="Add any additional details (optional)"
                  />
                  <div class="text-xs text-gray-500 mt-1 text-right">
                    {{ (addForm.details || "").length }}/512
                  </div>
                </div>

                <div v-if="addError" class="text-red-600 text-sm">
                  {{ addError }}
                </div>

                <div class="flex gap-4 justify-center">
                  <button
                    type="button"
                    @click="cancelAddForm"
                    class="px-4 py-2 text-gray-700 border-2 border-gray-300 rounded-lg font-medium hover:bg-gray-50 transition-colors"
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    :disabled="isAdding"
                    class="flex-1 px-4 py-2 bg-gradient-to-r from-purple-500 to-purple-700 text-white rounded-lg font-medium hover:shadow-lg transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    <span v-if="isAdding">Adding...</span>
                    <span v-else>Add Item</span>
                  </button>
                </div>
              </form>
            </div>
            <div v-else class="flex justify-center pt-6">
              <button
                @click="showAddForm = true"
                class="px-4 py-2 bg-gradient-to-r from-purple-500 to-purple-700 text-white rounded-lg font-medium hover:shadow-lg transition-all"
              >
                + Add
              </button>
            </div>
          </div>

          <!-- Shared With Section -->
          <div
            v-if="list.shared_with.length > 0"
            class="border-t border-gray-200 pt-6 mt-6"
          >
            <h2 class="text-xl font-semibold text-gray-900 mb-4">
              Shared With ({{ list.shared_with.length }})
            </h2>
            <div class="flex flex-wrap gap-2">
              <span
                v-for="sharedUser in list.shared_with"
                :key="sharedUser.id"
                class="px-3 py-1 bg-purple-100 text-purple-700 rounded-full text-sm"
              >
                {{ sharedUser.email || sharedUser.id }}
              </span>
            </div>
          </div>
        </div>
      </div>
      <div class="flex justify-center mt-10">
        <NuxtLink
          to="/dashboard"
          class="px-4 py-2 text-white border-2 border-white rounded-lg font-medium no-underline hover:bg-white hover:text-purple-500 transition-colors ml-4"
        >
          Back to Dashboard
        </NuxtLink>
      </div>
    </div>

    <!-- Edit Item Modal -->
    <EditItemModal
      :is-open="isEditModalOpen"
      :item="editingItem"
      :item-index="editingItemIndex"
      @close="closeEditModal"
      @item-updated="handleItemUpdated"
      @item-deleted="handleItemDeleted"
    />
  </PageContainer>
</template>

<script setup lang="ts">
import confetti from "canvas-confetti";

const route = useRoute();
const router = useRouter();
const { getList, updateList, updateListItemChecked, addListItem, deleteList } =
  useLists();
const { user } = useAuth();

useHead({
  title: 'GrocerMe | List',
  meta: [
    {
      name: 'description',
      content: 'View and manage your grocery list on GrocerMe.'
    },
    {
      property: 'og:title',
      content: 'GrocerMe | List'
    },
    {
      property: 'og:description',
      content: 'View and manage your grocery list on GrocerMe.'
    },
    {
      property: 'og:type',
      content: 'website'
    },
    {
      name: 'robots',
      content: 'noindex, nofollow'
    }
  ]
});

const list = ref<any>(null);
const isLoading = ref(true);
const error = ref<string | null>(null);
const isEditingName = ref(false);
const editingName = ref("");
const nameInput = ref<HTMLInputElement | null>(null);
const isSaving = ref(false);
const showAddForm = ref(false);
const isEditModalOpen = ref(false);
const editingItem = ref<any>(null);
const editingItemIndex = ref<number | null>(null);
const wasAllChecked = ref(false);
const addNameInput = ref<HTMLInputElement | null>(null);
const isAdding = ref(false);
const addError = ref<string | null>(null);
const isSearchOpen = ref(false);
const searchQuery = ref("");
const searchInput = ref<HTMLInputElement | null>(null);
const isDeletingList = ref(false);
const shareNotification = ref<string | null>(null);

const addForm = ref({
  name: "",
  quantity: 1,
  details: "",
});

// Computed property to sort items: unchecked items first, checked items at the bottom
// Also filters by search query if search is active
const sortedItems = computed(() => {
  if (!list.value || !list.value.items) {
    return [];
  }

  let items = list.value.items.map((item: any, originalIndex: number) => ({
    item,
    originalIndex,
  }));

  // Filter by search query if search is active
  if (isSearchOpen.value && searchQuery.value.trim()) {
    const query = searchQuery.value.trim().toLowerCase();
    items = items.filter(({ item }: any) =>
      item.name.toLowerCase().includes(query)
    );
  }

  return items.sort((a: any, b: any) => {
    // Unchecked items (false) come before checked items (true)
    if (a.item.checked === b.item.checked) {
      // If both have the same checked state, maintain original order
      return a.originalIndex - b.originalIndex;
    }
    return a.item.checked ? 1 : -1;
  });
});

// Check if current user is the list owner
const isListOwner = computed(() => {
  if (!list.value || !user.value) return false;
  return list.value.user_id === user.value.id;
});

// Confetti functions (defined early so they can be used in loadList)
const checkAllItemsChecked = (): boolean => {
  if (!list.value || !list.value.items || list.value.items.length === 0) {
    return false;
  }
  return list.value.items.every((item: any) => item.checked);
};

const triggerConfetti = () => {
  const duration = 3000;
  const animationEnd = Date.now() + duration;
  const defaults = { startVelocity: 30, spread: 360, ticks: 60, zIndex: 0 };

  function randomInRange(min: number, max: number) {
    return Math.random() * (max - min) + min;
  }

  const interval: any = setInterval(function () {
    const timeLeft = animationEnd - Date.now();

    if (timeLeft <= 0) {
      return clearInterval(interval);
    }

    const particleCount = 50 * (timeLeft / duration);
    confetti({
      ...defaults,
      particleCount,
      origin: { x: randomInRange(0.1, 0.3), y: Math.random() - 0.2 },
    });
    confetti({
      ...defaults,
      particleCount,
      origin: { x: randomInRange(0.7, 0.9), y: Math.random() - 0.2 },
    });
  }, 250);
};

const checkAndTriggerConfetti = () => {
  if (!list.value) return;

  const allChecked = checkAllItemsChecked();

  // Only trigger confetti when transitioning from "not all checked" to "all checked"
  if (allChecked && !wasAllChecked.value) {
    triggerConfetti();
  }

  wasAllChecked.value = allChecked;
};

const loadList = async () => {
  isLoading.value = true;
  error.value = null;

  try {
    const listId = route.params.id as string;
    list.value = await getList(listId);
    // Initialize the wasAllChecked state
    wasAllChecked.value = checkAllItemsChecked();
  } catch (err: any) {
    if (err.statusCode === 404) {
      error.value = "List not found";
    } else if (err.statusCode === 403) {
      error.value = "You do not have access to this list";
    } else {
      error.value = err.data?.error || err.message || "Failed to load list";
    }
  } finally {
    isLoading.value = false;
  }
};

const startEditName = () => {
  if (!list.value) return;
  editingName.value = list.value.name;
  isEditingName.value = true;
  nextTick(() => {
    nameInput.value?.focus();
    nameInput.value?.select();
  });
};

const cancelEditName = () => {
  isEditingName.value = false;
  editingName.value = "";
};

const saveName = async () => {
  if (!list.value || isSaving.value) return;

  const trimmedName = editingName.value.trim();
  if (!trimmedName) {
    cancelEditName();
    return;
  }

  if (trimmedName === list.value.name) {
    cancelEditName();
    return;
  }

  isSaving.value = true;
  try {
    const listId = route.params.id as string;
    const updatedList = await updateList(listId, { name: trimmedName });
    list.value = updatedList;
    isEditingName.value = false;
    editingName.value = "";
  } catch (err: any) {
    error.value =
      err.data?.error || err.message || "Failed to update list name";
    // Keep editing mode on error so user can retry
  } finally {
    isSaving.value = false;
  }
};

const handleAddItem = async () => {
  if (!addForm.value.name.trim()) {
    addError.value = "Item name is required";
    return;
  }

  isAdding.value = true;
  addError.value = null;

  try {
    const listId = route.params.id as string;
    const updatedList = await addListItem(listId, {
      name: addForm.value.name.trim(),
      quantity: addForm.value.quantity || 1,
      details: addForm.value.details?.trim() || undefined,
    });

    list.value = updatedList;
    checkAndTriggerConfetti();

    // Reset form and hide
    addForm.value = {
      name: "",
      quantity: 1,
      details: "",
    };
    showAddForm.value = false;
  } catch (err: any) {
    addError.value = err.data?.error || err.message || "Failed to add item";
  } finally {
    isAdding.value = false;
  }
};

const cancelAddForm = () => {
  showAddForm.value = false;
  addForm.value = {
    name: "",
    quantity: 1,
    details: "",
  };
  addError.value = null;
};

const openSearch = () => {
  isSearchOpen.value = true;
  nextTick(() => {
    searchInput.value?.focus();
  });
};

const closeSearch = () => {
  isSearchOpen.value = false;
  searchQuery.value = "";
};

// Watch for when add form is shown to focus input
watch(showAddForm, (isShown) => {
  if (isShown) {
    nextTick(() => {
      addNameInput.value?.focus();
    });
  }
});

const openEditModal = (index: number) => {
  if (!list.value || !list.value.items[index]) return;
  editingItem.value = { ...list.value.items[index] };
  editingItemIndex.value = index;
  isEditModalOpen.value = true;
};

const closeEditModal = () => {
  isEditModalOpen.value = false;
  editingItem.value = null;
  editingItemIndex.value = null;
};

const handleItemUpdated = (updatedList: any) => {
  list.value = updatedList;
};

const handleItemDeleted = (updatedList: any) => {
  list.value = updatedList;
  checkAndTriggerConfetti();
};

// Debounce timer for checkbox updates
const debounceTimers = new Map<number, ReturnType<typeof setTimeout>>();

const toggleItemChecked = (index: number) => {
  if (!list.value || !list.value.items[index]) return;

  const currentChecked = list.value.items[index].checked;
  const newChecked = !currentChecked;

  // Create a synthetic event to reuse existing handler
  const syntheticEvent = {
    target: { checked: newChecked },
  } as unknown as Event;

  handleItemCheckedChange(index, syntheticEvent);
};

const handleItemCheckedChange = async (index: number, event: Event) => {
  if (!list.value) return;

  const target = event.target as HTMLInputElement;
  const newChecked = target.checked;

  // Optimistically update the UI
  if (list.value.items[index]) {
    list.value.items[index].checked = newChecked;
  }

  // Check for confetti after optimistic update
  checkAndTriggerConfetti();

  // Clear existing debounce timer for this item
  const existingTimer = debounceTimers.get(index);
  if (existingTimer) {
    clearTimeout(existingTimer);
  }

  // Set new debounce timer
  const timer = setTimeout(async () => {
    try {
      const listId = route.params.id as string;
      const updatedList = await updateListItemChecked(
        listId,
        index,
        newChecked
      );
      // Update with server response to ensure sync
      list.value = updatedList;
      // Check for confetti after server response
      checkAndTriggerConfetti();
    } catch (err: any) {
      // Revert on error
      if (list.value.items[index]) {
        list.value.items[index].checked = !newChecked;
      }
      console.error("Failed to update item checked state:", err);
      // Re-check state after revert
      checkAndTriggerConfetti();
    } finally {
      debounceTimers.delete(index);
    }
  }, 500); // 500ms debounce

  debounceTimers.set(index, timer);
};

const handleShareList = async () => {
  if (!list.value) return;

  try {
    const shareUrl = `${window.location.origin}/lists/share/${list.value.id}`;
    await navigator.clipboard.writeText(shareUrl);

    shareNotification.value = "Share link copied to clipboard!";
    setTimeout(() => {
      shareNotification.value = null;
    }, 3000);
  } catch (err) {
    // Fallback for browsers that don't support clipboard API
    const shareUrl = `${window.location.origin}/lists/share/${list.value.id}`;
    const textArea = document.createElement("textarea");
    textArea.value = shareUrl;
    textArea.style.position = "fixed";
    textArea.style.left = "-999999px";
    document.body.appendChild(textArea);
    textArea.select();
    try {
      document.execCommand("copy");
      shareNotification.value = "Share link copied to clipboard!";
      setTimeout(() => {
        shareNotification.value = null;
      }, 3000);
    } catch (e) {
      shareNotification.value =
        "Failed to copy link. Please copy manually: " + shareUrl;
      setTimeout(() => {
        shareNotification.value = null;
      }, 5000);
    }
    document.body.removeChild(textArea);
  }
};

const handleDeleteList = async () => {
  if (!list.value) return;

  // Confirm deletion
  if (
    !confirm(
      `Are you sure you want to delete "${list.value.name}"? This action cannot be undone.`
    )
  ) {
    return;
  }

  isDeletingList.value = true;
  error.value = null;

  try {
    const listId = route.params.id as string;
    await deleteList(listId);
    // Redirect to dashboard after successful deletion
    await router.push("/dashboard");
  } catch (err: any) {
    error.value = err.data?.error || err.message || "Failed to delete list";
    isDeletingList.value = false;
  }
};

// Load list on page load
await loadList();

// Cleanup timers on unmount
onUnmounted(() => {
  debounceTimers.forEach((timer) => clearTimeout(timer));
  debounceTimers.clear();
});
</script>
