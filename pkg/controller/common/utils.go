/*
 * Copyright 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package common

import (
	"context"
	"fmt"
	"hash/fnv"
	"strconv"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	extapi "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gardener/gardener-extension-shoot-dns-service/pkg/service"
)

const (
	ANNOTATION_OPERATION         = constants.GardenerOperation
	ANNOTATION_OPERATION_MIGRATE = constants.GardenerOperationMigrate
	ANNOTATION_OPERATION_RESTORE = constants.GardenerOperationRestore
)

func CopyMap(m map[string]string) map[string]string {
	if m == nil {
		return nil
	}
	r := map[string]string{}
	for k, v := range m {
		r[k] = v
	}
	return r
}

func FindExtension(ctx context.Context, reader client.Reader, namespace string) (*extapi.Extension, error) {
	list := &extapi.ExtensionList{}
	if err := reader.List(ctx, list, client.InNamespace(namespace)); err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	for _, e := range list.Items {
		if e.Spec.Type == service.ExtensionType {
			return &e, nil
		}
	}

	return nil, nil
}

func IsMigrating(ex *extensionsv1alpha1.Extension) bool {
	if ex.Annotations == nil {
		return false
	}
	return ex.Annotations[ANNOTATION_OPERATION] == ANNOTATION_OPERATION_MIGRATE
}

func IsRestoring(ex *extensionsv1alpha1.Extension) bool {
	if ex.Annotations == nil {
		return false
	}
	return ex.Annotations[ANNOTATION_OPERATION] == ANNOTATION_OPERATION_RESTORE
}

// ShortenID shortens an identifier longer than maxlen characters by cutting the string
// and adding a hash suffix so that total length is maxlen. Identifiers are preserved
// if length < maxlen.
func ShortenID(id string, maxlen int) string {
	if maxlen < 16 {
		panic("maxlen < 16 for shortenID")
	}
	if len(id) <= maxlen {
		return id
	}

	hash := fnv.New64()
	_, _ = hash.Write([]byte(id))
	hashstr := strconv.FormatUint(hash.Sum64(), 36)
	return fmt.Sprintf("%s-%s", id[:62-len(hashstr)], hashstr)
}
