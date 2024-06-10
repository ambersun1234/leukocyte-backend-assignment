package k8s

import (
	"context"

	"leukocyte/src/types"

	"go.uber.org/zap"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type K8s struct {
	logger *zap.Logger

	clientset *kubernetes.Clientset
}

func NewK8s(logger *zap.Logger) *K8s {
	clientset, err := kubernetes.NewForConfig(readK8sConfig(logger))
	if err != nil {
		logger.Fatal("Failed to create clientset", zap.Error(err))
	}

	return &K8s{
		logger:    logger,
		clientset: clientset,
	}
}

func readK8sConfig(logger *zap.Logger) *rest.Config {
	config, err := rest.InClusterConfig()
	if err != nil {
		logger.Fatal("Failed to build config", zap.Error(err))
	}
	return config
}

func (k *K8s) Schedule(data types.JobObject) error {
	jobsClient := k.clientset.BatchV1().Jobs("default")
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: data.Name,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "schedule",
							Image: data.Image,
							Args:  data.Commands,
						},
					},
					RestartPolicy: corev1.RestartPolicy(data.RestartPolicy),
				},
			},
		},
	}

	result, err := jobsClient.Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		k.logger.Error("Failed to create job", zap.Error(err))

		return err
	}

	k.logger.Info("Created job", zap.String("job", string(result.GetObjectMeta().GetUID())))

	return nil
}
