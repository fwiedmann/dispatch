package dispatch

import "github.com/fwiedmann/dispatch/pkg/code"

func NewCallTaker(analyzer analyzer, dispatcher dispatcher, owners []code.Owner) CallTaker {
	return CallTaker{
		analyzer:   analyzer,
		dispatcher: dispatcher,
		owners:     owners,
	}
}

type analyzer interface {
	Analyze() ([]code.Info, error)
}

// Note holds the required information which a dispatcher will use to send the correct information.
type Note struct {
	code.Info
	code.Owner
}

type dispatcher interface {
	Dispatch([]Note) error
}

// CallTaker orchestrates the analyzing and the correct dispatching
type CallTaker struct {
	analyzer
	dispatcher
	owners []code.Owner
}

func (calltaker *CallTaker) Dispose() error {
	ownerRefs, err := calltaker.Analyze()
	if err != nil {
		return err
	}

	return calltaker.Dispatch(calltaker.findOwners(ownerRefs))
}

func (calltaker *CallTaker) findOwners(ownerRefs []code.Info) []Note {
	notes := make([]Note, 0)
	for _, ref := range ownerRefs {
		for _, owner := range calltaker.owners {
			if ref.OwnerId == owner.Id {
				notes = append(notes, Note{
					Info:  ref,
					Owner: owner,
				})
			}
		}
	}
	return notes
}
