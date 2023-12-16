.PHONY: install
install:
	cd precompiles && yarn install
	cd solgen && npm install

.PHONY: gen
gen:
	cd precompiles
	./gen.sh
	rm solidity/tests/contracts/*.sol || true
	cd solgen && npm run build

.PHONY: clean
clean:
	rm solidity/*
