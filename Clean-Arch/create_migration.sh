#!/bin/bash

if [ -z "$1" ]; then
  echo "Por favor, forneça uma descrição para a migração."
  exit 1
fi

TIMESTAMP=$(date +'%Y%m%d%H%M%S')
DESCRIPTION=$1
FILENAME="${TIMESTAMP}_${DESCRIPTION}.sql"

echo "-- +goose Up
CREATE TABLE IF NOT EXISTS orders
(
    id varchar(36) NOT NULL,
    price float NOT NULL,
    tax float NOT NULL,
    final_price float NOT NULL,
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE IF EXISTS orders;" > "./migrations/${FILENAME}"

echo "Migração criada: ./migrations/${FILENAME}"
