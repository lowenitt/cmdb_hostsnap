/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package publiser

import (
	"encoding/json"
	"fmt"

	redispkg "gopkg.in/redis.v5"

	"github.com/wusendong/cmdb_hostsnap/storage/redis"
)

type RedisPubliser struct {
	*redispkg.Client
	channel string
}

func NewRedisPubliser(channel string, conf redis.Config) (*RedisPubliser, error) {
	cli, err := redis.NewFromConfig(conf)
	if err != nil {
		return nil, err
	}
	return &RedisPubliser{Client: cli, channel: channel}, nil
}

func (p *RedisPubliser) PublishEvent(event map[string]interface{}) error {
	message, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal json failed. error: %v", err)
	}
	err = p.Client.Publish(p.channel, string(message)).Err()
	if err != nil {
		return fmt.Errorf("publish to redis failed. error: %v", err)
	}
	return nil
}

func (p *RedisPubliser) PublishEvents(events []map[string]interface{}) error {
	for _, event := range events {
		message, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("marshal json failed. error: %v", err)
		}
		err = p.Client.Publish(p.channel, string(message)).Err()
		if err != nil {
			return fmt.Errorf("publish to redis failed. error: %v", err)
		}
	}
	return nil
}

func (p *RedisPubliser) Close() error {
	return p.Client.Close()
}
