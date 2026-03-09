<script setup>
import { ref, computed, watch } from "vue"
import { appointmentStore } from "@/modules/AppointmentModule.js"

const emit = defineEmits(["next", "back"])

// Используем store
const store = appointmentStore

// Состояния
const submitting = ref(false)
const submitError = ref('')

// Локальная форма
const form = ref({
  phone: store.state.appointmentData.phone || '',
  contactPhone: store.state.appointmentData.preferredContact?.phone || false,
  contactTelegram: store.state.appointmentData.preferredContact?.telegram || false,
  contactWhatsApp: store.state.appointmentData.preferredContact?.whatsapp || false
})

// Ошибки валидации
const errors = ref({
  phone: '',
  contact: ''
})

// Валидация российского номера телефона
function validatePhone(phone) {
  if (!phone) return 'Номер телефона обязателен'
  
  const cleaned = phone.replace(/[^\d+]/g, '')
  const phoneRegex = /^(\+7|8)\d{10}$/
  
  if (!phoneRegex.test(cleaned)) {
    return 'Неверный формат номера (пример: +79123456789 или 89123456789)'
  }
  
  return ''
}

// Валидация выбора способа связи
function validateContact() {
  if (!form.value.contactPhone && !form.value.contactTelegram && !form.value.contactWhatsApp) {
    return 'Выберите хотя бы один способ связи'
  }
  return ''
}

// Проверка валидности формы
const isFormValid = computed(() => {
  const phoneError = validatePhone(form.value.phone)
  const contactError = validateContact()
  return !phoneError && !contactError
})

// Валидация при изменении
watch(form, () => {
  errors.value.phone = validatePhone(form.value.phone)
  errors.value.contact = validateContact()
}, { deep: true })

// Форматирование номера при вводе
function formatPhoneInput(e) {
  let value = e.target.value.replace(/[^\d+]/g, '')
  
  if (!value.startsWith('+') && !value.startsWith('8')) {
    value = '+7' + value
  }
  
  if (value.startsWith('+7') && value.length > 12) {
    value = value.slice(0, 12)
  } else if (value.startsWith('8') && value.length > 11) {
    value = value.slice(0, 11)
  }
  
  form.value.phone = value
}

// Маска для отображения
const displayPhone = computed({
  get: () => {
    const phone = form.value.phone
    if (!phone) return ''
    
    if (phone.startsWith('+7')) {
      return phone.replace(/(\+7)(\d{3})(\d{3})(\d{2})(\d{2})/, '$1 ($2) $3-$4-$5')
    } else if (phone.startsWith('8')) {
      return phone.replace(/(8)(\d{3})(\d{3})(\d{2})(\d{2})/, '$1 ($2) $3-$4-$5')
    }
    return phone
  },
  set: (value) => {
    form.value.phone = value.replace(/[^\d+]/g, '')
  }
})

// Подготовка данных для отправки
function prepareAppointmentData() {
  const data = store.state.appointmentData
  
  return {
    // Данные клиента
    phone: form.value.phone,
    last_name: data.lastName,
    first_name: data.firstName,
    middle_name: data.middleName,
    birth_date: data.birthDate,
    contact_phone: form.value.contactPhone,
    contact_telegram: form.value.contactTelegram,
    contact_whatsapp: form.value.contactWhatsApp,
    
    // Данные записи
    doctor_id: data.doctor?.id,
    time_slot_id: null, // Нужно будет получить ID слота по дате и времени
    appointment_date: data.appointmentDate,
    appointment_time: data.appointmentTime
  }
}

// Поиск ID слота по дате и времени
async function findTimeSlotId(doctorId, date, time) {
  try {
    const response = await fetch(`/api/doctors/${doctorId}/slots/${date}`)
    if (!response.ok) throw new Error('Ошибка загрузки слотов')
    const slots = await response.json()
    
    // Ищем слот с нужным временем
    const slot = slots.find(s => s.start_time === time)
    return slot?.id || null
  } catch (err) {
    console.error('Ошибка при поиске слота:', err)
    return null
  }
}

