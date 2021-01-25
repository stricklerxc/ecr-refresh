# ECR Refresh
A Kubernetes utility for refreshing expired ECR tokens in your Docker Registry secrets.

## Overview

Creates a CronJob object in the desired Kubernetes namespace. This CronJob executes this GO utility every 12 hours to refresh the selected secret with a new ECR token.

### To Deploy:

1. Clone repository

    ```bash
    $ git clone https://github.com/stricklerxc/ecr-refresh.git
    $ cd ecr-refresh
    ```

2. Deploy manifests

    ```bash
    $ kubectl apply -k . -n <namespace>
    serviceaccount/svc-ecr-refresh created
    role.rbac.authorization.k8s.io/ecr-refresh-role created
    rolebinding.rbac.authorization.k8s.io/svc-ecr-refresh created
    configmap/ecr-refresh created
    secret/ecr-creds configured
    cronjob.batch/ecr-refresh created
    ```

3. Configure ConfigMap
   - Create a programmatic IAM user with the following policy:
        ```json
        {
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Sid": "VisualEditor0",
                    "Effect": "Allow",
                    "Action": [
                        "ecr:GetDownloadUrlForLayer",
                        "ecr:BatchGetImage",
                        "ecr:BatchCheckLayerAvailability",
                        "ecr:CompleteLayerUpload",
                        "ecr:GetAuthorizationToken",
                        "ecr:UploadLayerPart",
                        "ecr:InitiateLayerUpload",
                        "ecr:PutImage"
                    ],
                    "Resource": "*"
                }
            ]
        }
        ```
    - Add AWS access keys to ConfigMap