package main

import (
	"context"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp создает новый экземпляр App
func NewApp() *App {
	return &App{}
}

// startup вызывается при запуске приложения
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Определим методы для счетчика, которые будут вызываться из Vue

// GetCounter возвращает текущее значение счетчика
func (a *App) GetCounter(counter int) int {
	return counter
}

// IncrementCounter увеличивает счетчик на 1
func (a *App) IncrementCounter(counter int) int {
	return counter + 1
}

// DecrementCounter уменьшает счетчик на 1
func (a *App) DecrementCounter(counter int) int {
	return counter - 1
}
