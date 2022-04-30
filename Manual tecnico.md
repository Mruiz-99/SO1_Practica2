# Manual técnico </center>

### Servicios utilizados:
* Clústeres de Kubernetes en Google Cloud Platform
* Virtual Machine
  
### Tecnologías utilizadas:
* Go
* Nodejs
* Docker
* Kafka
* Mongo DB
  

# <img src="https://cursosdedesarrollo.com/wp-content/uploads/2020/03/k-logo.png" width="50"> Kubernetes

Kubernetes es un sistema de código libre para la automatización del despliegue, ajuste de escala y manejo de aplicaciones en contenedores que fue originalmente 
diseñado por Google y donado a la Cloud Native Computing Foundation. Soporta diferentes entornos para la ejecución de contenedores, incluido Docker.

Kubernetes es una plataforma portable y extensible para administrar cargas de trabajo y servicios. Kubernetes tiene un ecosistema grande y en rápido
crecimiento. El soporte, las herramientas y los servicios para Kubernetes están ampliamente disponibles.

<br>

> ## Implementando Kubernetes 
En esta practica de utilizo Google Cloud como plataforma para poder utilizar sus servicios, tenemos que ingresar a cloud shell. 
Y luego seleccionamos el proyecto en el que deseamos trabajar(o se crea un proyecto) con kubernetes con el siguiente comando:

```
gcloud config set project <ID del proyecto>
```

luego seguimos la configuracion con el siguiente comando, en este comando le enviamos la zona donde queremos trabajar:

```
gcloud config set compute/zone us-central1-a
```
 Ahora vamos a crear nuestro cluster de kubernetes con el siguiente comando:

```
gcloud container clusters create cluster-practica2-mynor --num-nodes=1 --tags=allin,allout --machine-type=n1-standard-2 --no-enable-network-policy
```
### Donde:

* Nombre del cluster: cluster-practica2-mynor
* Numero de nodos: --num-nodes=1
* Tipo de VM (2 CPUs, 8GB RAM): --machine-type=n1-standard-2
* Networks rules (allin, allout): --tags=allin,allout
* Autenticacion con certificado: --enable-legacy-authorization

Luego se creara archivos Kafka_Service.yml, deployment.yml y Kafka_consume.yml.

Antes de crear los archivos se debe ingresar los siguientes comandos:

```
kubectl create namespace kafka
kubectl create -f ' https://strimzi.io/install/latest?namespace=kafka' -n kafka
```
En esta parte creamos los servicios de kafka para luego solo implementarlos en nuestro namspace del proyecto. Ahora si implementamos el archivo ***.yml*** el cual ingresamos el comando:

```
nano Kafka_Service.yml
```
en el cual ponemos el siguiente contenido:

```yml
apiVersion: kafka.strimzi.io/v1beta2
kind: Kafka
metadata:
  name: my-cluster
spec:
  kafka:
    version: 3.1.0
    replicas: 3
    listeners:
      - name: plain
        port: 9092
        type: internal
        tls: false
      - name: tls
        port: 9093
        type: internal
        tls: true
    config:
      offsets.topic.replication.factor: 3
      transaction.state.log.replication.factor: 3
      transaction.state.log.min.isr: 2
      log.message.format.version: "3.1"
      inter.broker.protocol.version: "3.1"
    storage:
      type: ephemeral
  zookeeper:
    replicas: 3
    storage:
      type: ephemeral
  entityOperator:
    topicOperator: {}
    userOperator: {}
```
Guardamos y procedemos a ejecutar con el siguiente comando:

```
kubectl apply -f kafka_service.yml -n kafka
```
Este comando nos transforma nuestro archivo yml a yaml el cual estara ejecutando kubernetes, una vez termina de hacer la compilacion, procedemos a probar nuestro servicio con los siguientes comandos:

```
kubectl -n kafka run kafka-producer -ti --image=quay.io/strimzi/kafka:0.28.0-kafka-3.1.0 --rm=true --restart=Never -- bin/kafka-console-producer.sh --bootstrap-server my-cluster-kafka-bootstrap:9092 --topic my-topic
```
Este comando nos siver para probar si el servicio del producer funciona el cual ingresamos cualquier cosa, para luego ingresar el siguiente comando:

```
kubectl -n kafka run kafka-consumer -ti --image=quay.io/strimzi/kafka:0.28.0-kafka-3.1.0 --rm=true --restart=Never -- bin/kafka-console-consumer.sh --bootstrap-server my-cluster-kafka-bootstrap:9092 --topic my-topic --from-beginning
```
Este nos mostrara la cola que recibe del producer. Una vez que tenemos esto implementado, pasamos a crear nuestro archivo de deployment, el cual contentra lo siguiente:

```apiVersion: apps/v1
kind: Deployment
metadata:
  name: grpc-practica2
  namespace: practica2-201801329
  labels:
    app: grpc-practica2
spec:
  selector:
    matchLabels:
      app: grpc-practica2
  replicas: 1
  template:
    metadata:
      labels:
        app: grpc-practica2
    spec:
      hostname: grpc-pod-host
      containers:
        - name: grpc-client
          image: mruiz01329/grpc-client-api
          env:
          - name: SERVER_ADDR
            value: grpc-pod-host
          - name: SERVER_PORT
            value: "50051"
          - name: CLIENTE_PORT
            value: "2000"
          ports:
            - containerPort: 2000
        - name: grpc-server
          image: mruiz01329/grpc-server
          env:
          - name: PORT
            value: "50051"
          - name: BROKER_ADDR
            value: "my-cluster-kafka-bootstrap.kafka:9092"
          ports:
            - containerPort: 50051
```
una vez tenemos termiando nuestro servicio de grpc junto con el producer de kafka procedemos a ejecutar la compilacion de nuestro archivo .yml con el comando:

```
kubectl apply -f deployment.yml -n practica2-201801329
```
Este nos creara nuestro servicio de grpc, el cual aun tenemos que exponer con el servicio de load balancer, el cual encontramos en Cargas de trabajo -> Seleccionamos 
nuestra carga de trabajo -> Accines -> Exponer -> asignamos el puero de nuestro cliente -> tipo de servicio seleccionamos Load balancer y listo.

Una vez terminamos de exponer nuestro deployment, agregamos nuestro consumer de la misma manera que agregamos nuestro deployment, con el siguiente contenido de nuestro archivo:

```apiVersion: apps/v1
kind: Deployment
metadata:
  name: kafka-consumer
  namespace: consumer
  labels:
    app: kafka-consumer
spec:
  selector:
    matchLabels:
      app: kafka-consumer
  replicas: 1
  template:
    metadata:
      labels:
        app: kafka-consumer
    spec:
      hostname: kafka-consumer
      containers:
        - name: grpc-client
          image: mruiz01329/consumer_kafka
          env:
          - name: BROKER_ADDR
            value: my-cluster-kafka-bootstrap.kafka

```
Compilamos el archivo de la misma manera
```
kubectl apply -f kafka-consumer.yml -n practica2-201801329
```
