## TODOs: ##

- [x] variadic functions as transitions
- [x] transitions can error
- [ ] persistance hooks -> wie soll das mit der ID laufen und so? Embedding der Machine in struct, das das macht?
- [ ] submachines: if current state type is StateMachine, pass call to PerformTransition through to submachine before processing it
- [ ] events: AddTransitionSucceededEvent, AddTransitionFailedEvent -> Event ist irgendeine Funktion die neuen & alten state & input übergeben bekommt
- [x] threadsafety
- [ ] more tests + examples
- [ ] add github CI pipeline
- [ ] docs
- [ ] create go module
- [ ] comparitive benchmarks to other frameworks
- [ ] add linters