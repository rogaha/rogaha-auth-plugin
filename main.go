package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/client"
	authz "github.com/docker/go-plugins-helpers/authorization"
	"os/user"
	"strconv"
)

const (
	socketAddress = "/run/docker/plugins/rogaha.sock"
)

type authzPlugin struct {
	authz.Plugin
}

func (p *authzPlugin) AuthZReq(r authz.Request) authz.Response {
	logrus.Infof("%s", r.RequestURI)
	if r.ResponseBody != nil {
		n := bytes.IndexByte(r.RequestBody, 0)
		s := string(r.ResponseBody[:n])
		logrus.Infof("%s", s)
	}
	if r.RequestURI == "/_ping" {
		cli, err := client.NewEnvClient()
		if err != nil {
			logrus.Error(err)
		} else {
			info, _ := cli.Info(context.Background())
			logrus.Infof("%+v", info)
		}
	}

	return authz.Response{
		Allow: true,
		Msg:   fmt.Sprintf("You are authorized %+v\n", r),
		Err:   "",
	}
}

func (p *authzPlugin) AuthZRes(r authz.Request) authz.Response {
	return authz.Response{
		Allow: true,
		Msg:   fmt.Sprintf("You are authorized %+v\n", r),
		Err:   "",
	}
}

func main() {
	d := &authzPlugin{}
	handler := authz.NewHandler(d)
	u, _ := user.Lookup("root")
	gid, _ := strconv.Atoi(u.Gid)
	logrus.Infof("listening on %s", socketAddress)
	logrus.Error(handler.ServeUnix(socketAddress, gid))
}
