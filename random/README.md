# Random Starlark module

Generates and manipulates pseudo-random numbers

## Usage

    predeclared := starlark.StringDict{
        "random": random.New()
        ...
    }
    starlark.ExecFile(thread, filename, nil, predeclared)

## Supported functions

    random.seed()

Initialize the random number generator.

    random.randint(a int, b int) int

Return a random integer N such that a <= N <= b for a <= b and b <= N <= a for b < a. Uses Rand.Int63n Go function.

    random.random() float

Return the next random floating point number in the range [0.0, 1.0).

    random.uniform(a float, b float) float

Return a random floating point number N such that a <= N <= b for a <= b and b <= N <= a for b < a.