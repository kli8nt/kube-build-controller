package utils

import (
	"context"
	"log"

	"github.com/Mohamed-Rafraf/kube-builder/config"
	"github.com/Mohamed-Rafraf/kube-builder/pkg"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Build(data *pkg.Data) error {
	// Create the Pod object
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: data.ApplicationName,
		},
		Spec: corev1.PodSpec{
			InitContainers: []corev1.Container{
				{
					Name:  "builder",
					Image: "mohamedrafraf/builder:pfa",
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "workdir",
							MountPath: "/workdir",
						},
					},
					Env: []corev1.EnvVar{
						{
							Name:  "TECHNOLOGY",
							Value: data.Technology,
						},
						{
							Name:  "Version",
							Value: data.Version,
						},
						{
							Name:  "REPO_URL",
							Value: data.RepositoryURL,
						},
						{
							Name:  "APP_NAME",
							Value: data.ApplicationName,
						},
						{
							Name:  "INSTALL_CMD",
							Value: data.InstallCommand,
						},
						{
							Name:  "RUN_CMD",
							Value: data.RunCommand,
						},
						{
							Name:  "STATIC",
							Value: data.IsStatic,
						},
						{
							Name:  "PORT",
							Value: data.Port,
						},
						{
							Name:  "ENVARS",
							Value: data.EnvironmentVariables,
						},
						{
							Name:  "DEPENDENCIES",
							Value: data.DependenciesFiles,
						},
					},
				},
			},
			Containers: []corev1.Container{
				{
					Name:  "kaniko",
					Image: "gcr.io/kaniko-project/executor:latest",
					Args: []string{
						"--dockerfile=/workdir/" + data.ApplicationName + "/Dockerfile",
						"--context=dir:///workdir/" + data.ApplicationName,
						"--destination=us-central1-docker.pkg.dev/esoteric-might-387308/kli8nt/" + data.ApplicationName,
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "kaniko-secret",
							MountPath: "/secret",
						},
						{
							Name:      "workdir",
							MountPath: "/workdir",
						},
					},
					Env: []corev1.EnvVar{
						{
							Name:  "GOOGLE_APPLICATION_CREDENTIALS",
							Value: "/secret/kaniko-secret.json",
						},
					},
				},
			},
			RestartPolicy: corev1.RestartPolicyNever,
			Volumes: []corev1.Volume{
				{
					Name: "workdir",
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{},
					},
				},
				{
					Name: "kaniko-secret",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "kaniko-secret",
						},
					},
				},
			},
		},
	}

	_, err := config.Clientset.CoreV1().Pods("default").Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	log.Println("Building The Application named ", data.ApplicationName)

	return nil

}

func Delete() (*pkg.Status, error) {
	var status *pkg.Status
	// Get the list of pods in the default namespace
	pods, err := config.Clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Find the first pod with the "Completed", "Error", or "Failed" status
	var podToDelete *corev1.Pod
	for _, pod := range pods.Items {
		if pod.Status.Phase == corev1.PodSucceeded ||
			pod.Status.Phase == corev1.PodFailed ||
			pod.Status.Phase == corev1.PodUnknown {
			podToDelete = &pod
			break
		}
	}

	// If a pod is found, delete it
	if podToDelete != nil {

		for _, env := range podToDelete.Spec.InitContainers[0].Env {
			if env.Name == "PORT" {
				status.Port = env.Value
			}
		}
		if podToDelete.Status.Phase == corev1.PodSucceeded {
			status.Status = "Succeeded"
		} else {
			status.Status = "Failed"
		}

		err = config.Clientset.CoreV1().Pods("default").Delete(context.TODO(), podToDelete.Name, metav1.DeleteOptions{})
		if err != nil {
			return nil, err
		}
		status.ApplicationName = podToDelete.Name

		log.Printf("Deleted pod: %s\n", podToDelete.Name)

		return status, nil
	}

	return nil, nil

}
