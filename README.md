# SQUID (Sortable Quasi-Unique Identifier)

SQUID is a Go package that generates unique, sortable identifiers. It combines features from CUID (Collision-resistant Unique Identifier) and KSUID (K-Sortable Unique Identifier) to create identifiers that are both collision-resistant and chronologically sortable.

## Features

- **Sortable**: SQUIDs are chronologically sortable, making them ideal for distributed systems where ordering is important.
- **Collision-resistant**: Designed to minimize the chance of collisions, even in distributed environments.
- **Time-based**: Includes a millisecond-precision timestamp for accurate sorting and tracking.
- **Flexible**: Can be used in various applications requiring unique, sortable identifiers.

## Installation

To use SQUID in your Go project, you can install it using `go get`:

```bash
go get github.com/aria3ppp/squid
```

Replace `aria3ppp` with your actual GitHub username.

## Usage

Here's a basic example of how to use SQUID in your Go code:

```go
package main

import (
    "fmt"
    "github.com/aria3ppp/squid"
)

func main() {
    generator, err := squid.NewGenerator()
    if err != nil {
        fmt.Printf("Failed to create SQUID generator: %v\n", err)
        return
    }

    id := generator.New()
    fmt.Printf("Generated SQUID: %s\n", id)

    // Parsing a SQUID
    timestamp, random, counter, err := generator.Parse(id)
    if err != nil {
        fmt.Printf("Failed to parse SQUID: %v\n", err)
        return
    }

    fmt.Printf("Timestamp: %v\n", timestamp)
    fmt.Printf("Random: %x\n", random)
    fmt.Printf("Counter: %d\n", counter)
}
```

## Structure of a SQUID

A SQUID consists of the following components:

1. Timestamp (10 bytes): Milliseconds since the Unix epoch
2. Random data (16 bytes): For uniqueness
3. Counter (2 bytes): To prevent collisions within the same millisecond

The total 28 bytes are then base32 encoded, resulting in a 45-character string.

## Performance

SQUID is designed to be efficient and can generate thousands of unique identifiers per second. However, actual performance may vary depending on your system and use case.

## Testing

To run the tests for SQUID, navigate to the project directory and run:

```bash
go test -v
```

## Contributing

Contributions to SQUID are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by [CUID](https://github.com/ericelliott/cuid) and [KSUID](https://github.com/segmentio/ksuid)
- Thanks to all contributors and users of SQUID

## Contact

If you have any questions or feedback, please open an issue on the GitHub repository.

---

Remember to star this repository if you find it useful!
