package debugger

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

const HelpMessage = `Available commands:
  list              - List all variables and their values.
  get <variable>    - Get the value of a specific variable.
  set <variable> <value> - Set the value of a variable.
  continue, c       - Continue execution.
  quit, q           - Quit the program.`

var mutex sync.Mutex

func Breakpoint(variables map[string]interface{}) {
	_, file, line, _ := runtime.Caller(1)

	fmt.Printf("🛑 Breakpoint reached at %s:%d\n", file, line)
	fmt.Println("Type 'help' for available commands.")

	handleDebugSession(variables)
}

func GoroutineBreakpoint(variables map[string]interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	mutex.Lock()

	defer mutex.Unlock()

	_, file, line, _ := runtime.Caller(1)

	fmt.Printf("🛑 Goroutine Breakpoint reached at %s:%d\n", file, line)
	fmt.Println("Type 'help' for available commands.")

	commandChan := make(chan string)

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("debug> ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			commandChan <- input
		}
	}()

	for {
		select {
		case input := <-commandChan:
			handleCommand(input, variables)
			if input == "continue" || input == "c" {
				return
			}
		}
	}
}

func handleCommand(input string, variables map[string]interface{}) {
	parts := strings.Split(input, " ")

	switch parts[0] {
	case "help":
		printHelp()
	case "list":
		listVariables(variables)
	case "get":
		if len(parts) < 2 {
			fmt.Println("Usage: get <variable>")
		} else {
			getVariable(variables, parts[1])
		}
	case "set":
		if len(parts) < 3 {
			fmt.Println("Usage: set <variable> <value>")
		} else {
			setVariable(variables, parts[1], parts[2])
		}
	case "continue", "c":
		return
	case "quit", "q":
		os.Exit(0)
	default:
		fmt.Println("Unknown command. Type 'help' for available commands.")
	}
}

func handleDebugSession(variables map[string]interface{}) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("debug> ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		parts := strings.Split(input, " ")

		switch parts[0] {
		case "help":
			printHelp()
		case "list":
			listVariables(variables)
		case "get":
			if len(parts) < 2 {
				fmt.Println("Usage: get <variable>")
			} else {
				getVariable(variables, parts[1])
			}
		case "set":
			if len(parts) < 3 {
				fmt.Println("Usage: set <variable> <value>")
			} else {
				setVariable(variables, parts[1], parts[2])
			}
		case "continue", "c":
			return
		case "quit", "q":
			os.Exit(0)
		default:
			fmt.Println("Unknown command. Type 'help' for available commands.")
		}
	}
}

func printHelp() {
	fmt.Println(HelpMessage)
}

func listVariables(variables map[string]interface{}) {
	for name, value := range variables {
		fmt.Printf(" - %s: %v (type: %s)\n", name,
			reflect.ValueOf(value).Elem(), reflect.TypeOf(value))
	}
}

func getVariable(variables map[string]interface{}, name string) {
	if value, ok := variables[name]; ok {
		fmt.Printf("%s: %v (type: %s)\n", name,
			reflect.ValueOf(value).Elem(), reflect.TypeOf(value))
	} else {
		fmt.Printf("Variable '%s' not found.\n", name)
	}
}

func setVariable(variables map[string]interface{}, name, newValue string) {
	if variable, ok := variables[name]; ok {
		switch v := variable.(type) {
		case *int:
			if newInt, err := strconv.Atoi(newValue); err == nil {
				*v = newInt
				fmt.Printf("Variable '%s' updated to %d\n", name, newInt)
			} else {
				fmt.Println("Invalid value for type int.")
			}
		case *string:
			*v = newValue
			fmt.Printf("Variable '%s' updated to %s\n", name, newValue)
		case *float64:
			if newFloat, err := strconv.ParseFloat(newValue, 64); err == nil {
				*v = newFloat
				fmt.Printf("Variable '%s' updated to %f\n", name, newFloat)
			} else {
				fmt.Println("Invalid value for type float64.")
			}
		default:
			fmt.Println("Unsupported type for modification.")
		}
	} else {
		fmt.Printf("Variable '%s' not found.\n", name)
	}
}
