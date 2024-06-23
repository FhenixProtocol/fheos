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

.PHONY: gencompile
compile: gen compile

.PHONY: lint
lint:
	# cd solidity && pnpm solhint FHE.sol FheOS.sol tests/contracts/*.sol tests/contracts/utils/*.sol
	cd solidity && pnpm solhint FheOS.sol tests/contracts/*.sol tests/contracts/utils/*.sol

check_network_is_running:
	@echo "Checking connection to 127.0.0.1:8547..."
	@nc -z -v 127.0.0.1 8547 2>/dev/null || (echo "Connection failed to localfhenix" && false)
	@echo "connected"

.PHONY: test
test: check_network_is_running gen compile
	cp solidity/.env.example solidity/.env
	cd solidity && pnpm install
	cd solidity && pnpm test

.PHONY: build
build:
	go build -o build/main ./cmd/

.PHONY: clean
clean:
	rm -r build/*

.PHONY: clean-gen
clean-gen:
	find solidity/tests/contracts -type f \
	-not -name 'Counter.sol' \
	-not -name 'Ownership.sol' \
	-not -name 'Tx.sol' \
	-not -name 'wERC20.sol' \
	-not -name 'Utils.sol' \
	-delete
	rm solidity/FHE.sol
	rm solidity/FheOS.sol
