package node

import (
	"github.com/zhengyangfeng00/woodenox/iterator"
)

type Processor interface {
	ProcessItem(stream string, item iterator.Item)
}

type ProcessorUtils interface {
	ProduceItem(stream string, item iterator.Item)
}
