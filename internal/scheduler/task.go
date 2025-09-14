package scheduler

// Task represents a schedulable unit of work.
type Task struct {
	id       int
	state    TaskState
	blocked  bool
	function func()
}

// TaskState represents the state of a task.
type TaskState int

const (
	TaskReady TaskState = iota
	TaskRunning
	TaskBlocked
	TaskDone
)

// NewTask creates a new task.
func NewTask(id int, fn func()) *Task {
	return &Task{
		id:       id,
		state:    TaskReady,
		function: fn,
	}
}

// Run executes the task.
func (t *Task) Run() {
	t.state = TaskRunning
	t.function()
	t.state = TaskDone
}