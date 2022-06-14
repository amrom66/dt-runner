package pkg

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"dt-runner/api"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GenerateJob is used to generate a job with given arguments
// job name is combined by ci.name and model.name
func GenerateJob(client *kubernetes.Clientset, ci api.Ci, model api.Model) (batchv1.Job, error) {

	namespace := ci.Namespace
	name := ci.Name
	repo := ci.Spec.Repo
	fmt.Println("initContainer, namespace:", namespace, "name:", name, "repo:", repo)

	if !check(ci, model) {
		return batchv1.Job{}, fmt.Errorf("ci and model are not matched")
	}
	// fix env variables by add ci variables to model variables
	variables := make(map[string]string)
	for k, v := range ci.Spec.Variables {
		variables[k] = v
	}
	for k, v := range model.Spec.Variables {
		variables[k] = v
	}
	model.Spec.Variables = variables

	job := batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      strings.Join([]string{name, model.Name}, "-"), // job name is combined by ci.name and model.name),
			Namespace: namespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: DefaultVolumeName,
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: DefaultSwapVolumeName,
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
					InitContainers: initContainers(ci.Spec.Repo, model),
					Containers:     containers(model),
					RestartPolicy:  corev1.RestartPolicyNever,
				},
			},
		},
	}
	result, err := client.BatchV1().Jobs(namespace).Create(context.TODO(), &job, metav1.CreateOptions{})
	if err != nil {
		err = fmt.Errorf("create job error %v", err)
	}
	return *result, err
}

// check is used to check ci and model are matched
func check(ci api.Ci, model api.Model) bool {
	if ci.Spec.Model != model.Name {
		return false
	}
	if ci.Namespace != model.Namespace {
		return false
	}
	re := regexp.MustCompile("(http|https):\\/\\/[\\w\\-_]+(\\.[\\w\\-_]+)+([\\w\\-\\.,@?^=%&:/~\\+#]*[\\w\\-\\@?^=%&/~\\+#])?")
	result := re.FindAllStringSubmatch(ci.Spec.Repo, -1)
	if result == nil {
		fmt.Println("repo is not matched, repo:", ci.Spec.Repo)
		return false
	}
	if len(model.Spec.Tasks) > 5 {
		fmt.Println("Number of model task should not be more than 5.")
		return false
	}
	return true
}

// container is used to generate container
func containers(model api.Model) []corev1.Container {
	containers := []corev1.Container{}

	envVars := []corev1.EnvVar{}
	for k, v := range model.Spec.Variables {
		envVars = append(envVars, corev1.EnvVar{Name: k, Value: v})
	}

	for _, task := range model.Spec.Tasks {
		containers = append(containers, corev1.Container{
			Name:       task.Name,
			Image:      task.Image,
			WorkingDir: DefaultContainerWorkspace,
			Env:        envVars,
			Args:       task.Args,
			Command:    task.Command,
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      DefaultVolumeName,
					MountPath: DefaultContainerWorkspace,
				},
				{
					Name:      DefaultSwapVolumeName,
					MountPath: DefaultSwapWorkspace,
				},
			},
		})
	}

	return containers
}

// initContainer is used to generate init container
func initContainers(repo string, model api.Model) []corev1.Container {

	envVars := []corev1.EnvVar{}
	for k, v := range model.Spec.Variables {
		envVars = append(envVars, corev1.EnvVar{Name: k, Value: v})
	}
	args := []string{
		"clone",
		"--single-branch",
		"--branch=main",
		"--",
		repo,
		DefaultInitContainerWorkspace,
	}
	initContainer := []corev1.Container{
		{
			Name:  DefaultInitContainerName,
			Image: DefaultInitContainerImage,
			Env:   envVars,
			Args:  args,
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      DefaultVolumeName,
					MountPath: DefaultInitContainerWorkspace,
				},
			},
		},
	}
	return initContainer
}
