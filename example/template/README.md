# Template contract

Do not run this code, it is meant as a starter for you.

1. copy this to another directory
2. change the name in the source code
3. run

Please replace you desired name everywhere you see foobar

```sh
cp -r template foobar
cd foobar
sed -i"" -e 's/TEMPLATE/foobar/' makefile
sed -i"" -e 's/TEMPLATE/foobar/' integration/integration_test.go 
```

Once that is working, you should be able to compile the contract and run tests:

```sh 
# fast pure go tests
make unit-test

# compile the contract with tinygo (needs docker)
make build

# compile the contract and run all tests (including integration tests on the compiled wasm)
make test
```