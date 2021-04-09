ENV_LOCAL_FILE := env.local
ENV_LOCAL = $(shell cat $(ENV_LOCAL_FILE))

serve:
	$(ENV_LOCAL) go run main.go

build:
	$(ENV_LOCAL) go build main.go

.PHONY: serve