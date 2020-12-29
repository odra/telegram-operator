# Telegram Operator

This is an example operator used for learning purposes.

The operator will create a pod for each `BotMessage` CR created which will send a message toa Telegram channel.

## Telegram Configuration

The following steps will help you create a telegram bot:

1. Start a conversation with bot father: https://telegram.me/botfather

2. Type `/newbot`

3. Follow bot father instructions to create a new bot
   
4. Create a new public channel (you can make it a private one later) and add/invite you newly created bot there, all it needs is post message permissions.
   
5. Send a dummy message into your channel: `curl 'https://api.telegram.org/bot$API_TOKEN/sendMessage' -d 'chat_id=@$CHANNEL_NAME&text=hello'`
   
6. Expected value: `{"ok":true,"result":{"message_id":4,"chat":{"id":$CHAT_ID,"title":"asdf","username":"asdf","type":"channel"},"date":1564334481,"text":"hello"}}` (you need the chat_id value for our channel definition)
   
7. OPTIONAL: you can now make it a private channel if you want to.

## Usage

Create a generic secret containing both telegram chat id and bot api token:

```sh
kubectl \
create secret generic telegram-bot \
--from-literal=TG_BOT_TOKEN=CHANGEME \
--from-literal=TG_CHAT_ID=CHANGEME
```

Deploy the operator:

```sh
make install && make deploy
```

You can check for the operator readiness by running:

```sh
kubectl get po -w -n telegram-operator-system
```

You can apply the botmessage custom resource once the pod is ready:

```sh 
kubectl apply -f config/samples/telegram_v1alpha1_botmessage.yaml
```

You should see a new pod being created in the default namespace which will run send the telegram message.

You can check your custom resource status by running:

```sh
$ kubectl get botmessage/sample -o yaml
```

The output should be something like this:

```sh
apiVersion: telegram.my.domain/v1alpha1
kind: BotMessage
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"telegram.my.domain/v1alpha1","kind":"BotMessage","metadata":{"annotations":{},"name":"sample","namespace":"default"},"spec":{"image":"quay.io/lrossett/telegram-send:latest","secret":{"name":"telegram-bot"},"text":"sample message from custom resource"}}
  creationTimestamp: "2020-12-29T17:32:37Z"
  generation: 1
  name: sample
  namespace: default
  resourceVersion: "1010"
  selfLink: /apis/telegram.my.domain/v1alpha1/namespaces/default/botmessages/sample
  uid: 42e7ea05-22bb-47dd-8adb-15d0302861a7
spec:
  image: quay.io/lrossett/telegram-send:latest
  secret:
    name: telegram-bot
  text: sample message from custom resource
status:
  message: Message successfully sent
  reason: Sent
  status: "True"
  type: Sent
```

You can create as many CRs you want, as long as each CR has a unique name.

