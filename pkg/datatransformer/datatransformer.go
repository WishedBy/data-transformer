package datatransformer

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type DataTransformerManager struct {
	sync.Mutex
	instances map[uuid.UUID]*DataTransformer
}

func (d *DataTransformerManager) NewTransformer() uuid.UUID {
	d.Lock()
	defer d.Unlock()
	id := uuid.New()
	d.instances[id] = &DataTransformer{}
	return id
}
func (d *DataTransformerManager) GetTransformer(id uuid.UUID) *DataTransformer {
	d.Lock()
	defer d.Unlock()
	return d.instances[id]
}
func (d *DataTransformerManager) RemoveTransformer(id uuid.UUID) {
	d.Lock()
	defer d.Unlock()
	delete(d.instances, id)
}

type Value struct {
	ValueString   string
	ValueInt      int64
	ValueFloat    float64
	ValueDuration time.Duration
	ValueTime     time.Time
	PrintedValue  string
}

type DataTransformer struct {
	DataPoints []DataPoint
	XAxises    []Axis
	YAxises    []Axis
}
type DataField struct {
	Value
}
type DataPoint struct {
	Fields map[FieldName]DataField
}
type FieldName struct {
}
type Axis struct {
	Labels Labels
	Name   string
}
type Label struct {
	Value
}
type Labels struct {
	ForField FieldName
	Labels   []Label
}
