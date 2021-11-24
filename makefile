MAKEFLAGS  := --silent --always-make
PAR        := $(MAKE) -j 128
VERB       := $(if $(filter $(verb), true), -v,)
FAIL       := $(if $(filter $(fail), false),, -failfast)
SHORT      := $(if $(filter $(short), true), -short,)
GO_FLAGS   := -tags=$(tags)
TEST_FLAGS := -count=1 $(VERB) $(FAIL) $(SHORT)
TEST       := test $(TEST_FLAGS) -timeout=1s -run=$(run)
BENCH      := test $(TEST_FLAGS) -run=- -bench=$(or $(run),.) -benchmem
WATCH      := watchexec -r -c -d=0 -n

default: test_w

watch:
	$(PAR) test_w lint_w

test_w:
	gow -c -v $(TEST) $(GO_FLAGS)

test:
	go $(TEST) $(GO_FLAGS)

bench_w:
	gow -c -v $(BENCH) $(GO_FLAGS)

bench:
	go $(BENCH) $(GO_FLAGS)

lint_w:
	$(WATCH) -- $(MAKE) lint

lint:
	golangci-lint run
	echo [lint] ok

prof:
	go $(BENCH) $(GO_FLAGS) -cpuprofile cpu.prof -memprofile mem.prof

cpu:
	go tool pprof -web cpu.prof

mem:
	go tool pprof -web mem.prof

clean:
	rm -f cpu.prof mem.prof oas.test
