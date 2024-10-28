package debugger

import (
	"bytes"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestListVariables(t *testing.T) {
	var1 := 10
	var2 := "test"
	var3 := 3.14

	variables := map[string]interface{}{
		"intVar":    &var1,
		"stringVar": &var2,
		"floatVar":  &var3,
	}

	output := captureOutput(func() {
		listVariables(variables)
	})

	if !strings.Contains(output, "intVar: 10 (type: *int)") {
		t.Errorf("Expected 'intVar' to be listed, got: %s", output)
	}
	if !strings.Contains(output, "stringVar: test (type: *string)") {
		t.Errorf("Expected 'stringVar' to be listed, got: %s", output)
	}
	if !strings.Contains(output, "floatVar: 3.14 (type: *float64)") {
		t.Errorf("Expected 'floatVar' to be listed, got: %s", output)
	}
}

func TestGetVariable(t *testing.T) {
	var1 := 42
	variables := map[string]interface{}{
		"testVar": &var1,
	}

	output := captureOutput(func() {
		getVariable(variables, "testVar")
	})

	expected := "testVar: 42 (type: *int)"
	if !strings.Contains(output, expected) {
		t.Errorf("Expected '%s', got: %s", expected, output)
	}
}

func TestSetVariable(t *testing.T) {
	var1 := 10
	var2 := "old"
	var3 := 3.14

	variables := map[string]interface{}{
		"intVar":    &var1,
		"stringVar": &var2,
		"floatVar":  &var3,
	}

	setVariable(variables, "intVar", "20")
	setVariable(variables, "stringVar", "new")
	setVariable(variables, "floatVar", "2.71")

	if var1 != 20 {
		t.Errorf("Expected intVar to be 20, got: %d", var1)
	}
	if var2 != "new" {
		t.Errorf("Expected stringVar to be 'new', got: %s", var2)
	}
	if var3 != 2.71 {
		t.Errorf("Expected floatVar to be 2.71, got: %f", var3)
	}
}

func TestSetVariableInvalid(t *testing.T) {
	var1 := 10
	variables := map[string]interface{}{
		"intVar": &var1,
	}

	output := captureOutput(func() {
		setVariable(variables, "intVar", "invalid")
	})

	expected := "Invalid value for type int."
	if !strings.Contains(output, expected) {
		t.Errorf("Expected '%s', got: %s", expected, output)
	}
}
func TestBreakpoint(t *testing.T) {
	var1 := 42
	variables := map[string]interface{}{
		"testVar": &var1,
	}

	setFakeStdin("continue\n")

	output := captureOutput(func() {
		Breakpoint(variables)
	})

	if !strings.Contains(output, "🛑 Breakpoint reached at") {
		t.Errorf("Expected breakpoint message, got: %s", output)
	}
}

func TestGoroutineBreakpoint(t *testing.T) {
	var1 := 100
	variables := map[string]interface{}{
		"testVar": &var1,
	}
	var wg sync.WaitGroup

	wg.Add(1)
	setFakeStdin("continue\n")

	output := captureOutput(func() {
		GoroutineBreakpoint(variables, &wg)
	})
	wg.Wait()

	if !strings.Contains(output, "🛑 Goroutine Breakpoint reached at") {
		t.Errorf("Expected goroutine breakpoint message, got: %s", output)
	}
}

func setFakeStdin(input string) {
	r, w, _ := os.Pipe()

	_, _ = w.WriteString(input)
	_ = w.Close()

	os.Stdin = r
}

func captureOutput(f func()) string {
	r, w, _ := os.Pipe()
	defer r.Close()

	stdout := os.Stdout
	os.Stdout = w

	outputChan := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = buf.ReadFrom(r)
		outputChan <- buf.String()
	}()

	f()
	_ = w.Close()
	os.Stdout = stdout

	return <-outputChan
}
