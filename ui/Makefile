help:
	@perl -nle'print $& if m{^[a-zA-Z_-]+:.*?## .*$$}' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

clean: ## Wipe dist and elm-stuff
	rm -rf elm-stuff dist

build: setup ## Build Production Dist
	npm run build

check-compile: ## Build Elm to /dev/null to see if it compiles
	elm-make src/elm/Main.elm --output /dev/null

clean-elm: ## Clear the build artifacts
	rm -rf elm-stuff/build-artifacts/0.17.1/user/project/1.0.0/

run: ## Run the dev server
	npm start

lint: ## Run linting
	npm run lint

test: ## Run tests
	npm run test

pre-commit: clean-elm check-compile test ## Run before commiting
	! grep -R "Debug.log" src

setup: ## Initialize the Developer Environment
	npm install
	npm install -g elm-test

.PHONY: clean build run setup clean-elm lint test pre-commit help
