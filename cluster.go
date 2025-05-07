package main

import (
    "context"
    "fmt"
    "github.com/go-redis/redis/v8"
)

// Redis context
var ctx = context.Background()

func fetchZSetsFromShards(redisClient *redis.ClusterClient, keys []string) (map[string][]string, error) {
    pipeline := redisClient.Pipeline()
    cmds := make(map[string]*redis.StringSliceCmd)

    // Queue pipeline requests
    for _, key := range keys {
        cmds[key] = pipeline.ZRange(ctx, key, 0, -1) // Fetch all sorted set members
    }

    // Execute pipeline
    _, err := pipeline.Exec(ctx)
    if err != nil {
        return nil, err
    }

    // Process results
    results := make(map[string][]string)
    for key, cmd := range cmds {
        values, err := cmd.Result()
        if err != nil {
            return nil, err
        }
        results[key] = values
    }

    return results, nil
}

func main() {
    // Connect to Redis cluster
    redisClient := redis.NewClusterClient(&redis.ClusterOptions{
        Addrs: []string{"shard1:6379", "shard2:6379", "shard3:6379"}, // Replace with actual shard addresses
    })

    keys := []string{"key1", "key2", "key3"}
    data, err := fetchZSetsFromShards(redisClient, keys)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    // Display results
    for k, v := range data {
        fmt.Printf("Key: %s -> Values: %v\n", k, v)
    }
}
