// src/stores/appointmentStore.js
import { reactive, readonly } from 'vue'

// Состояние
const state = reactive({
  // Данные с сервера
  specializations: [],
  doctors: [],
  
  // Данные записи
  appointmentData: {
    // Шаг 1
    specialization: null,
    doctor: null,
    
    // Шаг 2
    lastName: '',
    firstName: '',
    middleName: '',
    birthDate: '',
    
    // Шаг 3
    appointmentDate: '',
    appointmentTime: '',
    
    // Шаг 4
    phone: '',
    preferredContact: {
      phone: false,
      telegram: false,
      whatsapp: false
    }
  },

  // UI состояния
  loading: false,
  error: null,
  currentStep: 1
})

// Методы
const methods = {
  // Загрузка специализаций
  async loadSpecializations() {
    state.loading = true
    state.error = null
    
    try {
      const response = await fetch('/api/specializations')
      if (!response.ok) throw new Error('Ошибка загрузки')
      state.specializations = await response.json()
    } catch (err) {
      state.error = err.message
      console.error(err)
    } finally {
      state.loading = false
    }
  },

  // Загрузка врачей по специализации
  async loadDoctorsBySpecialization(specializationId) {
    if (!specializationId) return
    
    state.loading = true
    state.error = null
    
    try {
      const response = await fetch(`/api/doctors/by-specialization/${specializationId}`)
      if (!response.ok) throw new Error('Ошибка загрузки')
      state.doctors = await response.json()
    } catch (err) {
      state.error = err.message
      console.error(err)
    } finally {
      state.loading = false
    }
  },

  // Сеттеры
  setSpecialization(spec) {
    state.appointmentData.specialization = spec
    state.appointmentData.doctor = null
    if (spec) {
      this.loadDoctorsBySpecialization(spec.id)
    } else {
      state.doctors = []
    }
  },

  setDoctor(doctor) {
    state.appointmentData.doctor = doctor
  },

  setPatientData(data) {
    state.appointmentData.lastName = data.lastName || ''
    state.appointmentData.firstName = data.firstName || ''
    state.appointmentData.middleName = data.middleName || ''
    state.appointmentData.birthDate = data.birthDate || ''
  },

  setDateTime(data) {
    state.appointmentData.appointmentDate = data.date || ''
    state.appointmentData.appointmentTime = data.time || ''
  },

  setContactData(data) {
    state.appointmentData.phone = data.phone || ''
    state.appointmentData.preferredContact = {
      phone: data.phoneContact || false,
      telegram: data.telegram || false,
      whatsapp: data.whatsapp || false
    }
  },

  // Навигация
  nextStep() {
    if (state.currentStep < 5) state.currentStep++
  },

  prevStep() {
    if (state.currentStep > 1) state.currentStep--
  },

  // Сброс
  reset() {
    state.appointmentData = {
      specialization: null,
      doctor: null,
      lastName: '',
      firstName: '',
      middleName: '',
      birthDate: '',
      appointmentDate: '',
      appointmentTime: '',
      phone: '',
      preferredContact: {
        phone: false,
        telegram: false,
        whatsapp: false
      }
    }
    state.doctors = []
    state.currentStep = 1
    state.error = null
  }
}

// Геттеры (computed свойства)
const getters = {
  isStep1Complete: () => {
    return state.appointmentData.specialization && state.appointmentData.doctor
  },
  
  isStep2Complete: () => {
    return state.appointmentData.lastName && 
           state.appointmentData.firstName && 
           state.appointmentData.birthDate
  },
  
  isStep3Complete: () => {
    return state.appointmentData.appointmentDate && 
           state.appointmentData.appointmentTime
  },
  
  isStep4Complete: () => {
    const contact = state.appointmentData.preferredContact
    return state.appointmentData.phone && 
           (contact.phone || contact.telegram || contact.whatsapp)
  }
}

// Экспортируем store
export const appointmentStore = {
  state: readonly(state),
  ...methods,
  ...getters
}