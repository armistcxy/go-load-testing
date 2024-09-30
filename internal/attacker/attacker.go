package attacker

import (
	"log"
	"log/slog"
	"os"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type Attacker struct {
	attacker       vegeta.Attacker
	targeter       vegeta.Targeter
	saveResultPath string
}

func NewAttacker(targeter vegeta.Targeter, saveResultPath string) *Attacker {
	return &Attacker{
		attacker:       *vegeta.NewAttacker(),
		targeter:       targeter,
		saveResultPath: saveResultPath,
	}

}

func (a *Attacker) Attack(freq int, per time.Duration, duration time.Duration) {
	var metrics vegeta.Metrics
	rate := vegeta.Rate{
		Freq: freq,
		Per:  per,
	}

	binayResultFile, err := os.OpenFile(a.saveResultPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}

	resultEncoder := vegeta.NewEncoder(binayResultFile)

	for res := range a.attacker.Attack(a.targeter, rate, duration, "Attack") {
		metrics.Add(res)
		if err := resultEncoder.Encode(res); err != nil {
			slog.Error("failed to encode attack result", "error", err)
		}
	}

	metrics.Close()

	log.Printf("Success Rate: %.2f%%\n", metrics.Success*100)
}
