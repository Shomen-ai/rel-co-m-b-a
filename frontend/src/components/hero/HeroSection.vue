<template>
  <section class="mt-16">
    <Container class="grid md:grid-cols-2 items-center gap-10">

      <!-- Doctor image -->
      <div class="flex justify-center">
        <img :src="clinicImage" alt="doctors" class="rounded-[24px] max-h-[500px]" />
      </div>

      <!-- Card -->
      <div class="bg-white p-10 rounded-[24px] shadow-card">
        <h1 class="text-3xl font-bold leading-snug mb-4">
          Медицинский центр для тех,
          кто ценит качество
        </h1>

        <p class="text-gray-600 mb-6">
          Надежное и качественное обследование и лечение
          в уютной атмосфере в центре Бреста
        </p>

        <div class="flex gap-4">
          <button class="bg-primary text-white px-6 py-3 rounded-xl"
          @click="emit('open-appointment')">
            Записаться онлайн
          </button>

          <button class="border border-primary text-primary px-6 py-3 rounded-xl">
            Подробнее о центре
          </button>
        </div>
      </div>

    </Container>
  </section>
</template>

<script setup>
import Container from '../layout/Container.vue'
import clinicImage from '@/assets/clinic_1.png'

const emit = defineEmits(["open-appointment"])
</script>