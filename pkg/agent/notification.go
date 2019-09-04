package agent

import "github.com/n0rad/go-erlog/data"

type NotificationType string

//const Info		//
//const Warn		// label conflict, test failed
//const Error		// new pending failed
//const Critical	// disk lost
//

//const Action	// plug disk, ask label,

type Notify interface {
	Notify(fields data.Fields, message string)
}
