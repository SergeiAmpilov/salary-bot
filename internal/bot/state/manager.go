// internal/bot/state/manager.go
package state

import "sync"

type Manager struct {
	mu     sync.RWMutex
	states map[int64]*UserState // chatID → состояние
}

func NewManager() *Manager {
	return &Manager{
		states: make(map[int64]*UserState),
	}
}

func (m *Manager) Get(chatID int64) *UserState {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if s, ok := m.states[chatID]; ok {
		return s
	}
	return &UserState{Step: StepNone}
}

func (m *Manager) Set(chatID int64, state *UserState) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.states[chatID] = state
}

func (m *Manager) Clear(chatID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.states, chatID)
}
