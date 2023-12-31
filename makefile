app = todoapp
config = ./configs/config.yml

.PHONY: direnv
direnv:
	direnv allow .

.PHONY: build
build:
	go build -o ./bin/$(app) ./main.go

.PHONY: run
run:
	go run main.go -config-file=$(config)

# Testing
.PHONY: test
test:
	make -f makefile.test test-selected


# Tools (sudo reuired)
.PHONY: install/jq/arch
install/jq/arch:
	pacman -Sy jq

.PHONY: install/jq/debian
install/jq/debian:
	apt update;
	apt install -y jq

.PHONY: install/jq/centos
install/jq/centos:
	yum -y install https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm
	yum install jq -y

.PHONY: build/css
build/css:
	cd assets/web/src && npm run dev


