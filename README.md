# apiServer + Kubernetes + Linkerd + ingress
# 0. Create apiserver 
# 1. Containerize the golang app and 
 * Install Docker
 * Create a docker image from [docker file](https://github.com/jonathan34c/apiServer/blob/main/Dockerfile)
    * ```docker build -t apiserver:one .```
 *  run ```docker image ls``` in terminal to see the listed images 
 * Give the image a tag by run ``` docker tag [image name] [tagname]```
 # 2. deploy docker image to k8s 
 * login to azure ```az login```
 * create a azure group ```az group create --name [registory name] --location westus```
 * create azure container registory ```az acr create --resource-group [registory name] --name ```
 * login to azure container registry ```az acr login --name [registory name]```
 * tag image by docker ```docker tag apiserver:one [registory name].azurecr.io/[image name]:[tag name]```
 * push image to azure ```docker push [image]/[image]:[tag]```
 * list all the container images ```az acr repository list --name [registory name] --output table``` 
 ![Screenshot 2023-04-03 100243](https://user-images.githubusercontent.com/8307131/229578336-984f2d61-7a55-40ef-9561-fbc156b2342a.png)
# 3. create kubernetes pod
 * write up a [yaml file](https://github.com/jonathan34c/apiServer/blob/main/testpod.yaml)
 * create the pod in kubernetes using the yaml you just created ```kubectl create -f testpod.yaml```
 * test if the pod has created correctly, run ```kubectl get pods ``` to see all the current pods 
# 4. test communication between pods
  * create another [yaml file](https://github.com/jonathan34c/apiServer/blob/main/otherpod.yaml) using the same values and run ```kubectl create -f otherpod.yaml```
  * inspect target pod IP address by ``` kubectl describe pods testpod```
  
  ![Screenshot 2023-04-03 103135](https://user-images.githubusercontent.com/8307131/229584159-677aa658-e4de-41f9-872f-246b0e994f80.png)
  * test out pod to pod communication ```kubectl exec other -- curl http://[ip address get from previous step]:8080/apps```
  ![Screenshot 2023-03-31 202251](https://user-images.githubusercontent.com/8307131/229586125-721ad693-0a9e-474d-9695-bb7f6dc960dd.png)
# 5. inkect linkerd to kubernetes application
* install [linkerd](https://linkerd.io/2.12/tasks/install/)
* inject linkerd ```kubectl get deployments -n [namespace name] -o yaml | linkerd inject - | kubectl apply -f -```
* checkout linkerd dashboard at ```linkerd dashboard```
![Screenshot 2023-04-03 105506](https://user-images.githubusercontent.com/8307131/229588898-21778d87-32d4-495e-86b3-ce2854a67aa5.png)
* apply the authrization policy with https://github.com/jonathan34c/apiServer/blob/main/linker_auth_policy.yaml https://github.com/jonathan34c/apiServer/blob/main/linker_server.yaml https://github.com/jonathan34c/apiServer/blob/main/linker_service_account.yaml

* replace pod policy by adding server to pod yaml file
```
kubectl delete pod nginx2
kubectl run nginx2 --image=nginx --port=80 --overrides='{ "spec": { "serviceAccount": "curl-meshtls-test" }  }'
kubectl get pods -o yaml | linkerd inject - | kubectl replace --force -f -
``` 
* run the curl command in step 4 with the newly created pod to see the result 
![Screenshot 2023-04-03 163428](https://user-images.githubusercontent.com/8307131/229649078-af43773b-6104-45dc-a210-3ff29d04f102.png)


# 6. Certificates with Azure Key Vault and Nginx Ingress Controller
* install nginx ingress controller 
```
NAMESPACE=ingress-nginx

helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update

helm install ingress-nginx ingress-nginx/ingress-nginx \
  --create-namespace \
  --namespace $NAMESPACE \
  --set controller.service.annotations."service\.beta\.kubernetes\.io/azure-load-balancer-health-probe-request-path"=/healthz
 ```
 
* follow step to get the secret using [cli driver](https://learn.microsoft.com/en-us/azure/aks/csi-secrets-store-driver) 
* obtain xxx.crt and xxx.key by
```
export CERT_NAME=aks-ingress-cert
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -out aks-ingress-tls.crt \
    -keyout aks-ingress-tls.key \
    -subj "/CN=demo.azure.com/O=aks-ingress-tls"
```

```
export AKV_NAME="[YOUR AKV NAME]"
openssl pkcs12 -export -in aks-ingress-tls.crt -inkey aks-ingress-tls.key  -out $CERT_NAME.pfx
```

use cli to get the certificate
```
az keyvault certificate import --vault-name $AKV_NAME -n $CERT_NAME -f $CERT_NAME.pfx
```
* after obtain the crtificate and key, create your secret by 
```
kubectl create secret tls ingress-cert --namespace [namespace] --key=[keyname].key --cert=[certname].crt -o yaml
```
* double check if the scret has created correctly 
```
kubectl get secret
```

![Screenshot 2023-04-05 152438](https://user-images.githubusercontent.com/8307131/230225123-315755f5-88a1-4aba-a606-d829f98e061e.png)

* apply secret to your [ingress controller](https://github.com/jonathan34c/apiServer/blob/main/ingress.yaml). remember to add secret name in "secrectname" part

* apply changes, and see ithe external ip by ```kubectl get services -n [namespace]```

![Screenshot 2023-04-05 152705](https://user-images.githubusercontent.com/8307131/230225436-638864cd-b6e9-46be-a348-dbafce5df2d8.png)

* finally, check the coonection by ```curl -v -k --resolve [your difine address]:443:[external ip from last step] https://[your difine address]```

![Screenshot 2023-04-05 143740](https://user-images.githubusercontent.com/8307131/230225595-f16e4042-e773-4a1f-acb6-d727aec47db8.png)

# DNS
* get resource group name by 
```
kubectl --namespace ingress-basic get services -o wide -w ingress-nginx-controller
```
* create public static ip
```
az network public-ip create --resource-group MC_myResourceGroup_myAKSCluster_eastus --name myAKSPublicIP --sku Standard --allocation-method static --query publicIp.ipAddress -o tsv
```
* set the domain name 
```
DNS_LABEL="<DNS_LABEL>"
NAMESPACE="ingress-basic"
STATIC_IP=<STATIC_IP>

helm upgrade ingress-nginx ingress-nginx/ingress-nginx \
  --namespace $NAMESPACE \
  --set controller.service.annotations."service\.beta\.kubernetes\.io/azure-dns-label-name"=$DNS_LABEL \
  --set controller.service.loadBalancerIP=$STATIC_IP
```
* visit your set domaing by http://[domain name].westus3.cloudapp.azure.com



# Create secret for tls
* follows these steps before you start https://learn.microsoft.com/en-us/azure/aks/csi-secrets-store-driver https://learn.microsoft.com/en-us/azure/aks/csi-secrets-store-identity-access https://learn.microsoft.com/en-us/azure/aks/csi-secrets-store-nginx-tls
* create a secret provider class https://github.com/jonathan34c/apiServer/blob/main/ingress_secret_provider.yaml
* make sure you have keyvault yaml as well https://github.com/jonathan34c/apiServer/blob/main/keyvault.yaml
* make sure to include "userAssignedIdentityID:", the tutorial has left out this part. 
* folow [instruction](https://learn.microsoft.com/en-us/azure/aks/ingress-tls?tabs=azure-cli) to create secret for tls
* note that since the helm has updated, the secret type will be "helm.sh/release.v1" instead of  kubernetes.io/tls
* add the secret to tls [secret](https://github.com/jonathan34c/apiServer/blob/83cf8cc3e8c7c692a49ab0e8ac86a537b4efbf33/ingress.yaml#L18)
* apply ingress.yaml
* visit the dns domain to test it out

# Client certifitcate 
* create certificate authority``` openssl req -x509 -sha256 -newkey rsa:4096 -keyout ca.key -out ca.crt -days 356 -nodes -subj '/CN=My Cert Authority'''' 
* apply CA as secret ```kubectl create secret generic ca-secret --from-file=ca.crt=ca.crt -n [namespace]```
* create cert sigiing request ``` openssl req -new -newkey rsa:4096 -keyout client.key -out client.csr -nodes -subj ‘/CN=My Client’```
* add following to [ingress.yaml](https://github.com/jonathan34c/apiServer/blob/83cf8cc3e8c7c692a49ab0e8ac86a537b4efbf33/ingress.yaml#L9) ``` nginx.ingress.kubernetes.io/auth-tls-pass-certificate-to-upstream: "true"
    nginx.ingress.kubernetes.io/auth-tls-secret: default/ca-secret
    nginx.ingress.kubernetes.io/auth-tls-verify-client: "on"
    nginx.ingress.kubernetes.io/auth-tls-verify-depth: "1"```
* test tout by first visit the domain, it would request a certificate ``` curl -k https://aroga.westus3.cloudapp.azure.com/ga```  ![Screenshot 2023-04-10 201843](https://user-images.githubusercontent.com/8307131/231049772-32f6782e-75f3-4cb6-9626-616063a08f77.png) 
* add certificate to the curl ``` curl -k https://aroga.westus3.cloudapp.azure.com/ga --key client.key --cert client.crt``` you will be able to visit the pod

 ![Screenshot 2023-04-10 201836](https://user-images.githubusercontent.com/8307131/231049979-5385a505-b86a-48fb-9a71-8245397d1197.png)



# Trouble shooting 
* if encounter ```ImagePullBackOff``` error for pod. Remember need to verified your docker image on azure using ```az aks update -n [resource-group name] -g [registry name] --attach-acr [registry name] ```

* if you had configured the azure key vault certificate wrong, you will get a secret like this. #DO NOT use this kind of secret 
```
sh.helm.release.v1.ingress-nginx.v4   helm.sh/release.v1   1      15h
``` 
# Reference 

* https://devopscube.com/kubernetes-ingress-tutorial/
* https://www.digitalocean.com/community/tutorials/how-to-install-and-use-linkerd-with-kubernetes
* https://blog.baeke.info/2020/12/07/certificates-with-azure-key-vault-and-nginx-ingress-controller/
* https://blog.devops.dev/linkerd-authorization-policies-authorization-as-a-microservice-d48675512c0a
* https://snyk.io/blog/setting-up-ssl-tls-for-kubernetes-ingress/
