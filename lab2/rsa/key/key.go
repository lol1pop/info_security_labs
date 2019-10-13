package key

import "io"

type Key interface {
	IsValid() bool
	Save(writer io.Writer) error
}
