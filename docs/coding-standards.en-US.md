# Go Coding Standards

## 1. Naming Conventions

### Basic Principles

- Names must begin with a letter (A-Z, a-z) or an underscore, and can be followed by letters, underscores, or numbers (0-9).
- Special characters such as @, $, % are strictly prohibited in names.
- Go is case-sensitive. Identifiers starting with an uppercase letter are public (exported), while those starting with a lowercase letter are private (internal to the package).

### 1.1. Package Naming

- The package name must be consistent with the directory name. Choose a name that is concise, meaningful, and does not conflict with the standard library.
- Package names must be all lowercase. Multiple words can be separated by underscores or use mixed case (camelCase is not recommended).

```go
package demo
package main
```

### 1.2. File Naming

- Filenames should be clear, concise, and easy to understand.
- They must use lowercase letters, with words separated by underscores.

```go
my_test.go
```

### 1.3. Struct Naming

- Struct names must use CamelCase. The first letter's case depends on the desired access control (public/private).
- Struct declarations and initializations must use a multi-line format, as shown below:

```go
// Multi-line declaration
type User struct {
	Username string
	Email    string
}

// Multi-line initialization
user := User{
	Username: "admin",
	Email:    "admin@example.com",
}
```

### 1.4. Interface Naming

- Interface names must use CamelCase. The first letter's case depends on the desired access control.
- Interfaces with a single method should be named by the method name plus an "er" suffix (e.g., `Reader`, `Writer`).

```go
type Reader interface {
	Read(p []byte) (n int, err error)
}
```

### 1.5. Variable Naming

- Variable names must use CamelCase. The first letter's case depends on the desired access control.
- Rules for handling special nouns (e.g., API, ID):
  - If the variable is private and the special noun is the first word, use lowercase (e.g., `apiClient`).
  - In other cases, keep the original capitalization (e.g., `APIClient`, `repoID`, `UserID`).
  - Incorrect example: `UrlArray`. Correct examples: `urlArray` or `URLArray`.
- Boolean variable names must start with `Has`, `Is`, `Can`, or `Allow`.

```go
var isExist bool
var hasConflict bool
var canManage bool
var allowGitHook bool
```

### 1.6. Constant Naming

- Constant names must use CamelCase and be prefixed according to their category.

```go
// HTTP method constants
const (
	MethodGET  = "GET"
	MethodPOST = "POST"
)
```

- Enumerated constants should also follow this convention:

```go
type Scheme string

const (
    SchemeHTTP  Scheme = "http"
    SchemeHTTPS Scheme = "https"
)
```

### 1.7. Keywords

Go keywords: `break`, `case`, `chan`, `const`, `continue`, `default`, `defer`, `else`, `fallthrough`, `for`, `func`, `go`, `goto`, `if`, `import`, `interface`, `map`, `package`, `range`, `return`, `select`, `struct`, `switch`, `type`, `var`

## 2. Commenting Standards

Go supports C-style comments: `/**/` and `//`.

- Line comments (`//`) are the most common form.
- Block comments (`/* */`) are mainly used for package comments and cannot be nested. They are typically used for documentation or commenting out large blocks of code.

### 2.1. Package Comments

- Every package must have a package comment preceding the `package` clause.
- If a package has multiple files, the package comment only needs to appear in one file (preferably the one with the same name as the package).
- The package comment must include the following information in order:
  - A brief introduction to the package (name and functionality).
  - Creator information, format: `Creator: [GitHub Username]`
  - Creation date, format: `Created: YYYY-MM-DD`

```go
// Package biz_err provides business error codes and messages.
// Creator: Done-0
// Created: 2025-07-01
```

### 2.2. Struct and Interface Comments

- Every custom struct or interface must have a comment on the line preceding its definition.
- The format is: `// [Struct/Interface Name], [Description]`.
- Each field of a struct must have a comment, placed after the field and aligned.
- Example: `User` is the struct name, and `user object, defines basic user information` is the description.

```go
// User, user object, defines basic user information.
type User struct {
    Username  string // Username
    Email     string // Email
}
```

### 2.3. Function and Method Comments (Optional)

Each function or method should have a comment that includes (in order):

- A brief description: Start with the function name, followed by a space and the description.
- Parameters: One per line, starting with the parameter name, followed by `: ` and the description.
- Return values: One per line.

```go
// NewAttrModel is a factory method for the attribute data layer.
// Parameters:
//      ctx: Context information.
//
// Returns:
//      *AttrModel: A pointer to the attribute model.
func NewAttrModel(ctx *common.Context) *AttrModel {
}
```

