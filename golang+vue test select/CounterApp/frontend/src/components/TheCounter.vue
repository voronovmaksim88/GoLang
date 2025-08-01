<template>
  <div class="counter-container">
    <h1>Counter App</h1>
    <p>Current count: {{ counter }}</p>
    <div class="button-group">
      <button @click="increment">Increment</button>
      <button @click="decrement">Decrement</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const counter = ref<number>(0)

onMounted(async () => {
  try {
    counter.value = await window.go.main.App.GetCounter(counter.value)
  } catch (error) {
    console.error('Error fetching initial counter:', error)
  }
})

const increment = async () => {
  try {
    counter.value = await window.go.main.App.IncrementCounter(counter.value)
  } catch (error) {
    console.error('Error incrementing counter:', error)
  }
}

const decrement = async () => {
  try {
    counter.value = await window.go.main.App.DecrementCounter(counter.value)
  } catch (error) {
    console.error('Error decrementing counter:', error)
  }
}
</script>

<style scoped>
.counter-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
  font-family: Arial, sans-serif;
}

h1 {
  color: #333;
  margin-bottom: 20px;
}

p {
  font-size: 24px;
  margin: 10px 0;
}

.button-group {
  display: flex;
  gap: 10px;
}

button {
  padding: 10px 20px;
  font-size: 16px;
  cursor: pointer;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  transition: background-color 0.3s;
}

button:hover {
  background-color: #45a049;
}

button:active {
  background-color: #3d8b40;
}
</style>