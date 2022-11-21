Start the nacos service and open the ui interface in your browser `http://<ip>:8848/nacos/index.html`

- Namespaces: to distinguish between different services
- Groups: used to distinguish between different environments
- Configuration sets: specific service configurations

To create a development environment service configuration process.

Click the left menu bar [Namespace] --> Click [New Namespace], example: ID is empty, namespace name: serverName, description: xxx service, then click the left menu bar [Configuration Management] --> [Configuration List] --> [serverName] --> [+] New Configuration, enter Data ID: serverName .yml, Group: dev, Configuration format: yaml, fill in the configuration content.
