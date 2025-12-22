package resources

// LoadingTask is a function type that defines a single loading task within a loading operation.
type LoadingTask func(rs *ResourceSystem, batchID uint64) error

// LoadingOperation represents a batch of loading tasks to be executed within a resource system.
type LoadingOperation struct {
	batchID uint64
	rs      *ResourceSystem
	tasks   []LoadingTask
}

// AddTask adds a new loading task to the LoadingOperation.
func (lo *LoadingOperation) AddTask(task LoadingTask) {
	lo.tasks = append(lo.tasks, task)
}

// Execute performs all loading tasks within the LoadingOperation.
func (lo *LoadingOperation) Execute() error {
	for _, task := range lo.tasks {
		if err := task(lo.rs, lo.batchID); err != nil {
			return err
		}
	}
	return nil
}
