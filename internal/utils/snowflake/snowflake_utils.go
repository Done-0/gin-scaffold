// Package snowflake provides snowflake algorithm ID generation utilities
// Author: Done-0
// Created: 2025-09-25
package snowflake

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/snowflake"
)

var (
	node *snowflake.Node
	once sync.Once
)

// GenerateID generates snowflake algorithm ID
func GenerateID() (int64, error) {
	once.Do(func() {
		var err error
		node, err = snowflake.NewNode(1)
		if err != nil {
			fmt.Printf("failed to initialize snowflake node: %v", err)
		}
	})

	switch {
	case node != nil:
		return node.Generate().Int64(), nil

	default:
		// Snowflake format: 41-bit timestamp + 10-bit node ID + 12-bit sequence number
		// Standard snowflake epoch, node ID 1, sequence number uses lower 12 bits of current nanoseconds
		ts := time.Now().UnixMilli() - 1288834974657
		nodeID := int64(1)
		seq := time.Now().UnixNano() & 0xFFF

		return (ts << 22) | (nodeID << 12) | seq, nil
	}
}
