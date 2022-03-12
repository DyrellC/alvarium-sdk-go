package annotators

import (
	"context"
	"encoding/json"
	"github.com/dyrellc/alvarium-sdk-go/pkg/config"
	"github.com/dyrellc/alvarium-sdk-go/pkg/contracts"
	"github.com/dyrellc/alvarium-sdk-go/pkg/interfaces"
	"os"
)

// AutomationAnnotator is used to validate whether the data being provided was provided from a local file instance or
// automated through a retrieval api
type AutomationAnnotator struct {
	hash contracts.HashType
	kind contracts.AnnotationType
	sign config.SignatureInfo
}

func NewAutomationAnnotator(cfg config.SdkInfo) interfaces.Annotator {
	a := AutomationAnnotator{}
	a.hash = cfg.Hash.Type
	a.kind = contracts.AnnotationAuto
	a.sign = cfg.Signature
	return &a
}

func (a *AutomationAnnotator) Do(_ctx context.Context, data []byte) (contracts.Annotation, error) {
	key := deriveHash(a.hash, data)
	hostname, _ := os.Hostname()

	satisfied := false
	var sheet sheetReading
	err := json.Unmarshal(data, &sheet)
	if err != nil {
		var sensor sensorReading
		err := json.Unmarshal(data, &sensor)
		if err != nil {
			return contracts.Annotation{}, err
		}
		satisfied = true
	}

	annotation := contracts.NewAnnotation(string(key), a.hash, hostname, a.kind, satisfied)
	signed, err := signAnnotation(a.sign.PrivateKey, annotation)
	if err != nil {
		return contracts.Annotation{}, err
	}
	annotation.Signature = string(signed)
	return annotation, nil
}

// We need to determine of the reading is formatted as a SheetReading or a SensorReading for the annotation
type sheetReading struct {
	SheetId   string `json:"sheetId,omitempty"`
	Value     string `json:"value,omitempty"`
}

type sensorReading struct {
	SensorId   string        `json:"sensorId,omitempty"`
	Value      ReadingValues `json:"value,omitempty"`
}


type ReadingValues struct {
	Context  string      `json:"@odata.context"`
	Values []SensorValue  `json:"value"`
}

type SensorValue struct {
	FQN string      `json:"FQN"`
	DateTime string    `json:"DateTime"`
	OpcQuality uint32   `json:"OpcQuality"`
	Value float64    `json:"Value"`
	Text string      `json:"Text"`
}
