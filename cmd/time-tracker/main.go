package main

import (
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/klusekP/time_pr_tracking/internal/alarm"
	"github.com/klusekP/time_pr_tracking/internal/app"
	"github.com/klusekP/time_pr_tracking/internal/clock"
	"github.com/klusekP/time_pr_tracking/internal/mode"
	"github.com/klusekP/time_pr_tracking/internal/repository"
	"github.com/klusekP/time_pr_tracking/internal/ui"
)

func main() {
	repo, err := repository.NewSQLite("")
	if err != nil {
		log.Fatalf("database error: %v", err)
	}
	defer repo.Close()

	clk := clock.System{}
	sysAlarm := alarm.NewSystem()

	theme := ui.NewDefaultTheme()
	bigFont := ui.NewBigFont()
	renderer := app.NewRenderer(theme, bigFont)

	modes := []mode.Mode{
		mode.NewStopwatch(clk, repo),
		mode.NewPomodoro(clk, sysAlarm, 25*time.Minute),
		mode.NewHistory(),
	}

	program := tea.NewProgram(
		app.NewModel(repo, renderer, modes),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := program.Run(); err != nil {
		log.Fatalf("TUI error: %v", err)
	}
}
