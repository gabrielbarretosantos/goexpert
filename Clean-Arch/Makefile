RUN := go run


setup:
	docker-compose up -d

run-app:
	cd cmd/ordersystem && $(RUN) main.go wire_gen.go

run-migrations:
	goose -dir ./migrations mysql "root:root@tcp(localhost:3306)/orders" up

create-migrations:
	goose -dir ./migrations create  create_table_orders sql
