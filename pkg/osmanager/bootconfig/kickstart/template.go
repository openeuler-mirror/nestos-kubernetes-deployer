package kickstart

import (
	"fmt"
	"nestos-kubernetes-deployer/data"
	"nestos-kubernetes-deployer/pkg/configmanager/asset"
	"nestos-kubernetes-deployer/pkg/configmanager/runtime"
	"nestos-kubernetes-deployer/pkg/constants"
	"nestos-kubernetes-deployer/pkg/osmanager/bootconfig"
	"nestos-kubernetes-deployer/pkg/utils"

	"github.com/sirupsen/logrus"
)

type KsTempData struct {
	Hostname string
	Password string
	Files    []File
	Systemds []string
	IsDocker bool
	IsIsulad bool
}

type File struct {
	ChangeMod string
	Content   string
}

type template struct {
	clusterAsset    *asset.ClusterAsset
	config          []byte
	enabledServices []string
	enabledFiles    []string
}

func newTemplate(clusterAsset *asset.ClusterAsset, enabledServices, enabledFiles []string) *template {
	return &template{
		clusterAsset:    clusterAsset,
		enabledServices: enabledServices,
		enabledFiles:    enabledFiles,
	}
}

func (t *template) GenerateBootConfig(url string, nodeType string) error {
	files := []bootconfig.File{}
	systemds := bootconfig.Systemd{}
	ksData := KsTempData{Password: t.clusterAsset.Password}

	tmplData, err := bootconfig.GetTmplData(t.clusterAsset)
	if err != nil {
		return fmt.Errorf("failed to get template data: %v", err)
	}

	if nodeType == constants.Controlplane {
		tmplData.IsControlPlane = true
		tmplData.CertsUrl = utils.ConstructURL(url, constants.CertsFiles)

		ksData.Hostname = t.clusterAsset.Master[0].Hostname
	} else if nodeType == constants.Master {
		for i := 1; i < len(t.clusterAsset.Master); i++ {
			ksData.Hostname = t.clusterAsset.Master[i].Hostname
		}
	}

	engine, err := runtime.GetRuntime(t.clusterAsset.Runtime)
	if err != nil {
		return fmt.Errorf("failed to get runtime: %v", err)
	}

	tmplData.CriSocket = engine.GetRuntimeCriSocket()
	if runtime.IsIsulad(engine) {
		t.enabledFiles = append(t.enabledFiles, constants.IsuladConfig)
		ksData.IsIsulad = true
	} else if runtime.IsDocker(engine) {
		t.enabledFiles = append(t.enabledFiles, constants.DockerConfig)
		ksData.IsDocker = true
	}

	if err := bootconfig.AppendStorageFiles(&files, "/", constants.BootConfigFilesPath, tmplData, t.enabledFiles); err != nil {
		logrus.Errorf("failed to add files to a kickstart config: %v", err)
		return err
	}

	if err := bootconfig.AppendSystemdUnits(&systemds, constants.BootConfigSystemdPath, tmplData, t.enabledServices); err != nil {
		logrus.Errorf("failed to add systemd units to a kickstart config: %v", err)
		return err
	}

	for _, f := range files {
		ksData.Files = append(ksData.Files, File{
			Content:   fmt.Sprintf("cat <<'EOF'> %s\n%s\nEOF", f.Path, string(f.Contents.Source)),
			ChangeMod: fmt.Sprintf("chmod %o %s", f.Mode, f.Path),
		})
	}

	if len(t.clusterAsset.HookConf.ShellFiles) > 0 {
		for _, sf := range t.clusterAsset.HookConf.ShellFiles {
			hookFilePath := constants.HookFilesPath + sf.Name
			ksData.Files = append(ksData.Files, File{
				Content:   fmt.Sprintf("cat <<'EOF'> %s\n%s\nEOF", hookFilePath, string(sf.Content)),
				ChangeMod: fmt.Sprintf("chmod %o %s", sf.Mode, hookFilePath),
			})
		}
	}

	for _, u := range systemds.Units {
		fp := fmt.Sprintf("/etc/systemd/system/%s", u.Name)
		ksData.Files = append(ksData.Files, File{
			Content:   fmt.Sprintf("cat <<'EOF'> %s\n%s\nEOF", fp, string(u.Contents)),
			ChangeMod: fmt.Sprintf("chmod %o %s", constants.BootConfigFileMode, fp),
		})
		ksData.Systemds = append(ksData.Systemds, fmt.Sprintf("systemctl enable %s", u.Name))
	}
	ksData.Systemds = append(ksData.Systemds, fmt.Sprintf("systemctl enable %s", constants.KubeletService))

	file, err := data.Assets.Open("kickstart/" + nodeType + "/kickstart.cfg.template")
	if err != nil {
		return fmt.Errorf("failed to open kickstart template file: %v", err)
	}
	defer file.Close()

	_, data, err := utils.GetCompleteFile("kickstart.cfg.template", file, ksData)
	if err != nil {
		return fmt.Errorf("failed to generate kickstart file: %v", err)
	}

	t.config = data
	return nil
}
