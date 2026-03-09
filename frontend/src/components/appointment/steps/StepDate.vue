<script setup>
import { ref, computed, onMounted, watch } from "vue"
import { appointmentStore } from "@/modules/AppointmentModule.js"

const emit = defineEmits(["next", "back"])

// Состояния
const availableDates = ref([])
const availableSlots = ref([])
const selectedDate = ref('')
const selectedTime = ref('')
const loading = ref(false)
const openDateDropdown = ref(false)
const openTimeDropdown = ref(false)
const error = ref('')

// Получаем выбранного врача из store
const selectedDoctor = computed(() => appointmentStore.state.appointmentData.doctor)

// Загружаем доступные даты для врача
async function loadAvailableDates() {
  if (!selectedDoctor.value) return
  
  loading.value = true
  error.value = ''
  
  try {
    const response = await fetch(`/api/doctors/${selectedDoctor.value.id}/available-dates`)
    if (!response.ok) throw new Error('Ошибка загрузки')
    availableDates.value = await response.json()
  } catch (err) {
    error.value = err.message
    console.error(err)
  } finally {
    loading.value = false
  }
}

// Загружаем доступные слоты для выбранной даты
async function loadAvailableSlots(date) {
  if (!selectedDoctor.value || !date) return
  
  loading.value = true
  error.value = ''
  
  try {
    const response = await fetch(`/api/doctors/${selectedDoctor.value.id}/slots/${date}`)
    if (!response.ok) throw new Error('Ошибка загрузки')
    availableSlots.value = await response.json()
  } catch (err) {
    error.value = err.message
    console.error(err)
  } finally {
    loading.value = false
  }
}

// Следим за выбранной датой
watch(selectedDate, async (newDate) => {
  selectedTime.value = ''
  if (newDate) {
    await loadAvailableSlots(newDate)
  } else {
    availableSlots.value = []
  }
})

// При монтировании загружаем доступные даты
onMounted(() => {
  if (selectedDoctor.value) {
    loadAvailableDates()
  }
})

// Форматирование даты для отображения
function formatDate(dateStr) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'long',
    weekday: 'short'
  })
}

// Форматирование времени
function formatTime(time) {
  return time.slice(0, 5)
}

// Выбор даты
function selectDate(dateStr) {
  selectedDate.value = dateStr
  openDateDropdown.value = false
}

// Выбор времени
function selectTime(slot) {
  selectedTime.value = slot.start_time
  // Сохраняем в store
  appointmentStore.setDateTime({
    date: selectedDate.value,
    time: selectedTime.value
  })
  openTimeDropdown.value = false
}

// Проверка, можно ли продолжить
const canContinue = computed(() => {
  return selectedDate.value && selectedTime.value
})

// Переход к следующему шагу
function handleNext() {
  if (canContinue.value) {
    emit('next')
  }
}
</script>

<template>
  <div>
    <h2 class="text-xl font-bold mb-6">
      3. Укажите дату и время записи
    </h2>

    <!-- Информация о выбранном враче -->
    <div v-if="selectedDoctor" class="mb-4 p-3 bg-gray-50 rounded-lg text-sm">
      <span class="text-gray-600">Врач:</span>
      <span class="font-semibold ml-1">
        {{ selectedDoctor.last_name }} {{ selectedDoctor.first_name }}
      </span>
    </div>

    <div class="flex gap-4 mb-6">
      <!-- Первый dropdown (дата) -->
      <div class="relative w-[220px]">
        <button
          @click="openDateDropdown = !openDateDropdown"
          class="w-full flex items-center justify-between border border-gray-300 rounded-lg px-4 py-2 bg-white text-sm"
          :disabled="loading"
        >
          <span class="text-gray-700">
            {{ selectedDate ? formatDate(selectedDate) : 'Выберите дату' }}
          </span>
          <span class="text-gray-400 text-xs">▼</span>
        </button>

        <!-- Выпадающий список дат -->
        <div
          v-if="openDateDropdown"
          class="absolute left-0 mt-2 w-full bg-white border rounded-lg shadow-md z-10 max-h-60 overflow-y-auto"
        >
          <!-- Загрузка -->
          <div v-if="loading" class="px-4 py-2 text-gray-400">
            Загрузка...
          </div>
          
          <!-- Ошибка -->
          <div v-else-if="error" class="px-4 py-2 text-red-500">
            {{ error }}
          </div>
          
          <!-- Список дат -->
          <div
            v-for="date in availableDates"
            :key="date"
            @click="selectDate(date)"
            class="px-4 py-2 hover:bg-gray-100 cursor-pointer text-sm"
            :class="{ 'bg-blue-50': selectedDate === date }"
          >
            {{ formatDate(date) }}
          </div>
          
          <!-- Нет доступных дат -->
          <div v-if="!loading && availableDates.length === 0" class="px-4 py-2 text-gray-400">
            Нет доступных дат
          </div>
        </div>
      </div>

      <!-- Второй dropdown (время) -->
      <div class="relative w-[220px]">
        <button
          @click="openTimeDropdown = !openTimeDropdown"
          class="w-full flex items-center justify-between border border-gray-300 rounded-lg px-4 py-2 bg-white text-sm"
          :disabled="!selectedDate || availableSlots.length === 0"
        >
          <span class="text-gray-700">
            {{ selectedTime ? formatTime(selectedTime) : 'Время приёма' }}
          </span>
          <span class="text-gray-400 text-xs">▼</span>
        </button>

        <!-- Выпадающий список времени -->
        <div
          v-if="openTimeDropdown"
          class="absolute left-0 mt-2 w-full bg-white border rounded-lg shadow-md z-10 max-h-60 overflow-y-auto"
        >
          <!-- Загрузка -->
          <div v-if="loading" class="px-4 py-2 text-gray-400">
            Загрузка...
          </div>
          
          <!-- Список слотов -->
          <div
            v-for="slot in availableSlots"
            :key="slot.id"
            @click="selectTime(slot)"
            class="px-4 py-2 hover:bg-gray-100 cursor-pointer text-sm"
            :class="{ 'bg-blue-50': selectedTime === slot.start_time }"
          >
            {{ formatTime(slot.start_time) }} - {{ formatTime(slot.end_time) }}
          </div>
          
          <!-- Нет доступного времени -->
          <div v-if="!loading && availableSlots.length === 0" class="px-4 py-2 text-gray-400">
            Нет доступного времени
          </div>
        </div>
      </div>

      <!-- Кнопка "Продолжить" -->
      <button
        @click="handleNext"
        class="bg-primary text-white px-5 py-2 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed"
        :disabled="!canContinue"
      >
        Продолжить
      </button>
    </div>

    <!-- Кнопка "Вернуться" -->
    <button @click="emit('back')" class="text-gray-500 hover:text-gray-700 text-sm">
      ← Вернуться обратно
    </button>
  </div>
</template>