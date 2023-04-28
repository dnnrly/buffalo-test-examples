

.PHONY: cucumber-test
cucumber-test:
	docker compose -f docker-compose.yml up \
		--build \
		--remove-orphans \
		--abort-on-container-exit \
		--attach-dependencies \
		tests

.PHONY: test
test:
	docker compose -f docker-compose.yml up \
		--build \
		--remove-orphans \
		--abort-on-container-exit \
		--attach-dependencies \
		tests

.PHONY: down-all
down-all:
	docker compose -f docker-compose.yml down
	# docker compose -f docker-compose.gherkin.yml down

.PHONY: ps-all
ps-all:
	docker compose -f docker-compose.yml ps
	# docker compose -f docker-compose.gherkin.yml ps