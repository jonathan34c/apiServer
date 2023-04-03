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
* apply the authrization policy with https://github.com/jonathan34c/apiServer/blob/main/testauth.yaml https://github.com/jonathan34c/apiServer/blob/main/testcurl.yaml and https://github.com/jonathan34c/apiServer/blob/main/testnginxserver.yaml

* replace pod policy by adding server to pod yaml file
```
kubectl delete pod nginx2
kubectl run nginx2 --image=nginx --port=80 --overrides='{ "spec": { "serviceAccount": "curl-meshtls-test" }  }'
kubectl get pods -o yaml | linkerd inject - | kubectl replace --force -f -
``` 
* run the curl command in step 4 with the newly created pod to see the result 
![Screenshot 2023-04-03 163428](https://user-images.githubusercontent.com/8307131/229649078-af43773b-6104-45dc-a210-3ff29d04f102.png)


# 6. Certificates with Azure Key Vault and Nginx Ingress Controller
* [Install nginx ingress controller ](https://kubernetes.github.io/ingress-nginx/deploy/)
* [install akv2k8s controller](https://akv2k8s.io/)
* Creating the certificate in Azure Key Vault
* sync the certificate with [yaml](https://github.com/jonathan34c/apiServer/blob/main/testcertificate.yaml)
* check if the certificate has been sync using krew plugin ``` kubectl view-cert```

![Screenshot 2023-04-03 112732](https://user-images.githubusercontent.com/8307131/229596100-09696550-c3fc-4114-9b9f-d772c6a513aa.png)

# 7. Use the secret you sync with akv2k8s with ingress
* follow [testingress yaml](https://github.com/jonathan34c/apiServer/blob/main/testingress.yaml) and sync the certificate 
* after sync with akv2k8s, you can double check with ```kubectl get ingress```
![Screenshot 2023-04-03 113504](https://user-images.githubusercontent.com/8307131/229597025-ca925378-1dfe-4f9b-8cd2-53e97771f3ec.png)


# Trouble shooting 
* if encounter ```ImagePullBackOff``` error for pod. Remember need to verified your docker image on azure using ```az aks update -n [resource-group name] -g [registry name] --attach-acr [registry name] ```

# Reference 

* https://devopscube.com/kubernetes-ingress-tutorial/
* https://www.digitalocean.com/community/tutorials/how-to-install-and-use-linkerd-with-kubernetes
* https://blog.baeke.info/2020/12/07/certificates-with-azure-key-vault-and-nginx-ingress-controller/
* https://blog.devops.dev/linkerd-authorization-policies-authorization-as-a-microservice-d48675512c0a
