<template>
  <div class="mb-5">
    <label :for="id" class="block text-gray-900 mb-2 font-medium text-sm">
      {{ label }}
      <slot name="label-suffix"></slot>
    </label>
    <input
      :id="id"
      :type="type"
      v-model="model"
      :required="required"
      :minlength="minlength"
      class="w-full px-3 py-3 border-2 border-gray-200 rounded-lg text-base transition-colors focus:outline-none focus:border-purple-500"
    />
  </div>
</template>

<script setup lang="ts">
defineProps<{
  id: string;
  label: string;
  type?: string;
  required?: boolean;
  minlength?: number;
}>();

const model = defineModel<string>();
</script>

