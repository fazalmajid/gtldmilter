IMG=		postfix-milter-test
CONTAINER=	$(IMG)

all: gtldmilter # test

gtldmilter: go.mod go.sum gtldmilter.go
	go build gtldmilter.go

go.mod:
	go mod init gtldmilter
	go mod tidy

test: run
	python3 test.py

docker-img: gtldmilter.go main.cf master.cf gtlds.bad dests.bad run.sh aliases
	docker build -t $(IMG) .
	touch $@

stop:
	-docker kill $(CONTAINER)
	@while docker ps|grep $(CONTAINER); do echo waiting...; sleep 1; done
	-docker rm $(CONTAINER)

run: docker-img stop
	docker run --name $(CONTAINER) -d -p 127.0.0.1:2525:25/tcp $(IMG)

sh:
	docker exec -it $(CONTAINER) sh

logs:
	docker logs -f $(CONTAINER)

clean: stop
	-rm *~
