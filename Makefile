.PHONY: install
install:
	cd precompiles && yarn install
	cd solgen && npm install

.PHONY: gen
gen: compile-go-tfhe
	./gen.sh
	rm solidity/tests/contracts/*.sol || true
	cd solgen && npm run build

.PHONY: compile
compile:
	cd solidity && pnpm compile

.PHONY: compile-go-tfhe
compile-go-tfhe:
	if [ ! -e ./go-tfhe/internal/api/amd64/libtfhe_wrapper.x86_64.so ]; then \
  			cd go-tfhe && make build; \
	fi

check_network_is_running:
	if [ -z $$(netstat -tln | grep ":8547") ]; then \
		echo "FHENIX NETWORK IS NOT LISTENING ON PORT 8547."; \
		exit 1; \
	fi

.PHONY: test
test: check_network_is_running gen compile
	cp solidity/.env.example solidity/.env
	cd solidity && pnpm install
	cd solidity && pnpm compile
	cd solidity && npm test

.PHONY: clean
clean:
	rm solidity/*
