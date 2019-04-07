# Internal Design Note

## Data structure of composite types

```
type Slice struct {
	pointer int
	len     int
	cap     int
}

type Map struct {
	pointer int
	len     int
	cap     int
}

type Interface struct {
	pointer          int
	receiverTypeId   int
	dynamicTypeId    int
}

type MapData struct {
	elements []Element
}

type Element struct {
	key   *interface{}
	value *interface{}
}
```

