// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package appsignals

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"

	"github.com/aws/amazon-cloudwatch-agent/internal/util/collections"
	"github.com/aws/amazon-cloudwatch-agent/translator/translate/otel/common"
)

func TestTranslatorTraces(t *testing.T) {
	type want struct {
		receivers  []string
		processors []string
		exporters  []string
		extensions []string
	}
	tt := NewTranslator(component.DataTypeTraces)
	assert.EqualValues(t, "traces/app_signals", tt.ID().String())
	testCases := map[string]struct {
		input      map[string]interface{}
		want       *want
		wantErr    error
		detector   func() (common.Detector, error)
		isEKSCache func() common.IsEKSCache
	}{
		"WithoutTracesCollectedKey": {
			input:   map[string]interface{}{},
			wantErr: &common.MissingKeyError{ID: tt.ID(), JsonKey: fmt.Sprint(common.AppSignalsTraces)},
		},
		"WithAppSignalsEnabledTracesEKS": {
			input: map[string]interface{}{
				"traces": map[string]interface{}{
					"traces_collected": map[string]interface{}{
						"app_signals": map[string]interface{}{},
					},
				},
			},
			want: &want{
				receivers:  []string{"otlp/app_signals"},
				processors: []string{"resourcedetection", "awsappsignals"},
				exporters:  []string{"awsxray/app_signals"},
				extensions: []string{"awsproxy/app_signals", "agenthealth/traces"},
			},
			detector:   common.TestEKSDetector,
			isEKSCache: common.TestIsEKSCacheEKS,
		},
		"WithAppSignalsEnabledK8s": {
			input: map[string]interface{}{
				"traces": map[string]interface{}{
					"traces_collected": map[string]interface{}{
						"app_signals": map[string]interface{}{},
					},
				},
			},
			want: &want{
				receivers:  []string{"otlp/app_signals"},
				processors: []string{"awsappsignals"},
				exporters:  []string{"awsxray/app_signals"},
				extensions: []string{"awsproxy/app_signals", "agenthealth/traces"},
			},
			detector:   common.TestK8sDetector,
			isEKSCache: common.TestIsEKSCacheK8s,
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv(common.KubernetesEnvVar, "TEST")
			common.NewDetector = testCase.detector
			common.IsEKS = testCase.isEKSCache
			conf := confmap.NewFromStringMap(testCase.input)
			got, err := tt.Translate(conf)
			assert.Equal(t, testCase.wantErr, err)
			if testCase.want == nil {
				assert.Nil(t, got)
			} else {
				require.NotNil(t, got)
				assert.Equal(t, testCase.want.receivers, collections.MapSlice(got.Receivers.Keys(), component.ID.String))
				assert.Equal(t, testCase.want.processors, collections.MapSlice(got.Processors.Keys(), component.ID.String))
				assert.Equal(t, testCase.want.exporters, collections.MapSlice(got.Exporters.Keys(), component.ID.String))
				assert.Equal(t, testCase.want.extensions, collections.MapSlice(got.Extensions.Keys(), component.ID.String))
			}
		})
	}
}

func TestTranslatorMetrics(t *testing.T) {
	type want struct {
		receivers  []string
		processors []string
		exporters  []string
		extensions []string
	}
	tt := NewTranslator(component.DataTypeMetrics)
	assert.EqualValues(t, "metrics/app_signals", tt.ID().String())
	testCases := map[string]struct {
		input      map[string]interface{}
		want       *want
		wantErr    error
		detector   func() (common.Detector, error)
		isEKSCache func() common.IsEKSCache
	}{
		"WithoutMetricsCollectedKey": {
			input:   map[string]interface{}{},
			wantErr: &common.MissingKeyError{ID: tt.ID(), JsonKey: fmt.Sprint(common.AppSignalsMetrics)},
		},
		"WithAppSignalsEnabledMetrics": {
			input: map[string]interface{}{
				"logs": map[string]interface{}{
					"metrics_collected": map[string]interface{}{
						"app_signals": map[string]interface{}{},
					},
				},
			},
			want: &want{
				receivers:  []string{"otlp/app_signals"},
				processors: []string{"resourcedetection", "awsappsignals"},
				exporters:  []string{"awsemf/app_signals"},
				extensions: []string{"agenthealth/logs"},
			},
			detector:   common.TestEKSDetector,
			isEKSCache: common.TestIsEKSCacheEKS,
		},
		"WithAppSignalsEnabledK8s": {
			input: map[string]interface{}{
				"logs": map[string]interface{}{
					"metrics_collected": map[string]interface{}{
						"app_signals": map[string]interface{}{},
					},
				},
			},
			want: &want{
				receivers:  []string{"otlp/app_signals"},
				processors: []string{"awsappsignals"},
				exporters:  []string{"awsemf/app_signals"},
				extensions: []string{"agenthealth/logs"},
			},
			detector:   common.TestK8sDetector,
			isEKSCache: common.TestIsEKSCacheK8s,
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv(common.KubernetesEnvVar, "TEST")
			common.NewDetector = testCase.detector
			common.IsEKS = testCase.isEKSCache
			conf := confmap.NewFromStringMap(testCase.input)
			got, err := tt.Translate(conf)
			assert.Equal(t, testCase.wantErr, err)
			if testCase.want == nil {
				assert.Nil(t, got)
			} else {
				require.NotNil(t, got)
				assert.Equal(t, testCase.want.receivers, collections.MapSlice(got.Receivers.Keys(), component.ID.String))
				assert.Equal(t, testCase.want.processors, collections.MapSlice(got.Processors.Keys(), component.ID.String))
				assert.Equal(t, testCase.want.exporters, collections.MapSlice(got.Exporters.Keys(), component.ID.String))
				assert.Equal(t, testCase.want.extensions, collections.MapSlice(got.Extensions.Keys(), component.ID.String))
			}
		})
	}
}
