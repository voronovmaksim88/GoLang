<template>
  <div class="select-container">
    <select v-model="selectedItem" :disabled="isDisabled">
      <option disabled value="">Выберите пункт...</option>
      <option v-for="(item, index) in options" :key="index" :value="item">
        {{ item }}
      </option>
    </select>

    <button
        @click="fetchOptions"
        :disabled="isLoading"
        class="fetch-button"
    >
      {{ buttonText }}
    </button>
  </div>
</template>

<script>
export default {
  name: 'AsyncDropdown',
  data() {
    return {
      selectedItem: '',
      options: [],
      isDisabled: false,
      isLoading: false
    }
  },
  computed: {
    buttonText() {
      return this.isLoading ? 'Загрузка...' : 'Загрузить опции';
    }
  },
  methods: {
    fetchOptions() {
      // Блокируем список и активируем состояние загрузки
      this.isDisabled = true;
      this.isLoading = true;

      // Имитация запроса к бэкенду
      setTimeout(() => {
        // Получаем "ответ" от сервера с данными
        this.options = [
          'Пункт 1',
          'Пункт 2',
          'Пункт 3',
          'Пункт 4',
          'Пункт 5'
        ];

        // Сбрасываем состояние загрузки и разблокируем список
        this.isLoading = false;
        this.isDisabled = false;

        // Автоматически выбираем первый пункт
        if (this.options.length > 0) {
          this.selectedItem = this.options[0];
        }
      }, 1500); // Имитация задержки сети 1.5 секунды
    }
  }
}
</script>

<style scoped>
.select-container {
  display: flex;
  gap: 10px;
  max-width: 300px;
  margin: 20px;
}

select {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #ccc;
  border-radius: 4px;
  background-color: white;
  font-size: 16px;
}

select:disabled {
  background-color: #f5f5f5;
  cursor: not-allowed;
}

.fetch-button {
  padding: 8px 16px;
  background-color: #42b983;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
  transition: background-color 0.3s;
}

.fetch-button:hover:not(:disabled) {
  background-color: #359c6f;
}

.fetch-button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}
</style>