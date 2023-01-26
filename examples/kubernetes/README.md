
# Sample Traefik Kubernetes configuration

A kubernetes config for the image build with `Dockerfile`

```yaml
apiVersion: v1
items:
- apiVersion: apps/v1
  kind: Deployment
  metadata:
    labels:
      app: pub
    name: pub
    namespace: masto
  spec:
    replicas: 1
    selector:
      matchLabels:
        app: pub
    template:
	  labels:
          app: pub
      spec:
        containers:
        - image: myregistry/akh/pub:latest
          imagePullPolicy: Always
          name: pub
          ports:
          - containerPort: 9999
            name: http
            protocol: TCP
          resources:
            limits:
              cpu: 500m
              memory: 64Mi
            requests:
              cpu: 250m
              memory: 32Mi
          volumeMounts:
          - mountPath: /data
            name: data
        restartPolicy: Always
        securityContext:
          fsGroup: 2000
          runAsGroup: 3000
          runAsNonRoot: true
          runAsUser: 1000
        terminationGracePeriodSeconds: 30
        volumes:
        - hostPath:
            path: /opt/data/masto
            type: Directory
          name: data
---
apiVersion: v1
items:
- apiVersion: traefik.containo.us/v1alpha1
  kind: IngressRoute
  metadata:
    name: pub-ingress
    namespace: masto
  spec:
    entryPoints:
    - websecure
    routes:
    - kind: Rule
      match: Host(`masto.my.domain`)
      middlewares:
      - name: compress
      services:
      - name: pub
        port: 9999
    tls:
      store:
        name: default
kind: List
metadata:
  resourceVersion: ""
---