## About ##

ekstatic is tiny workflow engine for go. It enables you to build anything from finite state machines to complex workflows. To achieve that, the library is roughly modelled along the concept of a nondeterministic finite automaton: Harnessing the go type system it enables you to transition not only from on primitive state to another at a time, but to transform complex state data along with it.

In a nutshell, ekstatic is:

- Lightweight, transparent and maintainable. ekstatic's core funcionality is implemented in less than 300 lines – no third party modules needed.
- Unopinionated. ekstatic doesn't make assumptions about your business case or how you want to integrate a workflow.
- Idiomatic. This library is written in go for go, keeps its feature set minimal and follows the idea to offer only one way to do a thing.
- Tried and true. Every part of it is covered with extensive tests.

## Concurrency ##

ekstatic's `WorkflowInstance` is concurrency safe and can safely be used from multiple go routines at once. The `Workflow` struct however is intended to be configured once at application startup and to be kept around for the apps lifetime to spawn `WorkflowInstance` structs on demand. Therefore `WorkflowInstance` is not concurrency safe except for its `New` receiver method.

## Restoring workflow instances ##

By design, instances are completely persistable and restorable. ... (Explanation, example, blabla)

## Roadmap ##

## Contributing ##

Contributions are always welcome. Feel free to fork the project and open a PR. If you want to make a feature addition, please consider the design principles this library is built around. Carefully think if it's really not their yet or just another way to achieve something already implemented, and if the functionality really really fits the scope of the project.

## 