### 2.4. Code Logic Comments

- Add comments to explain critical sections or complex logic.

```go
// Batch read attributes from Redis. For IDs not found,
// record them in an array to be read from the DB later.
// ... subsequent code ...
```

### 2.5. Comment Style

- Use English for all comments.
- A space must separate Chinese and English characters, including between Chinese characters and English punctuation.
- It is recommended to use single-line comments exclusively.
- A single-line comment should not exceed 120 characters.

## 3. Code Style

### 3.1. Indentation and Line Breaks

- Indentation must be formatted with the `gofmt` tool (using tabs).
- Each line of code should not exceed 120 characters. Longer lines should be broken and formatted elegantly.

> In Goland, you can format code with the shortcut `Control + Alt + L`.

### 3.2. Statement Termination

- Go does not require semicolons at the end of statements; a new line implies a new statement.
- If multiple statements are on the same line, they must be separated by semicolons.

```go
package main

func main() {
  var a int = 5; var b int = 10
  // Multiple statements on the same line must be separated by semicolons.
  c := a + b; fmt.Println(c)
}
```

- While multi-statement lines are allowed for simple code, single-statement lines are recommended.

```go
package main

func main() {
    var a int = 5
    var b int = 10

    c := a + b
    fmt.Println(c)
}
```

### 3.3. Braces and Spaces

- The opening brace must not be on a new line (enforced by Go syntax).
- A space must be present between all operators and their operands.

```go
// Correct
if a > 0 {
    // Code block
}

// Incorrect
if a>0  // Spaces should be around >
{       // Opening brace cannot be on a new line
    // Code block
}
```

### 3.4. Import Standards

- For a single package, use the parenthesized format:

```go
import (
    "fmt"
)
```

- When importing multiple packages, they should be grouped in the following order, separated by blank lines:

  1. Standard library packages
  2. Third-party packages
  3. Internal project packages

- Aliased imports, blank imports (`_`), and dot imports (`.`) should be placed within their respective groups and sorted alphabetically. **Note:** Dot imports can pollute the current namespace and should be used with caution.

```go
import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"

	"github.com/Done-0/gin-scaffold/internal/global"
	"github.com/Done-0/gin-scaffold/pkg/vo"

    _ "github.com/go-sql-driver/mysql"   // Blank import (third-party)
	customname "github.com/pkg/errors"   // Aliased import (third-party)
    . "github.com/alecthomas/kingpin/v2" // Dot import (third-party)
)
```

- Do not use relative paths to import external packages:

```go
// Incorrect
import "../net" // Relative imports of external packages are forbidden.

// Correct
import "github.com/repo/proj/src/net"
```

- If the package name and import path do not match, use an alias:

```go
// Incorrect
import "github.com/Done-0/gin-scaffold/gin-scaffold/internal/model/account" // The actual package name is `model`.

// Correct
import model "github.com/Done-0/gin-scaffold/gin-scaffold/internal/model/account" // Use `model` as an alias.
```

### 3.5. Error Handling

- Never discard an error returned from a call. Do not use `_` to discard errors; they must all be handled.
- Error handling principles:
  - Return immediately upon error (return early).
  - Do not use `panic` unless you know the exact consequences.
  - Error messages in English must be all lowercase and should not end with punctuation.
  - Errors must be handled in a separate error flow.

```go
// Incorrect
if err != nil {
    // Error handling
} else {
    // Normal code
}

// Correct
if err != nil {
    // Error handling
    return // or continue, etc.
}
// Normal code
```

### 3.6. Testing Standards

- Test filenames must end with `_test.go` (e.g., `example_test.go`).
- Test function names must start with `Test` (e.g., `TestExample`).
- Every important function should have a test case submitted with the official code for regression testing.

## 4. Common Tools

Go provides several tools to help developers follow coding standards:

### gofmt

Most formatting issues can be resolved with `gofmt`. It automatically formats code to ensure consistency with official Go standards. All formatting questions are settled by the output of `gofmt`.

### goimports

`goimports` is highly recommended. It builds on `gofmt` by automatically adding and removing package imports.

```bash
go get golang.org/x/tools/cmd/goimports
```

### go vet

The `vet` tool statically analyzes source code for various issues, such as dead code, prematurely returned logic, and non-standard struct tags.

```bash
go get golang.org/x/tools/cmd/vet
```

Usage:

```bash
go vet .
```
