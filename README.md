# gravitonctl

## not ready for public use yet (still testing)!

 Launch & use a graviton instance in 5 seconds.

## Usage:
Configure gravitonctl:
```
gravitonctl configure
```

To start (and connect to) a Graviton instance:
```
gravitonctl start {name}
````

To connect to a Graviton instance:
```
gravitonctl connect {name}
```

To list all Graviton instances:
```
gravitonctl list
```

<!--
To stop a Graviton instance
```
gravitonctl stop {name}
``` 
-->

To terminate a Graviton instance:
```
gravitonctl terminate {name}
```

## Todos:
- clean up SSH Connection code
- handle duplicate names more gracefully
- cleaner argument/flag handling
