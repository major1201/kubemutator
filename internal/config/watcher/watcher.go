package watcher

import (
	"context"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	k8sv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

// ErrWatchChannelClosed indicates the watch channel close error
var ErrWatchChannelClosed = errors.New("watcher channel has closed")

// K8sConfigMapWatcher indicates the configmap watcher
type K8sConfigMapWatcher struct {
	ConfigMapLabels map[string]string
	Client          k8sv1.CoreV1Interface
}

// New creates a new K8sConfigMapWatcher
func New() (*K8sConfigMapWatcher, error) {
	k8sConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	configSet, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, err
	}

	c := new(K8sConfigMapWatcher)
	c.Client = configSet.CoreV1()
	return c, nil
}

// Watch watches the configmap change event
func (c *K8sConfigMapWatcher) Watch(ctx context.Context, notifyMe chan<- interface{}) error {
	watcher, err := c.Client.ConfigMaps("default").Watch(metav1.ListOptions{
		LabelSelector: labels.Set(c.ConfigMapLabels).String(),
	})

	if err != nil {
		return errors.WithMessage(err, "create watcher failed, RBAC disabled?")
	}

	defer watcher.Stop()

	for {
		select {
		case e, ok := <-watcher.ResultChan():
			if !ok {
				zap.L().Named("watcher").Error("watcher channel closed")
				return ErrWatchChannelClosed
			}

			if e.Type == watch.Error {
				return apierrs.FromObject(e.Object)
			}

			zap.L().Named("watcher").Info("new event", zap.String("type", string(e.Type)), zap.String("kind", e.Object.GetObjectKind().GroupVersionKind().String()))

			switch e.Type {
			case watch.Added, watch.Modified, watch.Deleted:
				notifyMe <- struct{}{}
			default:
				zap.L().Named("watcher").Error("unsupported event type", zap.String("type", string(e.Type)))
			}
		case <-ctx.Done():
			zap.L().Named("watcher").Info("stopping configmap watcher")
			return nil
		}
	}
}
