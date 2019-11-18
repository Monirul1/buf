# Managed by makego. DO NOT EDIT.

# Must be set
$(call _assert_var,MAKEGO)
$(call _conditional_include,$(MAKEGO)/base.mk)
$(call _conditional_include,$(MAKEGO)/dep_protoc.mk)
$(call _conditional_include,$(MAKEGO)/dep_protoc_gen_validate.mk)
$(call _conditional_include,$(MAKEGO)/protoc_gen_go.mk)
# Must be set
$(call _assert_var,PROTO_PATH)
# Must be set
$(call _assert_var,PROTOC_GEN_VALIDATE_OUT)
$(call _assert_var,CACHE_INCLUDE)
$(call _assert_var,PROTOC)
$(call _assert_var,PROTOC_GEN_VALIDATE)

# Not modifiable for now
PROTOC_GEN_VALIDATE_OPT := lang=go
PROTO_INCLUDE_PATHS := $(PROTO_INCLUDE_PATHS) third_party/proto

EXTRA_MAKEGO_FILES := $(EXTRA_MAKEGO_FILES) scripts/protoc_gen_plugin.bash

.PHONY: protocgenvalidate
protocgenvalidate: protocgengoclean $(PROTOC) $(PROTOC_GEN_VALIDATE)
	bash $(MAKEGO)/scripts/protoc_gen_plugin.bash \
		"--proto_path=$(PROTO_PATH)" \
		"--proto_include_path=$(CACHE_INCLUDE)" \
		$(patsubst %,--proto_include_path=%,$(PROTO_INCLUDE_PATHS)) \
		"--plugin_name=validate" \
		"--plugin_out=$(PROTOC_GEN_VALIDATE_OUT)" \
		"--plugin_opt=$(PROTOC_GEN_VALIDATE_OPT)"

pregenerate:: protocgenvalidate
