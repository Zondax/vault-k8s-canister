package common

const (
	TororuManagedAnnotation       = "tororu.zondax.io/managed"
	TororuResourceReqRWAnnotation = "tororu.zondax.io/secret-rw"
	TororuResourceReqROAnnotation = "tororu.zondax.io/secret-ro"
	// TODO: Resources will be created in diff namespace
	TororuEnvVarsAppliedAnnotation = "ro-applied"
	TororuAnnotationsTrueString    = "true"

	TororuResourceNamePlural = "tororu-resources"

	TororuResourcesGroup = "zondax.io"
	TororuResourcesName  = "tororu-resources.zondax.io"
)

const (
	ResourcePermissionTypeRW = "rw"
	ResourcePermissionTypeRO = "ro"
)

const (
	SidecarName = "tororu-sidecar"
)

type contextLoggerKey string

const (
	ContextLoggerKey contextLoggerKey = "logger"
)

const (
	PodVolumeName      = "shared-data"
	PodVolumeMountPath = "/pod-data"
)
