package backend

import (
	"context"
	"fmt"
	"installer/backend/discord"
	"log"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// var dialogs = NewDialogManager()
// var paths = NewPathManager()

// App struct
type Backend struct {
	ctx     context.Context
	dialogs *Dialogs
	actions *Actions
}

// NewApp creates a new App application struct
func CreateBackend() *Backend {
	created := &Backend{}
	created.dialogs = NewDialogManager()
	created.actions = NewActionsManager()
	return created
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (d *Backend) SetContext(ctx context.Context) {
	d.ctx = ctx
	d.dialogs.SetContext(ctx)
	d.actions.SetContext(ctx)
}

func (d *Backend) GetModules() []interface{} {
	return []interface{}{d.dialogs, d.actions}
}

func (d *Backend) GetDiscordPath(channel string) string {
	return discord.GetSuggestedPath(discord.ParseChannel(channel))
}

func (d *Backend) LogTest(input string) {
	fmt.Println("LOGTEST")
	runtime.EventsEmit(d.ctx, "log", "LOGTEST")
	log.Printf("testing")
	log.Print(input)
}
