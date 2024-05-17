# ekstatic #

ekstatic is tiny workflow engine for go. It enables you to build anything from finite state machines to complex workflows. To achieve that, the library 
is roughly modelled along the concept of a nondeterministic finite automaton: Harnessing the go type system it enables you to transition not only from 
on primitive state to another at a time, but to transform complex state data along with it.

In a nutshell, ekstatic is:

- Lightweight, transparent and maintainable. ekstatic's core funcionality is implemented in less than 200 lines – no third party modules needed.
- Unopinionated. ekstatic doesn't make assumptions about your business case or how you want to integrate a workflow.
- Idiomatic. This library is written in go for go, keeps its feature set minimal and follows the idea to offer only one way to do a thing.
- Tried and true. Every part of it is covered with tests.

## Roadmap

This project is still in an early development stage and might undergo major changes. From this point on, these are the next steps:

- [ ] Finish [ToDos](todo.md)
- [x] Publish contribution guidelines
- [ ] Refine with contributions
- [ ] release v1.0.0 with stable API

## Contributions ##

If you want to make an addition, please link it to a previously opened ticket:
- If you have found a bug and are looking to fix it, reference a bug report or create one if needed.
- If you want to add a feature that is not listed in the [ToDos](todo.md), link it to an existing discussion.

When you think that the project would benefit from a feature, please consider the minimalist concept this library is built around:
Does it fit the project's scope? Is it really not their yet or just another way to achieve something already implemented? Please
[open a new discussion](https://github.com/Metamogul/ekstatic/discussions/new/choose) first to present your feature.

## Reporting bugs ##

If you have found a bug and report it, please [create an issue for it](https://github.com/Metamogul/ekstatic/issues/new/choose) and include 
a description to reproduce the bug, of the expected behavior, the actual behavior, and the affected version of the module.

