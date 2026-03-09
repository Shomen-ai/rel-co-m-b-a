<template>
  <div class="bg-white p-6 rounded-[24px] shadow-card hover:-translate-y-1 transition">

    <h3 class="text-lg font-bold mb-3">
      {{ title }}
    </h3>

    <p class="text-gray-600 text-sm mb-6">
      {{ description }}
    </p>

    <div class="w-10 h-10 rounded-full bg-lightBg flex items-center justify-center">
      →
    </div>

  </div>
</template>

<script setup>
defineProps({
  title: String,
  description: String,
})
</script>