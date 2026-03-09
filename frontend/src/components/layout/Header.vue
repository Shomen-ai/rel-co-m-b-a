<template>
  <header class="bg-white shadow-card">
    <Container class="flex items-center justify-between py-4">

      <!-- Логотип -->
      <div class="flex items-center gap-3">
        <img :src="logo" alt="БьютиМедЛого" class="h-12 w-auto" />
      </div>

      <!-- Время работы (добавлено) -->
      <div class="hidden md:block text-sm text-gray-600">
        <div>пн-пт: 7:30 - 20:00</div>
        <div>сб-вс: 8:00 - 17:00</div>
      </div>

      <!-- Контакты (теперь в 2 строки) -->
      <div class="hidden md:block text-sm text-gray-600 text-right">
        <div class="font-medium">+7 (8422) 505 323</div>
        <div>г. Ульяновск, ул. Карбышева, 30</div>
      </div>

      <button 
        class="bg-primary text-white px-6 py-2 rounded-xl hover:opacity-90 transition whitespace-nowrap"
        @click="emit('open-appointment')"
      >
        Записаться онлайн
      </button>

    </Container>
  </header>
</template>

<script setup>
import Container from '../layout/Container.vue'
import logo from '@/assets/logo.jpg' 

const emit = defineEmits(["open-appointment"])
</script>