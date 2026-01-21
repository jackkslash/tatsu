package harness

type Harness interface {
	Name() string
	IsAvailable() bool
}
