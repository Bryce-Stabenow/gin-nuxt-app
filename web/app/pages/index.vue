<template>
  <PageContainer>
    <div class="flex items-start items-center justify-center">
      <div class="text-center">
        <h1 class="text-4xl font-bold text-white mb-8">Welcome to GrocerMe</h1>
        <p class="text-white text-lg mb-8">
          Your personal grocery management solution
        </p>
      </div>
    </div>
  </PageContainer>
</template>

<script setup lang="ts">
const { checkAuth } = useAuth();

// Check authentication on page load
onMounted(async () => {
  await checkAuth();
});
</script>
