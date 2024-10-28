# go-simple-debugger

A simple debugger library for Go, allowing you to inspect and modify variables during runtime.

## Features

- List all variables and their values.
- Inspect a specific variable.
- Modify variables at runtime.
- Continue or quit the program using commands.

## Installation

Make sure you have [Go installed](https://golang.org/dl/) (version 1.21+).

Clone the repository:

```sh
git clone https://github.com/xxpauloxx/go-simple-debugger.git
cd go-simple-debugger
```

## Usage

You can find a usage example in `cmd/example/main.go`.

Run the example with:

```sh
go run ./cmd/example
```

Run tests:

```sh
go test -v ./debugger
```

### Example Code

Here’s a basic example of how to use the debugger:

```go
package main

import (
	"fmt"
	"sync"

	"github.com/xxpauloxx/go-simple-debugger/debugger"
)


func main() {
	x := 10
	y := "y value"

	fmt.Println("Starting...")
	debugger.Breakpoint(map[string]interface{}{
		"x": &x,
		"y": &y,
	})

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		z := 20
		debugger.GoroutineBreakpoint(map[string]interface{}{
			"z": &z,
		}, &wg)
		fmt.Printf("z value: %d\n", z)
	}()

	wg.Wait()

	fmt.Printf("After debugger: x = %d, y = %s\n", x, y)
}
```

## Available Commands

| Command               | Description                                |
|-----------------------|--------------------------------------------|
| `list`                | Lists all variables and their values.      |
| `get <variable>`      | Shows the value of a specific variable.    |
| `set <variable> <value>` | Modifies the value of a variable.     |
| `continue` or `c`     | Continues program execution.               |
| `quit` or `q`         | Quits the program.                         |

## Project Structure

```plaintext
go-simple-debugger/
│
├── cmd/
│   └── example/
│       └── main.go         # Example usage of the library.
│
├── internal/
│   └── debugger/
│       └── debugger.go         # Core debugger implementation.
│       └── debugge_test_.go    # Core debugger tests implementation.
│
├── go.mod                  # Module and dependencies management.
├── go.sum                  # Dependencies checksum.
├── LICENSE                 # Project license.
└── README.md               # Documentation (this file).
```

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests with improvements.

## License

This project is licensed under the [MIT License](LICENSE).