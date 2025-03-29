FLYWAY_CLI = flyway

GO = go
GO_MAIN = main.go

all:
	@echo "Please specify a target (e.g., make migrate or make run)."

migrate:
	@$(FLYWAY_CLI) migrate

run:
	@$(GO) run $(GO_MAIN)

gen:
	@python3 ./scripts/config_autogen.py

.PHONY: all migrate run
