include .env

.PHONY: migrations-up
migrations-up:
	@echo 'Running up migrations...'
	@migrate -path=./migrations -database="postgres://bannerservice:pass777@localhost/bannerservice?sslmode=disable" up

.PHONY: migrations-down
migrations-down:
	@echo 'Running down migrations...'
	@migrate -path=./migrations -database="postgres://bannerservice:pass777@localhost/bannerservice?sslmode=disable" down
