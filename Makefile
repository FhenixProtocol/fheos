check_network_is_running:
	@echo "Checking connection to 127.0.0.1:8547..."
	@nc -z -v 127.0.0.1 8547 2>/dev/null || (echo "Connection failed to localfhenix" && false)
	@echo "connected"

.PHONY: build
build:
	go build -o build/main ./cmd/

.PHONY: start-engine
start-engine:
	cd warp-drive/fhe-engine && make server-no-sgx & echo $$! > engine.pid

.PHONY: stop-engine
stop-engine:
	if [ -f engine.pid ]; then kill `cat engine.pid` && rm -f engine.pid; fi

.PHONY: unit-test
unit-test:
	go test -failfast ./precompiles/

.PHONY: build-coprocessor
build-coprocessor:
	go build -o build/coprocessor ./http/
	chmod 0777 build/coprocessor

.PHONY: clean
clean:
	rm -r build/*
