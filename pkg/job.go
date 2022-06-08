package pkg

import (
	"context"
	"fmt"
	"log"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// generateJob is used to generate a job with given arguments
func generateJob(client *kubernetes.Clientset, name string, namespace string, repo string) (batchv1.Job, error) {
	fmt.Println("generateJob")
	job := batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: "data",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
					InitContainers: []corev1.Container{
						{
							Name:  "prepare",
							Image: "alpine/git",
							Env: []corev1.EnvVar{
								{
									Name:  "repo",
									Value: "",
								},
								{
									Name:  "http_proxy",
									Value: "http://192.168.90.110:1087",
								},
								{
									Name:  "https_proxy",
									Value: "http://192.168.90.110:1087",
								},
							},
							Args: []string{
								"clone",
								"--branch=main",
								"--depth=1",
								"--recursive",
								"--single-branch",
								"--progress",
								"--",
								repo,
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "data",
									MountPath: "/opt/workspace",
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:       "build",
							Image:      "maven:latest",
							WorkingDir: "/opt/workspace",
							Env: []corev1.EnvVar{
								{
									Name:  "repo",
									Value: repo,
								},
								{
									Name:  "http_proxy",
									Value: "http://192.168.90.110:1087",
								},
								{
									Name:  "https_proxy",
									Value: "http://192.168.90.110:1087",
								},
								{
									Name:  "PROXY_HOST",
									Value: "192.168.90.110",
								},
								{
									Name:  "PROXY_PORT",
									Value: "1087",
								},
								{
									Name:  "MAVEN_OPTS",
									Value: "-DproxySet=true -DproxyHost=192.168.90.110 -DproxyPort=1087",
								},
							},
							Args: []string{
								"mvn",
								"clean",
								"package",
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "data",
									MountPath: "/opt/workspace",
								},
							},
						},
					},
					RestartPolicy: "Never",
				},
			},
		},
	}
	result, err := client.BatchV1().Jobs(namespace).Create(context.TODO(), &job, metav1.CreateOptions{})
	if err != nil {
		log.Println("create job error ", err)
	}
	return *result, err
}
