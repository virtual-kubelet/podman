module github.com/virtual-kubelet/podman

go 1.12

require (
	9fans.net/go v0.0.2 // indirect
	contrib.go.opencensus.io/exporter/jaeger v0.1.0
	contrib.go.opencensus.io/exporter/ocagent v0.4.12
	github.com/BurntSushi/xgb v0.0.0-20160522181843-27f122750802
	github.com/acroca/go-symbols v0.1.1 // indirect
	github.com/buger/goterm v0.0.0-20181115115552-c206103e1f37
	github.com/chbmuc/cec v0.0.0-20170405204755-573ad0b0369b // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/davidrjenni/reftools v0.0.0-20190827201643-0605d60846fb // indirect
	github.com/docker/spdystream v0.0.0-20181023171402-6480d4af844c // indirect
	github.com/evanphx/json-patch v4.5.0+incompatible // indirect
	github.com/fatih/gomodifytags v1.0.1 // indirect
	github.com/fatih/structtag v1.1.0 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/golang/groupcache v0.0.0-20190129154638-5b532d6fd5ef // indirect
	github.com/googleapis/gnostic v0.3.0 // indirect
	github.com/gorilla/mux v1.7.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.9.2 // indirect
	github.com/haya14busa/goplay v1.0.0 // indirect
	github.com/josharian/impl v0.0.0-20190715203526-f0d59e96e372 // indirect
	github.com/mdempsky/gocode v0.0.0-20190203001940-7fb65232883f // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mjudeikis/go-podman v0.0.0-20191113175730-90d538e53252
	github.com/openshift/openshift-azure v10.1.1+incompatible
	github.com/pkg/errors v0.8.1
	github.com/ramya-rao-a/go-outline v0.0.0-20181122025142-7182a932836a // indirect
	github.com/rogpeppe/godef v1.1.1 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/skratchdot/open-golang v0.0.0-20190402232053-79abb63cd66e // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.3
	github.com/sqs/goreturns v0.0.0-20181028201513-538ac6014518 // indirect
	github.com/uudashr/gopkgs v2.0.1+incompatible // indirect
	github.com/varlink/go v0.0.0-20191018142704-4ecdbb8a36c2
	github.com/virtual-kubelet/virtual-kubelet v1.1.0
	github.com/zmb3/gogetdoc v0.0.0-20190228002656-b37376c5da6a // indirect
	go.opencensus.io v0.22.0
	go.uber.org/zap v1.12.0
	google.golang.org/api v0.6.0 // indirect
	google.golang.org/appengine v1.6.1 // indirect
	google.golang.org/genproto v0.0.0-20190620144150-6af8c5fc6601 // indirect
	google.golang.org/grpc v1.21.1 // indirect
	gotest.tools v2.2.0+incompatible
	k8s.io/api v0.0.0
	k8s.io/apimachinery v0.0.0
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/klog v0.3.3
	k8s.io/kube-openapi v0.0.0-20190603182131-db7b694dc208 // indirect
	k8s.io/kubernetes v1.15.2
	k8s.io/utils v0.0.0-20190607212802-c55fbcfc754a // indirect
	sourcegraph.com/sqs/goreturns v0.0.0-20181028201513-538ac6014518 // indirect
)

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.0.0-20190805144654-3d5bf3a310c1

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.0.0-20190805144409-8484242760e7

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20190805143448-a07e59fb081d

replace k8s.io/apiserver => k8s.io/apiserver v0.0.0-20190805142138-368b2058237c

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.0.0-20190805144531-3985229e1802

replace k8s.io/cri-api => k8s.io/cri-api v0.0.0-20190531030430-6117653b35f1

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20190805142416-fd821fbbb94e

replace k8s.io/kubelet => k8s.io/kubelet v0.0.0-20190805143852-517ff267f8d1

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.0.0-20190805144128-269742da31dd

replace k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190612205821-1799e75a0719

replace k8s.io/api => k8s.io/api v0.0.0-20190918195907-bd6ac527cfd2

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.0.0-20190805144246-c01ee70854a1

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.0.0-20190805143734-7f1675b90353

replace k8s.io/component-base => k8s.io/component-base v0.0.0-20190805141645-3a5e5ac800ae

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.0.0-20190805144012-2a1ed1f3d8a4

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190805143126-cdb999c96590

replace k8s.io/metrics => k8s.io/metrics v0.0.0-20190805143318-16b07057415d

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.0.0-20190805142637-3b65bc4bb24f

replace k8s.io/code-generator => k8s.io/code-generator v0.0.0-20190612205613-18da4a14b22b

replace k8s.io/client-go => k8s.io/client-go v0.0.0-20190805141520-2fe0317bcee0
