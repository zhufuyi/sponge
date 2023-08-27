## krand

Generate random strings, integers, floating point numbers.

<br>

## Example of use

### Generate a random string

```go
    import "github.com/zhufuyi/sponge/pkg/krand"

    /*
	R_NUM = 1      // only number
	R_UPPER = 2   // only capital letters
	R_LOWER = 4  // only lowercase letters
	R_All = 7	       // numbers, upper and lower case letters
    */

	// by | or combining different types
    kind := krand.R_NUM|krand.R_UPPER    // capital letters and numbers

	// a random string of length 10, consisting of upper case letters and numbers
    krand.String(kind, 10)
```

<br>

### Generate random integers

```go
    import "github.com/zhufuyi/sponge/pkg/krand"

    krand.Int(200)            // random number range 0 ~ 200
    krand.Int(1000, 2000)  // random number range 1000 ~ 2000
```

<br>

### Generate random floating point numbers

```go
    import "github.com/zhufuyi/sponge/pkg/krand"

    krand.Float64(1, 200)            // floating point number with 1 decimal point, range 0~200
    krand.Float64(2, 100,1000)            // floating point number with 2 decimal places, range 100~1000
```
