package network

type data struct {
	id               string
	parent           string
	neighbors        []string
	message          string
	letter           string
	letterOccurrence map[string]int
}

func (d *data) probesDiffusion(id string, parent string, message string) {
	echoCounter := 0
	// check if the probe is an echo
	if contains(d.neighbors, id) {
		echoCounter++
	}
	if echoCounter == 0 {
		return
	}
	d.letterOccurrence[id] = letterCounter(message, d.letter)

	// send probe to all neighbors
	for _, neighbor := range d.neighbors {
		if neighbor != parent {
			// send probe to neighbor
		}
	}

	// Send echo to parent

}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func probesDiffusion(message string, neihgbors []string) {
	for _, neighbor := range neihgbors {
		// Send message to neighbor

	}
}

func (d *data) sendEcho(parent string, letterOccurrence map[string]int) {
	// Send echo to parent
}

func letterCounter(text string, letter string) int {
	count := 0
	for _, char := range text {
		if string(char) == letter {
			count++
		}
	}
	return count
}
