all: zenad-env

zenad-env: zenad-rmi
	docker build --tag zenanetwork/zenad -f zenad-env/Dockerfile \
		$(shell git rev-parse --show-toplevel)

zenad-rmi:
	docker rmi zenanetwork/zenad 2>/dev/null; true

.PHONY: all zenad-env zenad-rmi
