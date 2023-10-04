MG_PATH=migration
URL="postgres://aot:1234@localhost:5432/msg_clone?sslmode=disable"

migrate_up:
	migrate -path $(MG_PATH) -database $(URL) -verbose up

migrate_down:
	migrate -path $(MG_PATH) -database $(URL) -verbose down

migrate_fix:
	migrate -path $(MG_PATH) -database $(URL) force $(version)

migrate_create:
	migrate create -ext sql -dir $(MG_PATH) -seq $(name)

