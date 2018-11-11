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

package command

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/wusendong/cmdb_hostsnap/collector"
	"github.com/wusendong/cmdb_hostsnap/pidfile"
	"gopkg.in/urfave/cli.v1"
)

func DaemonCmd() cli.Command {
	return cli.Command{
		Name:  "daemon",
		Usage: "start the hostsnapshot daemon",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "config, c",
				Usage: "Load configuration form `FILE`",
			},
		},
		Action: func(c *cli.Context) {
			if err := pidfile.SavePid(); err != nil {
				logrus.Fatalf("Error saving pid file %v", err)
			}

			configfile := c.String("config")
			hostsnap, err := collector.NewHostsnap(configfile)
			if err != nil {
				logrus.Fatalf("Error NewHostsnap %v", err)
			}

			ch := make(chan os.Signal, 1)
			signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT)

			go func() {
				for signal := range ch {
					switch signal {
					case syscall.SIGHUP:
						hostsnap.Reload()
					case syscall.SIGQUIT:
						hostsnap.Stop()
					}
				}
			}()

			if err := hostsnap.Run(); err != nil {
				logrus.Fatalf("Error starting collector: %v", err)
			}
		},
	}
}
