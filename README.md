# huawei-push-authorizator
Huawei Push Kit authorizator in time

#### Why?

To send push messages via huawei push kit access token must be taken not often (about 1000 via 5 minuts). To get access token for request before each operation huawei service delay about 3-5 seconds and not acceptable in some use cases for example send notification for media call.

This microservice can reduce delay for send push notification by get access token in background and always returns valid access token for clients.

Documentation:

https://developer.huawei.com/consumer/en/doc/development/HMSCore-Guides/oauth2-0000001212610981#section128682386159

For use in docker application use public image oh docker hub: ekudinov/huawei-push-authorizator
