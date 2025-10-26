// internal/bot/state/state.go
package state

type Step int

const (
	StepNone Step = iota
	StepAwaitingTech
	StepAwaitingExperience
)

// UserState хранит текущее состояние диалога с пользователем
type UserState struct {
	Step       Step
	Tech       string
	Experience string
}
