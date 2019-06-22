package connections

type Connection interface {
	Open() error
	Close()
	Send(cmd interface{}) (RemoteCommandResults, error)
}

type RemoteCommandResults struct {
	Stdout string // Stdout written by the command
	Stderr string // Stderr written by the command
}
