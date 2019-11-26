# Interfaces

Ways to achieve loose coupling in Go

    - Avoid the use of global state.
    - Use interfaces to describe the behaviour your functions or methods require:

If you want to reduce the coupling a global variable creates

    - Move the relevant variables as fields on structs that need them.
    - Use interfaces to reduce the coupling between the behaviour and the
    implementation of that behaviour:
