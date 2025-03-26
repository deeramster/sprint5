# Schema Registry в Yandex Cloud

Этот документ описывает развертывание и настройку **Schema Registry** для Apache Kafka в Yandex Cloud.

## Описание
Schema Registry используется для управления схемами сообщений, передаваемых через Kafka. В данной конфигурации используется **SASL_SSL** для аутентификации и защищенного подключения к Kafka-кластеру.

## Конфигурация Schema Registry

### Файл конфигурации
Пример конфигурации Schema Registry (`schema-registry.properties`):

```properties
listeners=http://0.0.0.0:8081
kafkastore.topic=_schemas
debug=false
kafkastore.bootstrap.servers=rc1a-ii0r6jluacd0pob4.mdb.yandexcloud.net:9091,rc1b-3j7j6kdodgbr87u6.mdb.yandexcloud.net:9091,rc1d-leu6lr8j2f4753q1.mdb.yandexcloud.net:9091
kafkastore.ssl.truststore.location=/etc/schema-registry/client.truststore.jks
kafkastore.ssl.truststore.password=111111
kafkastore.sasl.mechanism=SCRAM-SHA-512
kafkastore.security.protocol=SASL_SSL
```

### Запуск Schema Registry

```bash
schema-registry-start /etc/schema-registry/schema-registry.properties
```

## Проверка работы

### Проверка доступности Schema Registry

```bash
curl -X GET http://localhost:8081/subjects
```

## Возможные ошибки и их решения

**Решение:**
- Проверьте, запущен ли Schema Registry (`netstat -tulnp | grep 8081`).
- Убедитесь, что порты не блокируются firewall'ом (`iptables -L -n`).
- Проверьте доступность брокеров Kafka (`telnet rc1a-ii0r6jluacd0pob4.mdb.yandexcloud.net 9091`).

### 2. Проблемы с сертификатами
Если используется **SASL_SSL**, убедитесь, что **truststore** настроен корректно и содержит валидные сертификаты.

```bash
keytool -list -keystore /etc/schema-registry/client.truststore.jks -storepass 111111
```


### Регистрация схемы
```bash
curl -X POST http://localhost:8081/subjects/json-subject/versions \
  -H "Content-Type: application/vnd.schemaregistry.v1+json" \
  --data '{"schemaType": "JSON", "schema": "{\"type\":\"object\",\"properties\":{\"id\":{\"type\":\"string\"},\"text\":{\"type\":\"string\"},\"timestamp\":{\"type\":\"integer\"}},\"required\":[\"id\",\"text\",\"timestamp\"]}"}'
```

### Получение схемы
```bash
curl -X GET http://localhost:8081/subjects/json-subject/versions/latest
```

## Заключение
Этот документ описывает основные шаги настройки Schema Registry в Yandex Cloud. Если у вас возникли проблемы, проверьте конфигурацию и доступность компонентов.

