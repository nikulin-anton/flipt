package cache

import (
	"context"
	"fmt"

	"go.flipt.io/flipt/storage"
)

const (
	evaluationRulesCachePrefix         = "eval:rules:flag:"
	evaluationDistributionsCachePrefix = "eval:dist:rule:"
)

// GetEvaluationRules returns all rules applicable to the flagKey provided from the cache if they exist; delegating to the underlying store and caching the result if no error
func (c *Store) GetEvaluationRules(ctx context.Context, flagKey string) ([]*storage.EvaluationRule, error) {
	key := evaluationRulesCachePrefix + flagKey

	// check if rules exists in cache
	data, ok, err := c.cache.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("getting rules from cache: %w", err)
	}

	if ok {
		c.logger.Debugf("cache hit: %q", key)

		rules, ok := data.([]*storage.EvaluationRule)
		if !ok {
			// not rules slice, bad cache
			c.logger.Errorf("corrupt cache, deleting: %q", key)
			if err := c.cache.Delete(ctx, key); err != nil {
				c.logger.WithError(err).Error("deleting cache entry")
			}
			goto db
		}

		return rules, nil
	}

db:
	// rules not in cache, delegate to underlying store
	rules, err := c.Store.GetEvaluationRules(ctx, flagKey)
	if err != nil {
		return rules, err
	}

	if len(rules) > 0 {
		if err := c.cache.Set(ctx, key, rules); err != nil {
			return rules, err
		}

		c.logger.Debugf("cache miss; added: %q", key)
	}

	return rules, nil
}

// GetEvaluationDistributions returns all distributions applicable to the ruleID provided from the cache if they exist; delegating to the underlying store and caching the result if no error
func (c *Store) GetEvaluationDistributions(ctx context.Context, ruleID string) ([]*storage.EvaluationDistribution, error) {
	var (
		key   = evaluationDistributionsCachePrefix + ruleID
		cache = c.cache.String()
	)

	// check if distributions exists in cache
	data, ok, err := c.cache.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("getting distributions from cache: %w", err)
	}

	if ok {
		c.logger.Debugf("cache hit: %q", key)

		distributions, ok := data.([]*storage.EvaluationDistribution)
		if !ok {
			// not distributions slice, bad cache
			c.logger.Errorf("corrupt cache, deleting: %q", key)
			if err := c.cache.Delete(ctx, key); err != nil {
				c.logger.WithError(err).Error("deleting cache entry")
			}
			goto db
		}

		return distributions, nil
	}

db:
	// distributions not in cache, delegate to underlying store
	distributions, err := c.Store.GetEvaluationDistributions(ctx, ruleID)
	if err != nil {
		return distributions, err
	}

	if len(distributions) > 0 {
		if err := c.cache.Set(ctx, key, distributions); err != nil {
			return distributions, err
		}

		c.logger.Debugf("cache miss; added %q", key)
	}

	return distributions, nil
}
