<template>
  <h2 class="text-lg font-semibold mb-6">
    1. Выберите специализацию и лечащего врача
  </h2>

  <div class="flex gap-4 mb-6">
    <!-- Первый dropdown (специализация) -->
    <div class="relative w-[220px]">
      <button
        @click="openSpecialization = !openSpecialization"
        class="w-full flex items-center justify-between border border-gray-300 rounded-lg px-4 py-2 bg-white text-sm"
      >
        <span class="text-gray-700">
          {{ store.state.appointmentData.specialization?.name || 'Специализация' }}
        </span>
        <span class="text-gray-400 text-xs">▼</span>
      </button>

      <div
        v-if="openSpecialization"
        class="absolute left-0 mt-2 w-full bg-white border rounded-lg shadow-md z-10 max-h-60 overflow-y-auto"
      >
        <!-- Загрузка -->
        <div v-if="store.state.loading" class="px-4 py-2 text-gray-400">
          Загрузка...
        </div>
        
        <!-- Ошибка -->
        <div v-else-if="store.state.error" class="px-4 py-2 text-red-500">
          {{ store.state.error }}
        </div>
        
        <!-- Список специализаций -->
        <div
          v-for="spec in store.state.specializations"
          :key="spec.id"
          @click="selectSpecialization(spec)"
          class="px-4 py-2 hover:bg-gray-100 cursor-pointer text-sm"
          :class="{ 'bg-blue-50': store.state.appointmentData.specialization?.id === spec.id }"
        >
          {{ spec.name }}
        </div>
      </div>
    </div>

    <!-- Второй dropdown (врач) -->
    <div class="relative w-[220px]">
      <button
        @click="openDoctor = !openDoctor"
        class="w-full flex items-center justify-between border border-gray-300 rounded-lg px-4 py-2 bg-white text-sm"
        :disabled="!store.state.appointmentData.specialization"
      >
        <div class="text-left">
          <div class="text-gray-700">
            {{ formatDoctorName(store.state.appointmentData.doctor) || 'Лечащий врач' }}
          </div>
          <div v-if="store.state.appointmentData.doctor" class="text-xs text-gray-500">
            {{ store.state.appointmentData.doctor.specialization?.name || '' }}
          </div>
        </div>
        <span class="text-gray-400 text-xs">▼</span>
      </button>

      <div
        v-if="openDoctor"
        class="absolute left-0 mt-2 w-full bg-white border rounded-lg shadow-md z-10 max-h-60 overflow-y-auto"
      >
        <!-- Загрузка -->
        <div v-if="store.state.loading" class="px-4 py-2 text-gray-400">
          Загрузка врачей...
        </div>
        
        <!-- Список врачей -->
        <div
          v-for="doctor in store.state.doctors"
          :key="doctor.id"
          @click="selectDoctor(doctor)"
          class="px-4 py-2 hover:bg-gray-100 cursor-pointer text-sm"
          :class="{ 'bg-blue-50': store.state.appointmentData.doctor?.id === doctor.id }"
        >
          <div>{{ doctor.last_name }} {{ doctor.first_name }} {{ doctor.middle_name || '' }}</div>
          <div class="text-xs text-gray-500">Стаж: {{ doctor.experience_years }} лет</div>
        </div>
        
        <!-- Нет врачей -->
        <div v-if="!store.state.loading && store.state.doctors.length === 0" 
             class="px-4 py-2 text-gray-400">
          Нет доступных врачей
        </div>
      </div>
    </div>

    <button
      @click="nextStep"
      class="bg-primary text-white px-5 py-2 rounded-lg"
      :disabled="!store.isStep1Complete()"
    >
      Продолжить
    </button>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from "vue"
import { appointmentStore } from "@/modules/AppointmentModule.js"

const emit = defineEmits(["next"])

// Используем store
const store = appointmentStore

// Состояния для открытия/закрытия dropdown'ов
const openSpecialization = ref(false)
const openDoctor = ref(false)

// Форматирование имени врача
function formatDoctorName(doctor) {
  if (!doctor) return ''
  return `${doctor.last_name} ${doctor.first_name}`
}

// При монтировании загружаем специализации
onMounted(async () => {
  if (store.state.specializations.length === 0) {
    await store.loadSpecializations()
  }
})

// Выбор специализации
function selectSpecialization(spec) {
  store.setSpecialization(spec)
  openSpecialization.value = false
  openDoctor.value = false
}

// Выбор врача
function selectDoctor(doctor) {
  store.setDoctor(doctor)
  openDoctor.value = false
}

// Переход к следующему шагу
function nextStep() {
  if (store.isStep1Complete()) {
    emit('next')
  }
}

// Закрытие dropdown'ов при клике вне
function handleClickOutside(event) {
  if (!event.target.closest('.relative')) {
    openSpecialization.value = false
    openDoctor.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>