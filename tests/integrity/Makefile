GO         ?= go
JAVA       ?= java
LUCENE_JAR ?= lib/lucene-core-10.4.0.jar
INDEX_DIR  ?= build/index
TEST_RUN   ?= TestCheckIndex

.DEFAULT_GOAL := test

.PHONY: test
test: check-tools
	$(GO) test -v -count=1 -run $(TEST_RUN) .

.PHONY: check
check: check-tools
	INTEGRITY_INDEX_DIR=$(INDEX_DIR) $(GO) test -v -count=1 -run $(TEST_RUN) .

.PHONY: checkindex
checkindex: check-java
	$(JAVA) -cp $(LUCENE_JAR) org.apache.lucene.index.CheckIndex $(INDEX_DIR)

.PHONY: check-tools
check-tools: check-go check-java

.PHONY: check-go
check-go:
	@command -v $(GO) >/dev/null 2>&1 || { echo "error: '$(GO)' not found on PATH"; exit 1; }

.PHONY: check-java
check-java:
	@command -v $(JAVA) >/dev/null 2>&1 || { echo "error: '$(JAVA)' not found on PATH (Lucene 10.4 needs Java 21+)"; exit 1; }
	@test -f $(LUCENE_JAR) || { echo "error: lucene jar not found at $(LUCENE_JAR)"; exit 1; }

.PHONY: clean
clean:
	rm -rf build

