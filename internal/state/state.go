package state

type State struct {
	Go         *Go
	NPM        *NPM
	Hephaestus *Hephaestus
}

func New() (*State, error) {
	s := &State{}

	{
		// Get Hephaestus
		var err error
		if s.Hephaestus, err = NewHephaestus(); err != nil {
			return nil, err
		}
	}

	{
		// Get Go
		var err error
		if s.Go, err = NewGo(); err != nil {
			return nil, err
		}
	}

	{
		// Get NPM
		var err error
		if s.NPM, err = NewNPM(); err != nil {
			return nil, err
		}
	}

	return s, nil
}
