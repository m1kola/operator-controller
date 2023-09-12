/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	catalogd "github.com/operator-framework/catalogd/api/core/v1alpha1"
	"github.com/operator-framework/operator-registry/alpha/action"
	"github.com/operator-framework/operator-registry/alpha/declcfg"
)

type catalogmetadataReader struct {
	renderer      action.Render
	clientBuilder *fake.ClientBuilder
}

func newCatalogmetadataReader(indexRef string, clientBuilder *fake.ClientBuilder) *catalogmetadataReader {
	return &catalogmetadataReader{
		clientBuilder: clientBuilder,
		renderer: action.Render{
			Refs:           []string{indexRef},
			AllowedRefMask: action.RefDCImage | action.RefDCDir,
		},
	}
}

func (es *catalogmetadataReader) run(ctx context.Context) (*fake.ClientBuilder, error) {
	cfg, err := es.renderer.Run(ctx)
	if err != nil {
		return nil, err
	}

	// TODO: update fake catalog name string to be catalog name once we support multiple catalogs in CLI
	catalogName := "offline-catalog"

	objs := []client.Object{
		&catalogd.Catalog{
			ObjectMeta: metav1.ObjectMeta{Name: catalogName},
		},
	}

	for _, bundle := range cfg.Bundles {
		content, err := json.Marshal(bundle)
		if err != nil {
			return nil, err
		}

		obj := &catalogd.CatalogMetadata{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("%s-%s-%s-%s", catalogName, declcfg.SchemaBundle, bundle.Package, bundle.Name),
				Labels: map[string]string{
					"catalog":       catalogName,
					"name":          bundle.Name,
					"package":       bundle.Package,
					"packageOrName": bundle.Package,
					"schema":        declcfg.SchemaBundle,
				},
			},
			Spec: catalogd.CatalogMetadataSpec{
				Catalog: corev1.LocalObjectReference{Name: catalogName},
				Name:    bundle.Name,
				Package: bundle.Package,
				Schema:  declcfg.SchemaBundle,
				Content: content,
			},
		}

		objs = append(objs, obj)
	}

	for _, channel := range cfg.Channels {
		content, err := json.Marshal(channel)
		if err != nil {
			return nil, err
		}

		obj := &catalogd.CatalogMetadata{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("%s-%s-%s-%s", catalogName, declcfg.SchemaChannel, channel.Package, channel.Name),
				Labels: map[string]string{
					"catalog":       catalogName,
					"name":          channel.Name,
					"package":       channel.Package,
					"packageOrName": channel.Package,
					"schema":        declcfg.SchemaChannel,
				},
			},
			Spec: catalogd.CatalogMetadataSpec{
				Catalog: corev1.LocalObjectReference{Name: catalogName},
				Name:    channel.Name,
				Package: channel.Package,
				Schema:  declcfg.SchemaChannel,
				Content: content,
			},
		}

		objs = append(objs, obj)
	}

	return es.clientBuilder.WithObjects(objs...), nil
}
