.PHONY: install
install:
	cd precompiles && pnpm install --frozen-lockfile=false
	cd solgen && pnpm install --frozen-lockfile=false
	cd solidity && pnpm install --frozen-lockfile=false

.PHONY: gen
gen:
	./gen.sh
	cd solgen && pnpm build

.PHONY: compile
compile:
	cd solidity && pnpm compile

.PHONY: lint
lint:
	# cd solidity && pnpm solhint FHE.sol FheOS.sol tests/contracts/*.sol tests/contracts/utils/*.sol
	cd solidity && pnpm solhint FheOS.sol tests/contracts/*.sol tests/contracts/utils/*.sol

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
test: check_network_is_running compile-go-tfhe gen compile
	cp solidity/.env.example solidity/.env
	cd solidity && pnpm install
	cd solidity && pnpm test

.PHONY: build
build:
	go build -o build/main ./cmd/

.PHONY: clean
clean:
	rm -r build/*
