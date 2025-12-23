package resources

import "github.com/adm87/flinch/engine/flinch"

// LoadingTask is a function type that defines a single loading task within a loading operation.
type LoadingTask func(ctx *flinch.Context, rs *ResourceSystem, batchID uint64) error

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
func (lo *LoadingOperation) Execute(ctx *flinch.Context) error {
	for _, task := range lo.tasks {
		if err := task(ctx, lo.rs, lo.batchID); err != nil {
			return err
		}
	}
	return nil
}
