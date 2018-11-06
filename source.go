package appcast

// Sourcer is the interface that wraps the Source methods.
//
// This interface should be embedded by more use-case specific Sourcer
// interfaces.
type Sourcer interface {
	Load() error
	GenerateChecksum(algorithm ChecksumAlgorithm) *Checksum
	GuessProvider()
	Content() []byte
	SetContent(content []byte)
	Checksum() *Checksum
	Provider() Provider
	SetProvider(provider Provider)
	Appcast() Appcaster
	SetAppcast(appcast Appcaster)
}

// Source represents an appcast source which holds the information about the
// retrieved appcast.
//
// This struct should be embedded by more use-case specific Source structs.
type Source struct {
	// content specifies an appcast data from which the provider will be guessed
	// and releases will be extracted.
	content []byte

	// checksum specifies the Checksum pointer that hold a checksum data about
	// the Source.content.
	//
	// This value is set by calling Source.GenerateChecksum.
	checksum *Checksum

	// provider specifies the one of the supported providers guessed using the
	// Source.GuessProvider method.
	provider Provider

	// appcast specifies the provider-specific appcast guessed for the current
	// Source.
	appcast Appcaster
}

// Load should implement a way of loading an appcast content into the
// Source.content depending on the chosen Sourcer.
func (s *Source) Load() error {
	panic("implement me")
}

// GenerateChecksum creates a new Checksum instance pointer based on the
// provided algorithm and sets it as a Source.checksum.
func (s *Source) GenerateChecksum(algorithm ChecksumAlgorithm) *Checksum {
	c := NewChecksum(algorithm, s.content)
	s.checksum = c
	return c
}

// GuessProvider attempts to guess the supported provider based on the
// Source.content. By default returns an Unknown provider.
func (s *Source) GuessProvider() {
	s.provider = GuessProviderByContent(s.content)
}

// Content is a Source.content getter.
func (s *Source) Content() []byte {
	return s.content
}

// SetContent is a Source.content setter.
func (s *Source) SetContent(content []byte) {
	s.content = content
}

// Checksum is a Source.checksum getter.
func (s *Source) Checksum() *Checksum {
	return s.checksum
}

// Provider is a Source.provider getter.
func (s *Source) Provider() Provider {
	return s.provider
}

// SetProvider is a Source.provider setter.
func (s *Source) SetProvider(provider Provider) {
	s.provider = provider
}

// Appcast is a Source.appcast getter.
func (s *Source) Appcast() Appcaster {
	return s.appcast
}

// SetAppcast is a Source.appcast setter.
func (s *Source) SetAppcast(appcast Appcaster) {
	s.appcast = appcast
}
