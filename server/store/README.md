# Database Store Implementation

Este directorio contiene la implementación del store para la conexión a la base de datos de Mattermost.

## Estructura

```
store/
├── store.go              # Interface Store
├── sqlstore/
│   ├── sqlstore.go       # Implementación SQL del store
│   └── migrations/
│       └── 001_create_messages_table.sql  # Migración para crear tabla de mensajes
└── README.md             # Este archivo
```

## Uso

### Inicialización

El store se inicializa automáticamente en el método `OnActivate()` del plugin:

```go
// Initialize database store
store, err := sqlstore.New(p.API)
if err != nil {
    return errors.Wrap(err, "failed to initialize database store")
}
p.store = store
```

### Métodos disponibles

El store implementa la interfaz `Store` con los siguientes métodos:

- `GetMessage(messageId string) (*models.Message, error)` - Obtiene un mensaje por ID
- `CreateMessage(message *models.Message) error` - Crea un nuevo mensaje
- `UpdateMessage(message *models.Message) error` - Actualiza un mensaje existente
- `DeleteMessage(messageId string) error` - Elimina un mensaje
- `GetMessagesByChannel(channelId string, limit, offset int) ([]*models.Message, error)` - Obtiene mensajes por canal
- `GetMessagesByUser(userId string, limit, offset int) ([]*models.Message, error)` - Obtiene mensajes por usuario

### Helpers del Plugin

El plugin incluye métodos helper para facilitar el uso del store:

- `SaveMessageToStore(post *model.Post) error` - Guarda un mensaje de Mattermost en el store
- `GetMessageFromStore(messageId string) (*models.Message, error)` - Obtiene un mensaje del store
- `GetChannelMessagesFromStore(channelId string, limit, offset int) ([]*models.Message, error)` - Obtiene mensajes de un canal
- `GetUserMessagesFromStore(userId string, limit, offset int) ([]*models.Message, error)` - Obtiene mensajes de un usuario
- `UpdateMessageInStore(message *models.Message) error` - Actualiza un mensaje en el store
- `DeleteMessageFromStore(messageId string) error` - Elimina un mensaje del store

### Ejemplo de uso

```go
// Obtener mensajes de un canal
messages, err := p.GetChannelMessagesFromStore(channelId, 10, 0)
if err != nil {
    p.API.LogError("Failed to get channel messages", "error", err.Error())
    return
}

// Crear un nuevo mensaje
message := &models.Message{
    ID:        "unique-id",
    ChannelID: "channel-id",
    UserID:    "user-id",
    Content:   "Hello World",
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
}

err = p.store.CreateMessage(message)
if err != nil {
    p.API.LogError("Failed to create message", "error", err.Error())
    return
}
```

## Migraciones

Las migraciones SQL se encuentran en `sqlstore/migrations/`. Para aplicar las migraciones, ejecuta el archivo SQL correspondiente en tu base de datos de Mattermost.

### Tabla de mensajes

La tabla `messages` tiene la siguiente estructura:

```sql
CREATE TABLE messages (
    id VARCHAR(26) PRIMARY KEY,
    channel_id VARCHAR(26) NOT NULL,
    user_id VARCHAR(26) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
```

## Configuración

El store utiliza la configuración de base de datos de Mattermost automáticamente. No se requiere configuración adicional.

## Limpieza

El store se cierra automáticamente en el método `OnDeactivate()` del plugin para liberar recursos.
