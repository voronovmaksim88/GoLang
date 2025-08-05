<template>
  <div id="app">
    <div class="container">
      <h1>Счетчик на Wails и Vue</h1>
      <div class="counter-box">
        <div class="counter-value">{{ counter }}</div>
        <div class="counter-controls">
          <button @click="decrement" class="btn btn-decrement">-</button>
          <button @click="increment" class="btn btn-increment">+</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'App',
  data() {
    return {
      counter: 0
    }
  },
  methods: {
    // Увеличение счетчика на 1
    async increment() {
      this.counter = await window.go.main.App.IncrementCounter(this.counter)
    },
    // Уменьшение счетчика на 1
    async decrement() {
      this.counter = await window.go.main.App.DecrementCounter(this.counter)
    }
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
}

.container {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  background-color: white;
}

h1 {
  margin-bottom: 2rem;
  color: #3498db;
}

.counter-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1.5rem;
}

.counter-value {
  font-size: 6rem;
  font-weight: bold;
  color: #2c3e50;
}

.counter-controls {
  display: flex;
  gap: 1rem;
}

.btn {
  border: none;
  padding: 0.75rem 1.5rem;
  font-size: 1.25rem;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-increment {
  background-color: #2ecc71;
  color: white;
}

.btn-increment:hover {
  background-color: #27ae60;
}

.btn-decrement {
  background-color: #e74c3c;
  color: white;
}

.btn-decrement:hover {
  background-color: #c0392b;
}
</style>