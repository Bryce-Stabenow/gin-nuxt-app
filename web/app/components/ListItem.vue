<template>
  <div
    class="flex items-stretch gap-0 p-0 overflow-hidden rounded-lg"
  >
    <div 
      class="flex-1 p-4 border-2 border-r-0 rounded-l-lg transition-colors cursor-pointer"
      :class="item.checked 
        ? 'border-gray-300 bg-gray-50' 
        : 'border-gray-200 hover:border-purple-300'"
      @click="handleItemClick"
    >
      <div class="flex items-center gap-2">
        <span
          class="font-medium text-gray-900"
          :class="{ 'line-through text-gray-500': item.checked }"
        >
          {{ item.name }}
        </span>
      </div>
      <div class="text-sm text-gray-500 mt-1 space-y-1">
        <div v-if="item.quantity > 0">
          Quantity: {{ item.quantity }}
        </div>
        <div v-if="item.details" class="text-gray-600 italic">
          {{ item.details }}
        </div>
      </div>
    </div>
    <div
      @click.stop="handleToggle"
      class="w-1/4 flex items-center justify-center cursor-pointer transition-colors border-2 rounded-r-lg"
      :class="item.checked 
        ? 'bg-green-500 hover:bg-green-600 border-green-500 hover:border-green-600' 
        : 'bg-gray-200 hover:bg-gray-300 border-gray-200 hover:border-gray-300'"
    >
      <input
        type="checkbox"
        :checked="item.checked"
        @change="handleChange"
        class="sr-only"
        tabindex="-1"
      />
      <svg
        v-if="item.checked"
        xmlns="http://www.w3.org/2000/svg"
        class="h-8 w-8 text-green-100"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
        stroke-width="3"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M5 13l4 4L19 7"
        />
      </svg>
      <svg
        v-else
        xmlns="http://www.w3.org/2000/svg"
        class="h-8 w-8 text-gray-400"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
        stroke-width="2"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M5 13l4 4L19 7"
        />
      </svg>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  item: {
    name: string;
    checked: boolean;
    quantity: number;
    details?: string;
  };
  originalIndex: number;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  toggle: [index: number];
  change: [index: number, event: Event];
  click: [index: number];
}>();

const handleToggle = () => {
  emit('toggle', props.originalIndex);
};

const handleChange = (event: Event) => {
  emit('change', props.originalIndex, event);
};

const handleItemClick = () => {
  emit('click', props.originalIndex);
};
</script>

