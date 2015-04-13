package shared

import ()

type Exportable interface {
	Export(groupName string)
}
