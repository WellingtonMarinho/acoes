FLUTTER := /Users/wellingtonsoares/development/flutter/bin/flutter
ANDROID_DEVICE ?= emulator-5554

.PHONY: run-backend test-backend test-backend-integration run-mobile test-mobile

run-backend:
	docker compose up --build

test-backend:
	cd backend && go test ./...

test-backend-integration:
	cd backend && go test -tags=integration ./internal/postgres

run-mobile:
	cd mobile && $(FLUTTER) run -d $(ANDROID_DEVICE)

test-mobile:
	cd mobile && $(FLUTTER) test
