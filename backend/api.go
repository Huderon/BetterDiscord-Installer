package backend

import (
	"context"
	"installer/backend/utils"
	"os/exec"
	"runtime"
)

// var dialogs = NewDialogManager()
// var paths = NewPathManager()

// App struct
type Backend struct {
	ctx     context.Context
	dialogs *Dialogs
	paths   *Paths
	actions *Actions
}

// NewApp creates a new App application struct
func CreateBackend() *Backend {
	created := &Backend{}
	created.dialogs = NewDialogManager()
	created.paths = NewPathManager()
	created.actions = NewActionsManager()
	return created
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (d *Backend) SetContext(ctx context.Context) {
	d.ctx = ctx
	d.dialogs.SetContext(ctx)
	d.paths.SetContext(ctx)
	d.actions.SetContext(ctx)
}

func (d *Backend) GetModules() []interface{} {
	return []interface{}{d.dialogs, d.paths, d.actions}
}

func (d *Backend) GetPlatform() string {
	return runtime.GOOS
}

func (d *Backend) KillProcess(name string, shouldRestart bool) error {
	exeName := utils.GetProcessExe(name)
	// fmt.Println(exeName)
	err := utils.KillProcess(name)

	if err != nil {
		return err
	}

	if shouldRestart {
		cmd := exec.Command(exeName)
		cmd.Start()
	}

	return nil
}
