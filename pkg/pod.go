package pkg

import (
	"fmt"
	"regexp"
	"strings"

	appsv1 "dt-runner/api/apps/v1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

// DtJob 是对ci和model的封装
type DtJob struct {
	name        string
	ci          string //关联的ci名称
	project     string
	httpurl     string
	sshurl      string
	branch      string
	ref         string
	checkoutSHA string
}

// GeneratePod is used to generate a pod with given arguments
// GeneratePod will not create pod using kubernetes client, it just generate the pod spec
// Pod name is combined by ci.name and model.name
func GeneratePod(dtJob DtJob) (corev1.Pod, error) {

	ci := GetCi(DefaultNamespace, dtJob.ci)
	model := GetModel(DefaultNamespace, ci.Spec.Model)

	if !check(*ci, *model) {
		return corev1.Pod{}, fmt.Errorf("ci and model are not matched")
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

	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      strings.Join([]string{dtJob.name, ci.Name, model.Name}, "-"), // pod name is combined by dtjob.name ci.name and model.name),
			Namespace: ci.Namespace,
			Labels: map[string]string{
				DefaultLabelDtRunner:      "true",
				DefaultLabelDtRunnerCi:    ci.Name,
				DefaultLabelDtRunnerModel: model.Name,
			},
		},
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
			InitContainers: initContainer(ci.Spec.Repo, *model),
			Containers:     containers(*model),
			RestartPolicy:  corev1.RestartPolicyNever,
		},
	}
	return pod, nil
}

// check is used to check ci and model are matched
// TODO check name and namespace are matched with kubernetes rules
func check(ci appsv1.Ci, model appsv1.Model) bool {
	if ci.Spec.Model != model.Name {
		return false
	}
	if ci.Namespace != model.Namespace {
		return false
	}
	klog.Info("repo name:", ci.Spec.Repo)
	re := regexp.MustCompile("^http|https://github.com|gitlab.com|dtwave-inc.com/*")
	result := re.FindAllStringSubmatch(ci.Spec.Repo, -1)
	if result == nil {
		klog.Info("repo is not matched, repo:", ci.Spec.Repo)
		return false
	}
	if len(model.Spec.Tasks) > 5 {
		klog.Info("Number of model task should not be more than 5.")
		return false
	}
	return true
}

// container is used to generate container
func containers(model appsv1.Model) []corev1.Container {
	containers := []corev1.Container{}

	envVars := []corev1.EnvVar{}
	for k, v := range model.Spec.Variables {
		envVars = append(envVars, corev1.EnvVar{Name: k, Value: v})
	}

	for _, task := range model.Spec.Tasks {
		klog.Info("container name: ", task.Name)
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
func initContainer(repo string, model appsv1.Model) []corev1.Container {

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
