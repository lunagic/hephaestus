package state

func NewHephaestus() (*Hephaestus, error) {
	return &Hephaestus{}, nil
}

type Hephaestus struct {
	DefaultPort int
}

func (hephaestus *Hephaestus) GitCleanExcludes() []string {
	return []string{
		".env.local",
		".env.*.local",
	}
}

func (hephaestus *Hephaestus) HTTPEnabled() bool {
	return hephaestus.DefaultPort != 0
}
