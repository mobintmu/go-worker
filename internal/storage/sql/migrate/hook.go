package migrate

import "go-worker/internal/config"

func RunMigrations(runner *Runner, cfg *config.Config) {
	if cfg.IsTest() {
		return
	}
	runner.Run()
}
