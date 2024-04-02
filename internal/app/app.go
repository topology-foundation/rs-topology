package app

import "fmt"

type App struct {
	Name        string
	Description string
	Action      func() error
}

func NewApp() *App {
	return &App{}
}

func (app *App) Run(args []string) error {
	fmt.Println("Start running an app...")
	fmt.Println("Name:", app.Name)
	fmt.Println("Description:", app.Description)
	fmt.Println()

	if app.Action == nil {
		return nil
	}

	if err := app.Action(); err != nil {
		return err
	}

	return nil
}
