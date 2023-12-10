.PHONY: install
install:
	cd precompiles && yarn install
	cd solgen && npm install

.PHONY: solgen
solgen:
	cd precompiles && yarn build
	go run gen.go 1
	mv FheOps_gen.sol solidity/FheOS.sol
	cd solgen && npm run build

.PHONY: clean
clean:
	rm solidity/*