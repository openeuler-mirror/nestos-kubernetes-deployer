# Housekeeper

## overview

In the cloud-native domain, business deployment and operation are mainly achieved through container technology and container orchestration technology, which highly decouples them from the underlying environment. However, this also brings about the problem of inconsistent operation and maintenance (O&M) technology stacks, leading to separate management of Kubernetes (k8s) and the underlying operating system, as well as redundant construction of O&M platforms. To address these issues, NKD integrates the housekeeper module, ensuring consistency in operation and maintenance between business and the NestOS cloud base operating system. Operational management is conducted through containerization. The primary update process of housekeeper involves rebuilding the new version system image using image construction tools when the operating system or k8s basic components require upgrade and maintenance. After discovering the new version image, it creates housekeeper CR resources in the cluster. Housekeeper services in the cluster sequentially upgrade cluster nodes according to the configuration, completing the entire cluster's upgrade process.

## custom resource
### Update Resources
- Authorization: Permission control through RBAC.
- Explanation of CRD Resource Object Parameters:
  |  Parameter       | Type  |  Description                                          | Usage Note | Required         |
  | -------------- | ------  | -----------------------------------------------------------| ----- | ---------------- |
  | osImageURL | string  | Address for upgrading container images | Should be in the format REPOSITORY/NAME[:TAG@DIGEST] | Yes |
  | kubeVersion  | string  | Version number for upgrading Kubernetes | Leave empty if only upgrading the OS version | No         |
  | evictPodForce | bool | Force eviction of Pods, may lead to data loss or service interruption, use with caution | Default: false | No |
  | maxUnavailable  | int  | Maximum number of nodes for upgrade |Maximum number of nodes to be upgraded simultaneously  | No  |

## Architecture Introduction
housekeeper's architecture is shown:
![housekeeper-arch](/docs/en/figures/housekeeper-arch.jpg)
As shown in the diagram, housekeeper mainly consists of three components: housekeeper-operator-manager, housekeeper-controller-manager, and housekeeper-daemon.
- housekeeper-operator-manager: Running in the form of a Deployment on the Master node, responsible for coordinating all Machines for upgrades (not directly responsible for updates) and marking nodes ready for upgrade.
- housekeeper-controller-manager：Running in the form of a DaemonSet on all nodes in the cluster, responsible for evicting business pods and forwarding upgrade information to housekeeper-daemon.
- housekeeper-daemon: Receives information from housekeeper-controller-manager and performs atomic updates of the OS or upgrades Kubernetes version according to instructions
