# system-education Helm chart

Este diretório contém um chart Helm para automatizar o deploy dos recursos Kubernetes utilizados pelo projeto `system-education`.

## Pré-requisitos

- [Helm 3](https://helm.sh/docs/intro/install/)
- Um cluster Kubernetes com acesso a um registry que contenha a imagem da aplicação `system-education`
- Opcional: `kubectl` configurado no contexto do cluster para validações manuais

## Estrutura do chart

- `Chart.yaml`: metadados do chart
- `values.yaml`: valores padrão utilizados pelos templates. Aqui você pode definir o repositório, tag e política de pull da imagem da aplicação, além de configurações do banco de dados PostgreSQL
- Diretório `templates/`: contém todos os manifests renderizados pelo Helm. Os arquivos criam ConfigMap, Secret, Deployment e Service da aplicação, além de (opcionalmente) os recursos para executar uma instância do PostgreSQL dentro do cluster

## Instalação

1. Atualize o `values.yaml` ou utilize parâmetros `--set` para configurar o deploy. Por exemplo, para apontar para uma imagem publicada no registry `ghcr.io/<sua-org>/system-education` com a tag `1.0.0`:

   ```bash
   helm upgrade --install system-education helm/system-education \
     --namespace system-education --create-namespace \
     --set image.repository=ghcr.io/<sua-org>/system-education \
     --set image.tag=1.0.0
   ```

   > **Dica:** defina `image.pullPolicy=Always` (valor padrão) para garantir que o cluster sempre busque a última versão disponibilizada com a tag informada.

2. Para atualizar a aplicação após publicar uma nova imagem, execute novamente o comando `helm upgrade` informando a nova tag (por exemplo, `--set image.tag=1.1.0`). O Helm irá atualizar o Deployment e forçar o Kubernetes a recriar os pods com a nova imagem.

3. Caso utilize o PostgreSQL embarcado, os dados serão armazenados em um `PersistentVolumeClaim`. Para desabilitar a criação desse banco, utilize `--set postgresql.enabled=false` e ajuste as variáveis de ambiente conforme o seu banco externo.

## Como verificar o deploy

- Listar os pods:

  ```bash
  kubectl get pods -n system-education
  ```

- Fazer port-forward para acessar a aplicação localmente:

  ```bash
  kubectl port-forward svc/system-education 8080:8080 -n system-education
  ```

- Visualizar os valores atualmente aplicados pelo Helm:

  ```bash
  helm get values system-education -n system-education
  ```

Com isso, o deploy do Kubernetes passa a ser gerenciado pelo Helm, simplificando atualizações de versão da imagem e manutenção dos manifestos.
