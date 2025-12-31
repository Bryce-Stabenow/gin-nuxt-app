<template>
  <PageContainer>
    <div class="max-w-4xl mx-auto">
      <div class="bg-white rounded-xl shadow-2xl py-2 mb-6">
        <h1 class="text-3xl font-bold text-gray-900 mb-2 text-center">
          Dashboard
        </h1>
        
        <div v-if="isLoading" class="text-center text-gray-600 text-base py-5">
          Loading...
        </div>
        <div v-else-if="isAuthenticated && user">
          <p class="text-gray-900 text-lg leading-relaxed text-center">
            Welcome,
            <strong class="text-purple-600 font-semibold">{{
              user.profile.first_name
            }}</strong>
          </p>
        </div>
        <div v-else class="text-center text-red-800 py-5">
          <p class="mb-5 text-base">
            You are not authenticated. Please sign in.
          </p>
          <NuxtLink
            to="/signin"
            class="inline-block px-6 py-3 bg-gradient-to-r from-purple-500 to-purple-700 text-white rounded-lg font-semibold no-underline transition-transform hover:-translate-y-0.5 hover:shadow-lg"
          >
            Sign In
          </NuxtLink>
        </div>
      </div>

      <!-- Lists Section -->
      <div
        v-if="isAuthenticated && !isLoading"
        class="bg-white rounded-xl shadow-2xl py-10 px-4"
      >
        <div class="flex justify-between items-center mb-6">
          <h2 class="text-2xl font-bold text-gray-900">My Lists</h2>
          <NuxtLink
            to="/lists/new"
            class="inline-block px-6 py-3 bg-gradient-to-r from-purple-500 to-purple-700 text-white rounded-lg font-semibold no-underline transition-transform hover:-translate-y-0.5 hover:shadow-lg"
          >
            +
          </NuxtLink>
        </div>

        <div
          v-if="listsLoading"
          class="text-center text-gray-600 text-base py-5"
        >
          Loading lists...
        </div>
        
        <div v-else-if="listsError" class="text-center text-red-800 py-5">
          <p class="mb-5 text-base">Error: {{ listsError }}</p>
        </div>
        
        <div
          v-else-if="lists.length === 0"
          class="text-center text-gray-600 py-10"
        >
          <p class="mb-5 text-base">You don't have any lists yet.</p>
          <NuxtLink
            to="/lists/new"
            class="inline-block px-6 py-3 bg-gradient-to-r from-purple-500 to-purple-700 text-white rounded-lg font-semibold no-underline transition-transform hover:-translate-y-0.5 hover:shadow-lg"
          >
            Create Your First List
          </NuxtLink>
        </div>
        
        <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <DashboardList
            v-for="list in lists"
            :key="list.id"
            :list="list"
            :is-shared="isSharedList(list)"
          />
        </div>
      </div>
    </div>
  </PageContainer>
</template>

<script setup lang="ts">
const { isAuthenticated, user, isLoading, checkAuth } = useAuth();
const { getLists } = useLists();

// Set page title and meta tags
useHead({
  title: 'GrocerMe | Dashboard',
  meta: [
    {
      name: 'description',
      content: 'Manage your grocery lists, view your items, and organize your shopping with GrocerMe.'
    },
    {
      property: 'og:title',
      content: 'GrocerMe | Dashboard'
    },
    {
      property: 'og:description',
      content: 'Manage your grocery lists, view your items, and organize your shopping with GrocerMe.'
    },
    {
      property: 'og:type',
      content: 'website'
    },
  ]
});

const lists = ref<any[]>([]);
const listsLoading = ref(false);
const listsError = ref<string | null>(null);

// Check if a list is shared (not owned by current user)
const isSharedList = (list: any): boolean => {
  if (!user.value || !list) return false;
  return list.user_id !== user.value.id;
};

const loadLists = async () => {
  listsLoading.value = true;
  listsError.value = null;

  try {
    lists.value = await getLists();
  } catch (error: any) {
    listsError.value =
      error.data?.error || error.message || "Failed to load lists";
  } finally {
    listsLoading.value = false;
  }
};

// Check authentication on page load
await checkAuth();

if (isAuthenticated.value) {
  await loadLists();
}
</script>
