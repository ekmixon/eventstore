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

generators := deepcopy


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
