RUN := go run
GOOSE_BINARY := goose

.PHONY: setup run-migrations create-migrations

setup:
	curl -o $(GOOSE_BINARY) -L https://github.com/pressly/goose/releases/download/v3.18.0/goose_linux_x86_64
	chmod +x $(GOOSE_BINARY)
	docker-compose up -d

run-app:
	cd cmd/ordersystem && $(RUN) main.go wire_gen.go

run-migrations:
	$(GOOSE_BINARY) -dir ./migrations mysql "root:root@tcp(localhost:3306)/orders" up
	rm $(GOOSE_BINARY)

create-migrations:
	mkdir migrations
	chmod +x create_migration.sh
	./create_migration.sh create_table_orders
