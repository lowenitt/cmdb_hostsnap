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

package manager

import (
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/wusendong/cmdb_hostsnap/collector"
	"github.com/wusendong/cmdb_hostsnap/util"
)

type Manager struct {
	c collector.Collector
}

func New(c collector.Collector) *Manager {
	return &Manager{c: c}
}

func (m *Manager) Run() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT)

	stop := util.NewBool(false)
	go func() {
		for signal := range ch {
			switch signal {
			case syscall.SIGHUP:
				if err := m.c.Reload(); err != nil {
					logrus.Errorf("reload collector error: %v", err)
				}
			case syscall.SIGQUIT:
				stop.Set()
				if err := m.c.Stop(); err != nil {
					logrus.Errorf("stop collector error: %v", err)
				}
				return
			}
		}
	}()
	for {
		if stop.IsSet() {
			return
		}
		if err := m.runCollector(m.c); err != nil {
			logrus.Errorf("collector return with error: %v, we will retry 5s later", err)
		}
		time.Sleep(time.Second * 5)
	}
}

func (m *Manager) runCollector(c collector.Collector) error {
	defer func() {
		if syserr := recover(); syserr != nil {
			logrus.Errorf("collector panic: %v, we will retry 5s later, stack: \n%s\n", syserr, debug.Stack())
		}
	}()
	return c.Run()
}
