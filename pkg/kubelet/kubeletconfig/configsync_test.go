package kubeletconfig

import (
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	api "k8s.io/kubernetes/pkg/apis/core"
	"reflect"
	"testing"
)

func Test_latestNodeConfigSource(t *testing.T) {

	indexer := cache.NewStore(cache.MetaNamespaceKeyFunc)

	configResource := &apiv1.NodeConfigSource{
		ConfigMap: &apiv1.ConfigMapNodeConfigSource{
			Name:             "c-name",
			Namespace:        "c-namespace",
			UID:              "c-uid",
			KubeletConfigKey: "c-key",
		},
	}

	node := &apiv1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node",
			UID:  "node-uid",
		},
		Spec: apiv1.NodeSpec{
			ConfigSource: configResource,
		},
	}
	_ = indexer.Add(node)

	nodeApi := &api.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nodeApi",
			UID:  "nodeApi-uid",
		},
	}
	_ = indexer.Add(nodeApi)

	type args struct {
		store    cache.Store
		nodeName string
	}
	tests := []struct {
		name    string
		args    args
		want    *apiv1.NodeConfigSource
		wantErr bool
	}{
		{
			name: "node name exists in cache",
			args: args{
				store:    indexer,
				nodeName: "node",
			},
			want:    configResource,
			wantErr: false,
		},
		{
			name: "node name doesn't exist in cache",
			args: args{
				store:    indexer,
				nodeName: "node-none",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "node name exists in cache,but failed to cast object from informer's store to Node",
			args: args{
				store:    indexer,
				nodeName: "nodeApi",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := latestNodeConfigSource(tt.args.store, tt.args.nodeName)
			if (err != nil) != tt.wantErr {
				t.Errorf("latestNodeConfigSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("latestNodeConfigSource() got = %v, want %v", got, tt.want)
			}
		})
	}
}

