package pkg

const (
	DefaultNamespace                = "default"
	DefaultVolumeName               = "data"
	DefaultInitContainerImage       = "alpine/git"
	DefaultInitContainerName        = "init"
	DefaultInitContainerWorkspace   = "/opt/workspace"
	DefaultContainerWorkspace       = "/opt/workspace"
	DefaultSwapWorkspace            = "/dtswap"
	DefaultSwapVolumeName           = "dtswap"
	DefaultAnnotationsDtRunner      = "apps.dtwave.com/dt-runner"
	DefaultAnnotationsDtRunnerCi    = "apps.dtwave.com/dt-runner/ci"
	DefaultAnnotationsDtRunnerModel = "apps.dtwave.com/dt-runner/model"
	DefaultLabelDtRunner            = "apps.dtwave.com/dt-runner"
	DefaultLabelDtRunnerCi          = "ci.apps.dtwave.com"
	DefaultLabelDtRunnerModel       = "model.apps.dtwave.com"
)
