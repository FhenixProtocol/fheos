.PHONY: install
install:
	cd precompiles && pnpm install --frozen-lockfile=false
	cd solgen && pnpm install --frozen-lockfile=false
	cd solidity && pnpm install --frozen-lockfile=false

.PHONY: gen
gen:
	./gen.sh
	cd solgen && pnpm build

.PHONY: gen-fheops
gen-fheops:
	./gen.sh -g true -n true
	cd solgen && pnpm build

.PHONY: compile
compile:
	cp solidity/.env.example solidity/.env
	cd solidity && pnpm compile

.PHONY: gencompile
gencompile: gen compile

.PHONY: lint
lint:
	# cd solidity && pnpm solhint FHE.sol FheOS.sol tests/contracts/*.sol tests/contracts/utils/*.sol
	cd solidity && pnpm solhint --ignore-path .solhintignore FheOS.sol tests/contracts/*.sol tests/contracts/utils/*.sol

check_network_is_running:
	@echo "Checking connection to 127.0.0.1:8547..."
	@nc -z -v 127.0.0.1 8547 2>/dev/null || (echo "Connection failed to localfhenix" && false)
	@echo "connected"

.PHONY: test
test: check_network_is_running gen compile
	cd solidity && pnpm install
	cd solidity && pnpm test

.PHONY: test-precomp
test-precomp: check_network_is_running
	cp solidity/.env.example solidity/.env
	cd solidity && pnpm test -- precompiles.test.ts

.PHONY: test-tx
test-tx: check_network_is_running
	cp solidity/.env.example solidity/.env
	cd solidity && pnpm test -- transaction.test.ts

.PHONY: build
build:
	go build -o build/main ./cmd/

.PHONY: build-coprocessor
build-coprocessor:
	go build -o build/coprocessor ./http/
	chmod 0777 build/coprocessor

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
