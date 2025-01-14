package main

import (
	"os"

	corev1 "k8s.io/api/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	"context"
	"fmt"
	"strings"

	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"k8s.io/apimachinery/pkg/runtime"
)

func init() {
	log.SetLogger(zap.New())
}

func main() {
	entryLog := log.Log.WithName("entrypoint")

	// Setup a Manager
	entryLog.Info("setting up manager")
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		WebhookServer: webhook.NewServer(webhook.Options{
			CertDir: "/tmp/k8s-webhook-server/serving-certs",
		}),
	})
	if err != nil {
		entryLog.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	entryLog.Info("Setting up controller")

	if err := builder.WebhookManagedBy(mgr).
		For(&corev1.Pod{}).
		WithDefaulter(NewImageTransformer()).
		Complete(); err != nil {
		entryLog.Error(err, "unable to create webhook", "webhook", "Pod")
		os.Exit(1)
	}

	entryLog.Info("starting manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "unable to run manager")
		os.Exit(1)
	}
}

// +kubebuilder:webhook:path=/mutate-v1-pod,mutating=true,failurePolicy=fail,groups="",resources=pods,verbs=create;update,versions=v1,image-transform-webhook
type imageTransformer struct {
	NewRepo       string
	OriginalRepos []string
}

func NewImageTransformer() *imageTransformer {
	originalRepos := []string{"docker.io", "gcr.io", "ghcr.io", "registry.k8s.io"}
	newRepo := "m.daocloud.io"

	originalRepo := os.Getenv("ORIGINAL_REPO")
	newRepoEnv := os.Getenv("NEW_REPO")
	if originalRepo != "" {
		originalRepos = strings.Split(originalRepo, ",")
	}
	if newRepoEnv != "" {
		newRepo = newRepoEnv
	}

	return &imageTransformer{
		OriginalRepos: originalRepos,
		NewRepo:       newRepo,
	}
}

func (a *imageTransformer) Default(ctx context.Context, obj runtime.Object) error {
	logf := log.FromContext(ctx)
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return fmt.Errorf("expected a Pod but got a %T", obj)
	}

	var initChanged, tempChanged bool

	// change init image
	pod.Spec.InitContainers, initChanged = a.modifyContainers(pod.Spec.InitContainers)

	// change template image
	pod.Spec.Containers, tempChanged = a.modifyContainers(pod.Spec.Containers)

	if (initChanged || tempChanged) && os.Getenv("WITH_SECRET") == "True" {
		pod.Spec.ImagePullSecrets = []corev1.LocalObjectReference{
			{
				Name: os.Getenv("SECRET_NAME"),
			},
		}
	}

	logf.Info("image transformed")
	return nil
}

// modifyContainers modifies the containers
func (a *imageTransformer) modifyContainers(containers []corev1.Container) ([]corev1.Container, bool) {
	changed := false
	for i, container := range containers {
		for _, originalRepo := range a.OriginalRepos {
			newImage, changed := needChangeImage(container.Image, a.NewRepo, originalRepo)
			if changed {
				containers[i].Image = newImage
			}
		}
	}
	return containers, changed
}

// needChangeImage checks if the image needs to be changed
func needChangeImage(image, newRepo, originalRepo string) (string, bool) {
	newImage := ""

	if strings.HasPrefix(image, originalRepo) || strings.Count(image, "/") == 1 || !strings.Contains(image, "/") {
		// docker.io/busybox:latest -> aaa.bbb.ccc/docker.io/busybox:latest
		// grafana/loki:2.0.0 -> aaa.bbb.ccc/grafana/loki:2.0.0
		// busybox:latest -> aaa.bbb.ccc/busybox:latest
		newImage = fmt.Sprintf("%s/%s", newRepo, image)
	}
	return newImage, newImage != ""
}
