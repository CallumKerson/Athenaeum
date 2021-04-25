
build:
	@echo "Building droplets at './bin/library' ..."
	@go build -o bin/droplets

clean:
	rm -rf ./bin

all: lint	vet	cyclo	test	build

test:
	@echo "Running unit tests..."
	@go test -cover ./...

cyclo:
	@echo "Checking cyclomatic complexity..."
	@gocyclo -over 7 ./

vet:
	@echo "Running vet..."
	@go vet ./...

lint:
	@echo "Running golint..."
	~/go/bin/golint ./...
