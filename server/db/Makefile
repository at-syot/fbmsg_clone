#!make

include ../.env

MG_PATH=migration
URL="postgres://$(DB_USER):$(DB_PW)@$(DB_HOST):$(DB_PORT)/msg_clone?sslmode=disable"

migrate_up:
	migrate -path $(MG_PATH) -database $(URL) -verbose up

migrate_down:
	migrate -path $(MG_PATH) -database $(URL) -verbose down

migrate_down_latest:
	migrate -path $(MG_PATH) -database $(URL) -verbose down 1

migrate_fix:
	migrate -path $(MG_PATH) -database $(URL) force $(version)

migrate_create:
	migrate create -ext sql -dir $(MG_PATH) -seq $(name)

