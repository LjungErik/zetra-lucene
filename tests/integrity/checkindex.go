package integrity

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const luceneCoreJar = "lucene-core-10.4.0.jar"

const checkIndexMainClass = "org.apache.lucene.index.CheckIndex"

type CheckIndexResult struct {
	Output   string
	ExitCode int
	OK       bool
}

func LocateJar() (string, error) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("integrity: unable to resolve source file location")
	}

	jar := filepath.Join(filepath.Dir(thisFile), "lib", luceneCoreJar)
	if _, err := os.Stat(jar); err != nil {
		return "", fmt.Errorf("integrity: lucene-core jar not found at %q: %w", jar, err)
	}

	return jar, nil
}

func JavaAvailable() bool {
	_, err := exec.LookPath("java")
	return err == nil
}

func RunCheckIndex(ctx context.Context, indexDir string) (*CheckIndexResult, error) {
	jar, err := LocateJar()
	if err != nil {
		return nil, err
	}

	absIndexDir, err := filepath.Abs(indexDir)
	if err != nil {
		return nil, fmt.Errorf("integrity: cannot resolve index dir %q: %w", indexDir, err)
	}

	args := []string{"-cp", jar}
	if extra := strings.Fields(os.Getenv("JAVA_OPTS")); len(extra) > 0 {
		args = append(args, extra...)
	}
	args = append(args, checkIndexMainClass, absIndexDir)

	cmd := exec.CommandContext(ctx, "java", args...)
	outBytes, runErr := cmd.CombinedOutput()

	result := &CheckIndexResult{
		Output:   string(outBytes),
		ExitCode: -1,
	}
	if cmd.ProcessState != nil {
		result.ExitCode = cmd.ProcessState.ExitCode()
	}

	if runErr != nil {
		var exitErr *exec.ExitError
		if errors.As(runErr, &exitErr) {
			result.OK = false
			return result, nil
		}
		return result, fmt.Errorf("integrity: failed to run %s: %w", checkIndexMainClass, runErr)
	}

	result.OK = result.ExitCode == 0
	return result, nil
}
