.PHONY: gen/tbls lint/tbls

# tblsでスキーマ図を生成する
gen/tbls:
	go tool tbls doc --rm-dist --config config/.tbls.yml

# tblsでDBテーブルのlintを行う
lint/tbls:
	go tool tbls lint --config config/.tbls.yml

.PHONY: db/up db/down db/migrate db/wait db/init

db/up:
	docker run -d --name mysql-db -p 3306:3306 \
		-e MYSQL_ROOT_PASSWORD=root \
		-e MYSQL_DATABASE=test_service \
		mysql:8

db/down:
	docker stop mysql-db || true
	docker rm mysql-db   || true

db/migrate:
	go tool sql-migrate up -config=config/dbconfig.yml -env="local"

db/wait:
	until (echo 'SELECT 1' | mysql -h 127.0.0.1 -P 3306 -uroot -proot --silent &> /dev/null); do echo 'waiting for mysqld to be connectable...' && sleep 2; done

db/init: db/down db/up db/wait db/migrate

.PHONY: gen/model

gen/model:
	go run tools/modelgen/main.go

.PHONY: app/run
app/run:
	go run cmd/app/main.go
