package cmd

import (
	dtapi "dt-runner/api/apps/v1"
	"dt-runner/pkg"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var ci = dtapi.Ci{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "test-job",
		Namespace: pkg.DefaultNamespace,
	},
	Spec: dtapi.CiSpec{
		Repo:   "https://github.com/linjinbao666/vrmanager.git",
		Model:  "model-sample",
		Branch: "main",
		Variables: map[string]string{
			"http_proxy":  "http://192.168.90.110:1087",
			"https_proxy": "http://192.168.90.110:1087",
			"no_proxy":    "localhost",
			"MAVEN_OPTS":  "-DproxySet=true -DproxyHost=192.168.90.110 -DproxyPort=1087",
		},
	},
}
var model = dtapi.Model{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "model-sample",
		Namespace: pkg.DefaultNamespace,
	},
	Spec: dtapi.ModelSpec{
		Tasks: []dtapi.Task{
			{
				Name:    "build",
				Image:   "docker.io/linjinbao66/dt-maven:0.0.4",
				Command: []string{"/bin/sh", "-c"},
				Args:    []string{"mvn clean package && echo 'success' > /dtswap/build.log"},
			},
			{
				Name:    "archive",
				Image:   "docker.io/linjinbao66/dt-mc:0.0.1",
				Command: []string{"/bin/sh", "-c"},
				Args:    []string{"while true; do sleep 30 && [ -f '/dtswap/build.log' ] && echo 'success' > /dtswap/archive.log && exit 0 || echo 'file not exists'; done"},
			},
			{
				Name:    "clean",
				Image:   "docker.io/busybox:latest",
				Command: []string{"/bin/sh", "-c"},
				Args:    []string{"while true; do sleep 30 && [ -f '/dtswap/archive.log' ] && exit 0 || echo 'file not exists'; done"},
			},
		},
	},
}

func Test_ExecuteCommand(t *testing.T) {
	cmd := serverCmd

	cmd.Execute()
}