// Отправка данных на сервер
async function submitAppointment() {
  submitting.value = true
  submitError.value = ''
  
  try {
    const data = prepareAppointmentData()
    
    // Получаем ID слота
    const timeSlotId = await findTimeSlotId(
      data.doctor_id,
      data.appointment_date,
      data.appointment_time
    )
    
    if (!timeSlotId) {
      throw new Error('Выбранный слот больше не доступен')
    }
    
    // Формируем финальные данные для отправки
    const appointmentData = {
      doctor_id: data.doctor_id,
      time_slot_id: timeSlotId,
      last_name: data.last_name,
      first_name: data.first_name,
      middle_name: data.middle_name,
      birth_date: data.birth_date,
      phone: data.phone,
      contact_phone: data.contact_phone,
      contact_telegram: data.contact_telegram,
      contact_whatsapp: data.contact_whatsapp
    }
    
    // Отправляем на сервер
    const response = await fetch('/api/appointments', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(appointmentData)
    })
    
    if (!response.ok) {
      const error = await response.text()
      throw new Error(error || 'Ошибка при создании записи')
    }
    
    const result = await response.json()
    console.log('Запись создана:', result)
    
    // Сохраняем данные в store перед переходом
    store.setContactData({
      phone: form.value.phone,
      phoneContact: form.value.contactPhone,
      telegram: form.value.contactTelegram,
      whatsapp: form.value.contactWhatsApp
    })
    
    // Переходим на следующий шаг (StepSuccess)
    emit('next')
    
  } catch (err) {
    submitError.value = err.message
    console.error('Ошибка отправки:', err)
  } finally {
    submitting.value = false
  }
}

// Обработчик кнопки "Продолжить"
async function handleNext() {
  const phoneError = validatePhone(form.value.phone)
  const contactError = validateContact()
  
  errors.value.phone = phoneError
  errors.value.contact = contactError
  
  if (!phoneError && !contactError) {
    await submitAppointment()
  }
}
</script>

<template>
  <div>
    <h2 class="text-xl font-bold mb-6">
      4. Укажите номер телефона
    </h2>

    <!-- Поле ввода телефона -->
    <div class="mb-2">
      <input
        v-model="displayPhone"
        @input="formatPhoneInput"
        class="w-full max-w-md border rounded-lg px-4 py-2"
        placeholder="+7 (XXX) XXX-XX-XX"
        type="tel"
        :disabled="submitting"
      />
      <p v-if="errors.phone" class="text-red-500 text-sm mt-1">
        {{ errors.phone }}
      </p>
      <p v-else class="text-gray-400 text-xs mt-1">
        Формат: +7XXXXXXXXXX или 8XXXXXXXXXX (10 цифр после кода)
      </p>
    </div>

    <!-- Способы связи -->
    <div class="mb-6">
      <p class="text-sm font-medium text-gray-700 mb-2">
        Предпочитаемый способ связи *
      </p>
      <div class="flex gap-6">
        <label class="flex items-center gap-2">
          <input 
            type="checkbox" 
            v-model="form.contactPhone"
            :disabled="submitting"
          />
          <span>Телефон</span>
        </label>

        <label class="flex items-center gap-2">
          <input 
            type="checkbox" 
            v-model="form.contactTelegram"
            :disabled="submitting"
          />
          <span>Telegram</span>
        </label>

        <label class="flex items-center gap-2">
          <input 
            type="checkbox" 
            v-model="form.contactWhatsApp"
            :disabled="submitting"
          />
          <span>WhatsApp</span>
        </label>
      </div>
      <p v-if="errors.contact" class="text-red-500 text-sm mt-1">
        {{ errors.contact }}
      </p>
    </div>

    <!-- Ошибка отправки -->
    <div v-if="submitError" class="mb-4 p-3 bg-red-50 text-red-600 rounded-lg">
      {{ submitError }}
    </div>

    <!-- Кнопки навигации -->
    <div class="flex gap-4">
      <button 
        @click="emit('back')" 
        class="text-gray-500 hover:text-gray-700"
        :disabled="submitting"
      >
        ← Вернуться обратно
      </button>

      <button
        @click="handleNext"
        class="bg-primary text-white px-6 py-2 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
        :disabled="!isFormValid || submitting"
      >
        <span v-if="submitting" class="inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
        {{ submitting ? 'Отправка...' : 'Продолжить →' }}
      </button>
    </div>
  </div>
</template>