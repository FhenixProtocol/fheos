.PHONY: install
install:
	cd precompiles && yarn install
	cd solgen && npm install

.PHONY: gen
gen:
	./gen.sh
	rm solidity/tests/contracts/*.sol || true
	cd solgen && npm run build

.PHONY: compile
compile:
	cd solidity && pnpm compile

.PHONY: test
test: gen compile
	cd solidity && pnpm compile
	cd solidity && npm test

.PHONY: clean
clean:
	rm solidity/*
