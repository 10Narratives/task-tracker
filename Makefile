who ?= tracker
config ?= tracker.yaml

run:
	@go run cmd/$(who)/main.go --config config/$(config)