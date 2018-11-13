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

package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"

	"github.com/wusendong/cmdb_hostsnap/command"
)

// app info
const (
	AppName = "cmdb_hostsnap"
	Usage   = "take snapshot for host information and report it via redis channel "
)

// build info
var (
	Version     = "0.1.0"
	BuildCommit = ""
	BuildTime   = ""
	GoVersion   = ""
)

func cmdNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "Unrecognized command: %s", command)
}

func onUsageError(c *cli.Context, err error, isSubcommand bool) error {
	err = fmt.Errorf("Usage error, please check your command: %s", err)
	fmt.Fprintf(os.Stderr, err.Error())
	return err
}

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	a := cli.NewApp()
	a.Version = Version
	a.Name = Usage

	a.Description = fmt.Sprintf(`
	BuildCommit : %s
	BuildTime   : %s
	GoVersion   : %s`, BuildCommit, BuildTime, GoVersion)

	a.Before = func(c *cli.Context) error {
		if c.GlobalBool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}
		return nil
	}

	a.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug, d",
			Usage:  "enable debug logging level",
			EnvVar: "CMDB_DEBUG",
		},
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "Load configuration from `FILE`",
			EnvVar: "CMDB_HOSTSNAP_CONFIG",
		},
	}
	a.Commands = []cli.Command{
		command.ReloadCmd(),
		command.StopCmd(),
	}
	a.CommandNotFound = cmdNotFound
	a.OnUsageError = onUsageError

	a.Action = command.DaemonAction

	if err := a.Run(os.Args); err != nil {
		message := fmt.Sprintf("Critical error: %v", err)
		// log to std error
		fmt.Fprintln(os.Stderr, message)
		// log to log file
		logrus.Fatal(message)
	}
}
