build-mysql:
	docker build -t natanchagas/gin-crud-mysql -f .\deploy\repositories\mysql\Dockerfile .
run-mysql:
	docker run -p 3306:3306 -d --name mysql --rm natanchagas/gin-crud-mysql
exec-mysql:
	docker exect -it mysql /
stop-mysql:
	docker stop mysql