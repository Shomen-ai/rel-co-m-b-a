<template>
  <section class="mt-20">
    <Container>
      <div class="grid md:grid-cols-3 gap-8">

        <ServiceCard
          title="Стоматология"
          description="Инновационные технологии клиники обеспечивают передовые решения..."
        />

        <ServiceCard
          title="Отоларингология"
          description="Профессиональная диагностика и лечение патологий уха, горла, носа"
        />

        <ServiceCard
          title="Дерматокосметология"
          description="Эффективное и безопасное лечение и поддержание безупречного вида"
        />

        <ServiceCard
          title="Рентгенодиагностика"
          description="Компьютерная визуализация в трехмерном формате"
        />

        <ServiceCard
          title="Онкология, маммология"
          description="Грамотная и точная диагностика онкологических заболеваний"
        />

        <ServiceCard
          title="Массаж"
          description="Поддержание тонуса мышц и их укрепление"
        />

      </div>
    </Container>
  </section>
</template>

<script setup>
import Container from '../layout/Container.vue'
import ServiceCard from './ServiceCard.vue'
</script>