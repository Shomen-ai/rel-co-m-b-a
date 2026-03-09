<script setup>
const emit = defineEmits(["close"])
</script>

<template>
  <div class="text-center">

    <div class="flex justify-center mb-6">

      <div class="w-20 h-20 bg-green-500 rounded-full flex items-center justify-center text-white text-3xl">
        ✓
      </div>

    </div>

    <h2 class="text-xl font-bold mb-2">
      Ваша запись сформирована!
    </h2>

    <p class="text-gray-500 mb-6">
      Наш менеджер свяжется с вами
      удобным для вас способом
    </p>

    <button
      @click="emit('close')"
      class="bg-primary text-white px-6 py-3 rounded-lg"
    >
      Продолжить
    </button>

  </div>
</template>