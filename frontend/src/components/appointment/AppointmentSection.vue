<template>
  <section class="mt-24">
    <div class="max-w-[1200px] mx-auto px-4">

      <div class="bg-white rounded-[24px] shadow-card p-12 grid md:grid-cols-2 gap-16 items-center">

        <!-- LEFT -->
        <div>

          <h2 class="text-2xl font-bold mb-4">
            Записаться на прием
          </h2>

          <p class="text-gray-500 mb-8 max-w-[420px]">
            Мы перезвоним вам и ответим на все вопросы, подберем
            удобное для визита время и лечащего врача
          </p>

          <!-- inputs -->
          <div class="space-y-4 mb-6">

            <input
              type="text"
              placeholder="Ваше имя"
              class="w-full bg-gray-100 rounded-lg px-4 py-3 outline-none focus:ring-2 focus:ring-primary"
            />

            <input
              type="tel"
              placeholder="+375"
              class="w-full bg-gray-100 rounded-lg px-4 py-3 outline-none focus:ring-2 focus:ring-primary"
            />

          </div>

          <!-- buttons -->
          <div class="flex gap-4 mb-4">

            <button
              class="bg-primary text-white px-6 py-3 rounded-lg font-medium hover:opacity-90"
            >
              Заказать звонок
            </button>

            <button
              @click="emit('open-appointment')"
              class="border border-primary text-primary px-6 py-3 rounded-lg font-medium"
            >
              Записаться онлайн
            </button>

          </div>

          <!-- checkbox -->
          <label class="flex items-start gap-2 text-sm text-gray-500">

            <input
              type="checkbox"
              checked
              class="mt-1 accent-primary"
            />

            <span>
              Отправляя сообщение, Вы соглашаетесь
              с публичным договором
            </span>

          </label>

        </div>

        <!-- RIGHT IMAGES -->
        <div class="flex justify-center">

  <img
    :src="appointmentImage"
    class="rounded-[24px] object-cover w-full max-w-[420px]"
  />

</div>

      </div>

    </div>
  </section>
</template>

<script setup>
import appointmentImage from "@/assets/clinic_3.png"
const emit = defineEmits(["open-appointment"])
</script>