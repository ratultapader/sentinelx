package intelligence

type MitreMapping struct {
	EventType    string
	Tactic       string
	Technique    string
	TechniqueID  string
	SubTechnique string
	Confidence   float64
}

type EnrichedMitreContext struct {
	EventType    string
	Tactic       string
	Technique    string
	TechniqueID  string
	SubTechnique string
	Confidence   float64
	Raw          *MitreMapping
}