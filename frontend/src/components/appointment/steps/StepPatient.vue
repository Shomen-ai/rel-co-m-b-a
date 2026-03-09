<template>
  <div>

    <h2 class="text-xl font-bold mb-6">
      2. Введите ваши данные
    </h2>

    <div class="grid grid-cols-2 gap-4 mb-6">
      <input 
        v-model="form.lastName"
        class="border rounded-lg px-4 py-2" 
        placeholder="Фамилия *" 
        type="text"
      />
      <input 
        v-model="form.firstName"
        class="border rounded-lg px-4 py-2" 
        placeholder="Имя *" 
        type="text"
      />
      <input 
        v-model="form.middleName"
        class="border rounded-lg px-4 py-2" 
        placeholder="Отчество" 
        type="text"
      />
      <input 
        v-model="form.birthDate"
        class="border rounded-lg px-4 py-2" 
        placeholder="Дата рождения *" 
        type="date"
      />
    </div>

    <div class="flex gap-4">
      <button 
        @click="emit('back')" 
        class="text-gray-500 hover:text-gray-700"
      >
        ← Вернуться обратно
      </button>

      <button
        @click="handleNext"
        class="bg-primary text-white px-6 py-2 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed"
        :disabled="!isFormValid"
      >
        Продолжить →
      </button>
    </div>

    <!-- Подсказка об обязательных полях -->
    <p class="text-xs text-gray-400 mt-4">
      * — обязательные поля
    </p>

  </div>
</template>

<script setup>
import { ref, watch } from "vue"
import { appointmentStore } from "@/modules/AppointmentModule.js"

const emit = defineEmits(["next", "back"])

// Используем store
const store = appointmentStore

// Локальная форма для ввода
const form = ref({
  lastName: store.state.appointmentData.lastName || '',
  firstName: store.state.appointmentData.firstName || '',
  middleName: store.state.appointmentData.middleName || '',
  birthDate: store.state.appointmentData.birthDate || ''
})

// Проверка заполненности обязательных полей
const isFormValid = ref(false)

// Следим за изменениями формы
watch(form, () => {
  isFormValid.value = form.value.lastName && 
                      form.value.firstName && 
                      form.value.birthDate
}, { deep: true, immediate: true })

// Сохраняем данные и переходим дальше
function handleNext() {
  store.setPatientData(form.value)
  emit('next')
}
</script>