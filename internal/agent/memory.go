package agent

type Memory interface {
	Add(step Step)
	Recent(n int) []Step
	Clear()
}

type Step struct {
	Type      string
	Thought   string // LLM reasoning text
	ToolName  string
	ToolUseID string
	Input     string
	Output    string
}

type ShortTermMemory struct {
	steps []Step
}

func (m *ShortTermMemory) Add(step Step) {
	m.steps = append(m.steps, step)
}
func (m *ShortTermMemory) Recent(n int) []Step {
	l := len(m.steps)
	if n > l {
		return m.steps
	}
	return m.steps[l-n:]
}
func (m *ShortTermMemory) Clear() {
	m.steps = m.steps[:0]
}
