# Code generation
#
# see https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api_changes.md#generate-code

# Name of the Go package for this repository
PKG := github.com/triggermesh/eventstore

# List of API groups to generate code for
# e.g. "eventstores/v1alpha1 eventstores/v1alpha2"
API_GROUPS := eventstores/v1alpha1
# generates e.g. "PKG/pkg/apis/eventstores/v1alpha1 PKG/pkg/apis/eventstores/v1alpha2"
api-import-paths := $(foreach group,$(API_GROUPS),$(PKG)/pkg/apis/$(group))

generators := deepcopy client lister informer injection


.PHONY: codegen $(generators)

codegen: $(generators)

# http://blog.jgc.org/2007/06/escaping-comma-and-space-in-gnu-make.html
comma := ,
space :=
space +=

deepcopy:
	@echo "+ Generating deepcopy funcs for $(API_GROUPS)"
	@go run k8s.io/code-generator/cmd/deepcopy-gen \
		--go-header-file hack/boilerplate/boilerplate.go.txt \
		--input-dirs $(subst $(space),$(comma),$(api-import-paths))

client:
	@echo "+ Generating clientsets for $(API_GROUPS)"
	@rm -rf pkg/generated/client/clientset
	@go run k8s.io/code-generator/cmd/client-gen \
		--go-header-file hack/boilerplate/boilerplate.go.txt \
		--input $(subst $(space),$(comma),$(API_GROUPS)) \
		--input-base $(PKG)/pkg/apis \
		--output-package $(PKG)/pkg/generated/client/clientset

lister:
	@echo "+ Generating listers for $(API_GROUPS)"
	@rm -rf pkg/generated/client/listers
	@go run k8s.io/code-generator/cmd/lister-gen \
		--go-header-file hack/boilerplate/boilerplate.go.txt \
		--input-dirs $(subst $(space),$(comma),$(api-import-paths)) \
		--output-package $(PKG)/pkg/generated/client/listers

informer:
	@echo "+ Generating informers for $(API_GROUPS)"
	@rm -rf pkg/generated/client/informers
	@go run k8s.io/code-generator/cmd/informer-gen \
		--go-header-file hack/boilerplate/boilerplate.go.txt \
		--input-dirs $(subst $(space),$(comma),$(api-import-paths)) \
		--output-package $(PKG)/pkg/generated/client/informers \
		--versioned-clientset-package $(PKG)/pkg/generated/client/clientset/internalclientset \
		--listers-package $(PKG)/pkg/generated/client/listers

injection:
	@echo "+ Generating injection for $(API_GROUPS)"
	@rm -rf pkg/generated/client/injection
	@go run knative.dev/pkg/codegen/cmd/injection-gen \
		--go-header-file hack/boilerplate/boilerplate.go.txt \
		--input-dirs $(subst $(space),$(comma),$(api-import-paths)) \
		--output-package $(PKG)/pkg/generated/client/injection \
		--versioned-clientset-package $(PKG)/pkg/generated/client/clientset/internalclientset \
		--listers-package $(PKG)/pkg/generated/client/listers \
		--external-versions-informers-package $(PKG)/pkg/generated/client/informers/externalversions
