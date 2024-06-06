package k8s

import (
	"context"
	"flag"
	"path/filepath"

	"leukocyte/src/types"

	"go.uber.org/zap"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type K8s struct {
	logger *zap.Logger

	clientset *kubernetes.Clientset
}

func NewK8s(logger *zap.Logger, configPath string) *K8s {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", configPath, "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		logger.Fatal("Failed to build config", zap.Error(err))
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Fatal("Failed to create clientset", zap.Error(err))
	}

	return &K8s{
		logger:    logger,
		clientset: clientset,
	}
}

func (k *K8s) Schedule(data types.JobObject) error {
	jobsClient := k.clientset.BatchV1().Jobs(data.Namespace)
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
