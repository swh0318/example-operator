# sample-operator
This is an sample operator created by using the kubebuilder

## step by step to build your operator
### 1、# kubebuilder init --domain github.com --skip-go-version-check 

### 2、# kubebuilder create api --group crd --version v1beta1 --kind MyCRD
### 3、定义CRD：对应 api/v1beta1/mycrd_types.go文件
### 4、开发控制器逻辑：controllers/mycrd_controller.go 
注意：该 Operator 须要创立 Pod，因而须要给该 Operator 创立 Pod 的权限，Kubebuilder 反对主动生成 Operator 的 RBAC，
然而须要开发者在管制逻辑加上标识，此处咱们加上对 Pod 有读写的权限的标识：
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete
### 5、# go run main.go -kubeconfig=/Users/sunwenhua/.kube/config
所有需要自定义的地方，我都在sample工程中注释'// TODO(user):用户自定义逻辑'，大家可以通过搜索该注释，修改自定义自己的结构或逻辑。
