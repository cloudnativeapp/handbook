### Helm Chart 创作指南

Helm作为当前最流行的Kubernetes应用管理工具之一，整合应用部署所需的K8s资源（包括Deployment，Service等）到Chart中。今天，本文会带领大家学习如何创建一个简单的Chart。

#### 准备

本文基于的是最新的[Helm v3.0.0-alpha.1](https://v3.helm.sh/)，Helm v3相较于Helm v2有较大的变动，比如Helm v3没有服务器端的Tiller组件，在Chart的使用上接口也有一些变化。建议没有Helm v3的同学先从Helm的Github仓库中下载最新的[Helm v3 release](https://github.com/helm/helm/releases/tag/v3.0.0-alpha.1)。想要仔细钻研Helm v3各种特性的同学，可以参考Helm v3的完整[官方文档](https://v3.helm.sh/docs/using_helm/)。

在下载得到最新的Helm v3后，同学们可以运行`helm init`来初始化Helm。

#### 开始创作

首先，我们需要有一个要部署的应用。这里我们使用一个简单的基于golang的[hello world HTTP服务](https://github.com/cloudnativeapp/handbook/tree/master/helm-chart-creation-tutorial/src/main.go)。该服务通过读取环境变量`USERNAME`获得用户自己定义的名称，然后监听80端口。对于任意HTTP请求，返回`Hello ${USERNAME}。`比如如果设置`USERNAME=world`（默认场景），该服务会返回`Hello world`。

准备好要部署的应用镜像后，运行`helm create my-hello-world`，便会得到一个helm自动生成的空chart。这个chart里的名称是`my-hello-world`。
**需要注意的是，Chart里面的my-hello-world名称需要和生成的Chart文件夹名称一致。如果修改my-hello-world，则需要做一致的修改。**
现在，我们看到Chart的文件夹目录如下

```yaml
my-hello-world
├── charts
├── Chart.yaml
├── templates
│   ├── deployment.yaml
│   ├── _helpers.tpl
│   ├── ingress.yaml
│   ├── NOTES.txt
│   └── service.yaml
└── values.yaml
```

在根目录下的Chart.yaml文件内，声明了当前Chart的名称、版本等基本信息，这些信息会在该Chart被放入仓库后，供用户浏览检索。比如我们可以把Chart的Description改成"My first hello world helm chart"。

#### 走近Chart

Helm Chart对于应用的打包，不仅仅是将Deployment和Service以及其它资源整合在一起。我们看到deployment.yaml和service.yaml文件被放在templates/文件夹下，相较于原生的Kubernetes配置，多了很多渲染所用的可注入字段。比如在deployment.yaml的`spec.replicas`中，使用的是`.Values.replicaCount`而不是Kubernetes本身的静态数值。这个用来控制应用在Kubernetes上应该有多少运行副本的字段，在不同的应用部署环境下可以有不同的数值，而这个数值便是由注入的`Values`提供。

在根目录下我们看到有一个`values.yaml`文件，这个文件提供了应用在安装时的默认参数。在默认的`Values`中，我们看到`replicaCount: 1`说明该应用在默认部署的状态下只有一个副本。

为了使用我们要部署应用的镜像，我们看到deployment.yaml里在`spec.template.spec.containers`里，`image`和`imagePullPolicy`都使用了`Values`中的值。其中`image`字段由`.Values.image.repository`和`.Chart.AppVersion`组成。看到这里，同学们应该就知道我们需要变更的字段了，一个是位于values.yaml内的`image.repository`，另一个是位于Chart.yaml里的`AppVersion`。我们将它们与我们需要部署应用的docker镜像匹配起来。这里我们把values.yaml里的`image.repository`设置成`somefive/hello-world`，把Chart.yaml里的`AppVersion`设置成`1.0.0`即可。

类似的，我们可以查看service.yaml内我们要部署的服务，其中的主要配置也在values.yaml中。默认生成的服务将80端口暴露在Kubernetes集群内部。我们暂时不需要对这一部分进行修改。

由于部署的hello-world服务会从环境变量中读取`USERNAME`环境变量，我们将这个配置加入deployment.yaml。相关部分如下：

```yaml
- name: {{ .Chart.Name }}
  image: "{{ .Values.image.repository }}:{{ .Chart.AppVersion }}"
  imagePullPolicy: {{ .Values.image.pullPolicy }}
  env:
    - name: USERNAME
      value: {{ .Values.Username }}
```

现在我们的deployment.yaml模版会从values.yaml中加载`Username`字段，因此相应的，我们也在values.yaml中添加`Username: AppHub`。

#### 打包使用

完成上述配置后，我们可以使用`helm lint --strict my-hello-world`进行格式检查。如果显示

```bash
1 chart(s) linted, 0 chart(s) failed
	[INFO] Chart.yaml: icon is recommended
```

那么我们就已经离成功只差一步之遥了。

接下来，我们运行`helm package my-hello-world`指令对我们的Chart文件夹进行打包。现在我们就得到了`my-hello-world-0.1.0.tgz`的Chart包。到这一步我们的Chart便已经完成了。

之后，运行`helm install my-hello-world-chart-test my-hello-world-0.1.0.tgz`来将本地的chart安装到my-hello-world-chart-test的Release中。运行`kubectl get pods`我们可以看到要部署的pod已经处于运行状态

```bash
NAME                                         READY   STATUS    RESTARTS   AGE
my-hello-world-chart-test-65d6c7b4b6-ptk4x   1/1     Running   0          4m3s
```

运行`kubectl port-forward my-hello-world-chart-test-65d6c7b4b6-ptk4x 8080:80`后，就可以直接在本地运行`curl localhost:8080`看到`Hello AppHub`了！

#### 进阶使用

上述提到values.yaml只是Helm install参数的默认设置，我们可以在安装Chart的过程中使用自己的参数覆盖。比如我们可以运行`helm install my-hello-world-chart-test2 my-hello-world-0.1.0.tgz --set Username="Cloud Native"`来安装一个新Chart。同样运行`kubectl port-forward`进行端口映射，这时可以得到`Hello Cloud Native`。

我们注意到在安装Chart指令运行后，屏幕的输出会出现

```bash
NOTES:
1. Get the application URL by running these commands:
  export POD_NAME=$(kubectl get pods -l "app=my-hello-world,release=my-hello-world-chart-test2" -o jsonpath="{.items[0].metadata.name}")
  echo "Visit http://127.0.0.1:8080 to use your application"
  kubectl port-forward $POD_NAME 8080:80
```

这里的注释是由Chart中的`templates/NOTES.txt`提供的。我们注意到原始的NOTES中，所写的`"app={{ template "my-hello-world.name" . }},release={{ .Release.Name }}"`和我们的deployment.yaml中所写的配置不太一样。我们可以把它改成`"app.kubernetes.io/name={{ template "my-hello-world.name" . }},app.kubernetes.io/instance={{ .Release.Name }}"`，将values.yaml中的`version`更新成`0.1.1`。然后重新打包Chart（运行`helm package`）。得到新的my-hello-world-0.1.1.tgz之后，重新安装Chart（运行`helm install my-hello-world-chart-test3 my-hello-world-0.1.1.tgz --set Username="New Chart"`），就能看到更新过后的NOTES了。

```bash
NAME: my-hello-world-chart-test3
LAST DEPLOYED: 2019-07-10 14:02:55.321468411 +0800 CST m=+0.091032750
NAMESPACE: default
STATUS: deployed

NOTES:
1. Get the application URL by running these commands:
  export POD_NAME=$(kubectl get pods -l "app.kubernetes.io/name=my-hello-world,app.kubernetes.io/instance=my-hello-world-chart-test3" -o jsonpath="{.items[0].metadata.name}")
  echo "Visit http://127.0.0.1:8080 to use your application"
  kubectl port-forward $POD_NAME 8080:80
```

#### 其他

Helm Chart还有诸如dependency等其他功能，更加详细的资料可以参考Helm官方文档的[相关章节](https://v3.helm.sh/docs/topics/chart_template_guide/)。
