<script lang="ts">
import { defineComponent, ref } from 'vue'

export default defineComponent({
  name: 'HelloWorld',
  props: {
    msg: {
      type: String,
      required: true
    }
  },
  setup() {
    const count = ref(0)
    const increment = () => {
      count.value++
    }

    return { count, increment }
  }
})
</script>

<template>
  <h1>{{ msg }}</h1>
  <button @click="increment">Count is: {{ count }}</button>

  <div>
    <p class="text-indigo-400">Цвет text-indigo-400</p>
    <p class="text-indigo-500">Цвет text-indigo-500</p>
    <p class="text-indigo-600">Цвет text-indigo-600</p>
    <hr class="mb-5">
  </div>

</template>

<style scoped>
/* Ваши стили */
</style>