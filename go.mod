module github.com/openyurtio/kole

go 1.16

require (
	github.com/aliyunmq/mq-http-go-sdk v1.0.3
	github.com/eclipse/paho.golang v0.10.0
	github.com/eclipse/paho.mqtt.golang v1.3.5
	github.com/evanphx/json-patch v4.11.0+incompatible // indirect
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/frankban/quicktest v1.14.0 // indirect
	github.com/go-logr/logr v0.4.0 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.5 // indirect
	github.com/go-openapi/swag v0.19.14 // indirect
	github.com/gogap/errors v0.0.0-20210818113853-edfbba0ddea9
	github.com/gogap/stack v0.0.0-20150131034635-fef68dddd4f8 // indirect
	github.com/golang/snappy v0.0.3
	github.com/google/uuid v1.3.0
	github.com/googleapis/gnostic v0.5.5 // indirect
	github.com/imdario/mergo v0.3.9 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/onsi/gomega v1.14.0 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.10.1
	github.com/valyala/fasthttp v1.31.0 // indirect
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba // indirect
	k8s.io/api v0.22.2
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v1.5.2
	k8s.io/code-generator v0.22.2
	k8s.io/klog/v2 v2.30.0
	k8s.io/kube-openapi v0.0.0-20210421082810-95288971da7e // indirect
	k8s.io/utils v0.0.0-20210819203725-bdf08cb9a70a // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.1.2 // indirect
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/eclipse/paho.golang => github.com/eclipse/paho.golang v0.10.0
	github.com/eclipse/paho.mqtt.golang => github.com/eclipse/paho.mqtt.golang v1.3.5
	github.com/google/uuid => github.com/google/uuid v1.3.0
	github.com/mitchellh/go-homedir => github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra => github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag => github.com/spf13/pflag v1.0.5
	github.com/spf13/viper => github.com/spf13/viper v1.10.1
	k8s.io/api => k8s.io/api v0.20.4
	k8s.io/apimachinery => k8s.io/apimachinery v0.20.4
	k8s.io/client-go => k8s.io/client-go v0.20.4
	k8s.io/code-generator => k8s.io/code-generator v0.20.4
	k8s.io/klog/v2 => k8s.io/klog/v2 v2.4.0
	sigs.k8s.io/controller-tools => sigs.k8s.io/controller-tools v0.7.0
	sigs.k8s.io/yaml => sigs.k8s.io/yaml v1.2.0
)
