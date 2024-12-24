check_network_is_running:
	@echo "Checking connection to 127.0.0.1:8547..."
	@nc -z -v 127.0.0.1 8547 2>/dev/null || (echo "Connection failed to localfhenix" && false)
	@echo "connected"

.PHONY: build
build:
	go build -o build/main ./cmd/

.PHONY: start-engine-async
start-engine-async:
	cd warp-drive/fhe-engine && make server-no-sgx & echo $$! > engine.pid
	for i in {1..20}; do \
		if nc -z localhost 50051; then echo "Engine is up!"; break; fi; \
		echo "Waiting for engine..."; sleep 1; \
	done

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
