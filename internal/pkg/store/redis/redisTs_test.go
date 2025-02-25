// Copyright 2021-2022 EMQ Technologies Co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build redisdb || !core

package redis

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/lf-edge/ekuiper/internal/pkg/store/test/common"
	ts2 "github.com/lf-edge/ekuiper/pkg/kv"
	"github.com/redis/go-redis/v9"
	"testing"
)

func TestRedisTsSet(t *testing.T) {
	ks, db, minRedis := setupTRedisKv()
	defer cleanRedisKv(db, minRedis)

	common.TestTsSet(ks, t)
}

func TestRedisTsLast(t *testing.T) {
	ks, db, minRedis := setupTRedisKv()
	defer cleanRedisKv(db, minRedis)

	common.TestTsLast(ks, t)
}

func TestRedisTsGet(t *testing.T) {
	ks, db, minRedis := setupTRedisKv()
	defer cleanRedisKv(db, minRedis)

	common.TestTsGet(ks, t)
}

func TestRedisTsDelete(t *testing.T) {
	ks, db, minRedis := setupTRedisKv()
	defer cleanRedisKv(db, minRedis)

	common.TestTsDelete(ks, t)
}

func TestRedisTsDeleteBefore(t *testing.T) {
	ks, db, minRedis := setupTRedisKv()
	defer cleanRedisKv(db, minRedis)

	common.TestTsDeleteBefore(ks, t)
}

func setupTRedisKv() (ts2.Tskv, *redis.Client, *miniredis.Miniredis) {
	minRedis, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	redisDB := redis.NewClient(&redis.Options{
		Addr: minRedis.Addr(),
	})

	builder := NewTsBuilder(redisDB)
	var ks ts2.Tskv
	ks, err = builder.CreateTs("test")
	if err != nil {
		panic(err)
	}
	return ks, redisDB, minRedis
}
