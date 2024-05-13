## TODOs: ##

- [x] variadic functions as transitions
- [x] transitions can error
- [x] persistance hooks -> wie soll das mit der ID laufen und so? Embedding der Machine in struct, das das macht?
- [x] submachines: if current state type is StateMachine, pass call to PerformTransition through to submachine before processing it
- [x] events: AddTransitionSucceededEvent, AddTransitionFailedEvent -> Event ist irgendeine Funktion die neuen & alten state & input übergeben bekommt
- [x] threadsafety
- [x] structure module according to established standards (https://go.dev/doc/modules/layout)
- [x] more tests + examples
- [x] create go module
- [ ] docs
- [ ] add github CI pipeline
- [ ] comparitive benchmarks to other frameworks
- [ ] add linters