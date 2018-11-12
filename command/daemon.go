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
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"

	"github.com/wusendong/cmdb_hostsnap/collector"
	"github.com/wusendong/cmdb_hostsnap/manager"
	"github.com/wusendong/cmdb_hostsnap/pidfile"
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
		Action: DaemonAction,
	}
}

func DaemonAction(c *cli.Context) {
	if err := pidfile.SavePid(); err != nil {
		logrus.Fatalf("Error saving pid file %v", err)
	}

	configfile := c.String("config")
	hostsnap, err := collector.NewHostsnap(configfile)
	if err != nil {
		logrus.Fatalf("Error NewHostsnap %v", err)
	}

	man := manager.New(hostsnap)
	man.Run()
}
