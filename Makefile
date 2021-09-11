DOCKER_COMPOSE_TEST_RUN ?= docker-compose -f docker-compose-test.yaml

.PHONY: test
test:
	${DOCKER_COMPOSE_TEST_RUN} run --rm test /bin/sh -c "make internal-test"
	${DOCKER_COMPOSE_TEST_RUN} down


.PHONY: lint
lint:
	-${DOCKER_COMPOSE_TEST_RUN} run --rm test /bin/sh -c "make internal-lint"
	${DOCKER_COMPOSE_TEST_RUN} down

.PHONY: internal-test
internal-test:
	CGO_ENABLED="1" \
	CGO_LDFLAGS="$(shell mapnik-config --libs)" \
	CGO_CXXFLAGS="$(shell mapnik-config --cxxflags --includes --dep-includes | tr '\n' ' ')" \
	go test ./...

internal-lint:
	CGO_ENABLED="1" \
	CGO_LDFLAGS="$(shell mapnik-config --libs)" \
	CGO_CXXFLAGS="$(shell mapnik-config --cxxflags --includes --dep-includes | tr '\n' ' ')" \
	golangci-lint run -v
