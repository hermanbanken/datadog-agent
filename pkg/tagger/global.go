// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package tagger

import (
	"sync"

	"github.com/DataDog/datadog-agent/cmd/agent/api/response"
	"github.com/DataDog/datadog-agent/pkg/config"
	"github.com/DataDog/datadog-agent/pkg/tagger/collectors"
	"github.com/DataDog/datadog-agent/pkg/tagger/local"
	"github.com/DataDog/datadog-agent/pkg/tagger/utils"
	"github.com/DataDog/datadog-agent/pkg/util"
	"github.com/DataDog/datadog-agent/pkg/util/containers"
	"github.com/DataDog/datadog-agent/pkg/util/containers/providers"
	"github.com/DataDog/datadog-agent/pkg/util/log"
)

// defaultTagger is the shared tagger instance backing the global Tag and Init functions
var defaultTagger Tagger
var initOnce sync.Once

// ChecksCardinality defines the cardinality of tags we should send for check metrics
// this can still be overridden when calling get_tags in python checks.
var ChecksCardinality collectors.TagCardinality

// DogstatsdCardinality defines the cardinality of tags we should send for metrics from
// dogstatsd.
var DogstatsdCardinality collectors.TagCardinality

// Init must be called once config is available, call it in your cmd
func Init() {
	initOnce.Do(func() {
		var err error
		checkCard := config.Datadog.GetString("checks_tag_cardinality")
		dsdCard := config.Datadog.GetString("dogstatsd_tag_cardinality")

		ChecksCardinality, err = collectors.StringToTagCardinality(checkCard)
		if err != nil {
			log.Warnf("failed to parse check tag cardinality, defaulting to low. Error: %s", err)
			ChecksCardinality = collectors.LowCardinality
		}
		DogstatsdCardinality, err = collectors.StringToTagCardinality(dsdCard)
		if err != nil {
			log.Warnf("failed to parse dogstatsd tag cardinality, defaulting to low. Error: %s", err)
			DogstatsdCardinality = collectors.LowCardinality
		}

		if config.IsCLCRunner() {
			log.Infof("Tagger not started on CLC")
			return
		}

		err = defaultTagger.Init()
		if err != nil {
			log.Errorf("failed to start the tagger: %s", err)
		}
	})
}

// Tag queries the defaultTagger to get entity tags from cache or sources.
// It can return tags at high cardinality (with tags about individual containers),
// or at orchestrator cardinality (pod/task level).
func Tag(entity string, cardinality collectors.TagCardinality) ([]string, error) {
	return defaultTagger.Tag(entity, cardinality)
}

// TagBuilder queries the defaultTagger to get entity tags from cache or
// sources and appends them to the TagsBuilder.  It can return tags at high
// cardinality (with tags about individual containers), or at orchestrator
// cardinality (pod/task level).
func TagBuilder(entity string, cardinality collectors.TagCardinality, tb *util.TagsBuilder) error {
	return defaultTagger.TagBuilder(entity, cardinality, tb)
}

// TagWithHash is similar to Tag but it also computes and returns the hash of the tags found
func TagWithHash(entity string, cardinality collectors.TagCardinality) ([]string, string, error) {
	tags, err := Tag(entity, cardinality)
	if err != nil {
		return tags, "", err
	}
	return tags, utils.ComputeTagsHash(tags), nil
}

// GetEntityHash returns the hash for the tags associated with the given entity
// Returns an empty string if the tags lookup fails
func GetEntityHash(entity string, cardinality collectors.TagCardinality) string {
	_, hash, _ := TagWithHash(entity, cardinality)
	return hash
}

// StandardTags queries the defaultTagger to get entity
// standard tags (env, version, service) from cache or sources.
func StandardTags(entity string) ([]string, error) {
	return defaultTagger.Standard(entity)
}

// AgentTags returns the agent tags
// It relies on the container provider utils to get the Agent container ID
func AgentTags(cardinality collectors.TagCardinality) ([]string, error) {
	ctrID, err := providers.ContainerImpl().GetAgentCID()
	if err != nil {
		return nil, err
	}

	entityID := containers.BuildTaggerEntityName(ctrID)
	return Tag(entityID, cardinality)
}

// OrchestratorScopeTag queries tags for orchestrator scope (e.g. task_arn in ECS Fargate)
func OrchestratorScopeTag() ([]string, error) {
	return defaultTagger.Tag(collectors.OrchestratorScopeEntityID, collectors.OrchestratorCardinality)
}

// OrchestratorScopeTagBuilder queries tags for orchestrator scope (e.g.
// task_arn in ECS Fargate) and appends them to the TagsBuilder
func OrchestratorScopeTagBuilder(tb *util.TagsBuilder) error {
	return defaultTagger.TagBuilder(collectors.OrchestratorScopeEntityID, collectors.OrchestratorCardinality, tb)
}

// Stop queues a stop signal to the defaultTagger
func Stop() error {
	return defaultTagger.Stop()
}

// List the content of the defaulTagger
func List(cardinality collectors.TagCardinality) response.TaggerListResponse {
	return defaultTagger.List(cardinality)
}

// SetDefaultTagger sets the global Tagger instance
func SetDefaultTagger(tagger Tagger) {
	defaultTagger = tagger
}

// GetDefaultTagger returns the global Tagger instance
func GetDefaultTagger() Tagger {
	return defaultTagger
}

func init() {
	SetDefaultTagger(local.NewTagger(collectors.DefaultCatalog))
}
