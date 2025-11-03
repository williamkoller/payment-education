# Kubernetes deployment

Este diretório contém os manifestos Kubernetes necessários para executar o projeto `system-education` em um cluster local (por exemplo, usando o Minikube).

## Recursos

- `configmap.yaml`: define as variáveis de ambiente da aplicação e da conexão com o banco de dados.
- `secret.yaml`: armazena a senha do banco de dados utilizada por aplicação e banco.
- `deployment.yaml`: cria o deployment que executa o container com a aplicação na porta 8080.
- `service.yaml`: expõe a aplicação internamente no cluster na porta 8080.
- `postgres-deployment.yaml`: provisiona um banco PostgreSQL com `PersistentVolumeClaim` e serviço interno.


## Executando com Minikube

1. Inicie o Minikube e configure seu shell para utilizar o Docker interno do cluster:

   ```bash
   minikube start
   eval "$(minikube -p minikube docker-env)"
   ```

2. Faça o build da imagem da aplicação utilizando o Docker do Minikube (o Dockerfile expõe a porta 8080):

   ```bash
   docker build -t system-education:latest .
   ```

3. Ajuste os valores de banco de dados conforme o seu ambiente (por exemplo, edite `configmap.yaml` e `secret.yaml` ou utilize `kubectl create secret generic system-education-secret --from-literal=DB_PASSWORD=<sua-senha>` e atualize o ConfigMap antes de aplicar).

4. Aplique os manifestos Kubernetes:

   ```bash
   kubectl apply -f k8s/
   ```

5. Verifique se os pods estão prontos (aplicação e banco):

   ```bash
   kubectl get pods
   ```

6. Para acessar a aplicação fora do cluster, exponha o serviço via `port-forward` ou crie um `NodePort`. Exemplo com `port-forward`:


   ```bash
   kubectl port-forward service/system-education 8080:8080
   ```

   Depois disso, a aplicação estará acessível em `http://localhost:8080`.

> **Observação:** os manifestos assumem um banco de dados PostgreSQL acessível a partir do cluster com as credenciais definidas no `ConfigMap` e `Secret`. Ajuste os valores conforme o ambiente desejado.
