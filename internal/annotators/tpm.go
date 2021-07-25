/*******************************************************************************
 * Copyright 2021 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/
package annotators

import (
	"context"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/config"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/contracts"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/interfaces"
	"os"
)

// TpmAnnotator is used to attest whether or not the host machine has TPM capability for managing secrets
type TpmAnnotator struct {
	hash contracts.HashType
	kind contracts.AnnotationType
	sign config.SignatureInfo
}

func NewTpmAnnotator(cfg config.SdkInfo) interfaces.Annotator {
	a := TpmAnnotator{}
	a.hash = cfg.Hash.Type
	a.kind = contracts.AnnotationTPM
	a.sign = cfg.Signature
	return &a
}

func (a *TpmAnnotator) Do(ctx context.Context, data []byte) (contracts.Annotation, error) {
	key := deriveHash(a.hash, data)
	hostname, _ := os.Hostname()

	// For now, the satisfied property will always be false until I can get virtual TPM support in the lab
	// TODO: On the other hand, I could implement the TPM check with the understanding it will be false in any
	// deployment without a TPM.
	annotation := contracts.NewAnnotation(key, a.hash, hostname, a.kind, false)
	sig, err := signAnnotation(a.sign.PrivateKey, annotation)
	if err != nil {
		return contracts.Annotation{}, err
	}
	annotation.Signature = string(sig)
	return annotation, nil
}