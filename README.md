# manifest-summary

Read a Kubernetes manifest and print a table of its contents.

```
$ curl -s https://raw.githubusercontent.com/kubernetes/website/master/content/en/examples/application/wordpress/mysql-deployment.yaml | manifest-summary
+-------------+-----------------------+-----------------+
| API VERSION |         KIND          |      NAME       |
+-------------+-----------------------+-----------------+
| v1          | Service               | wordpress-mysql |
| v1          | PersistentVolumeClaim | mysql-pv-claim  |
| apps/v1     | Deployment            | wordpress-mysql |
+-------------+-----------------------+-----------------+
```
