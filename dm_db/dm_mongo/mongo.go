package dm_mongo

import (
	_ "dm_server/dm_configuration"
	conf "dm_server/dm_configuration"
	_ "dm_server/dm_helper"
	h "dm_server/dm_helper"

	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getEntityStruct(entityName string) reflect.Type {
	entityMap := map[string]reflect.Type{
		"Projects":      reflect.TypeOf(Project{}),
		"Objects":       reflect.TypeOf(Object{}),
		"Requests":      reflect.TypeOf(Request{}),
		"Verifications": reflect.TypeOf(Verification{}),
		"Analogues":     reflect.TypeOf(Analogue{}),
		"RequestItems":  reflect.TypeOf(RequestItem{}),
		"ARCs":          reflect.TypeOf(ARC{}),
		"ARCWorks":      reflect.TypeOf(ARCWork{}),
		"ARCWorkItems":  reflect.TypeOf(ARCWorkItem{}),
		"Invoices":      reflect.TypeOf(Invoice{}),
		"InvoiceItems":  reflect.TypeOf(InvoiceItem{}),
		"Counteragents": reflect.TypeOf(Counteragent{}),
	}

	return entityMap[entityName]
}

var MongoDBClient *mongo.Client
var MongoDatabase *mongo.Database

func UpdateWithJSONFilterAndPayload(collectionName string, mdbQuery string, mdbPayload string) error {
	// Parse mdbQuery into bson.M filter
	var filter bson.M
	if mdbQuery != "" {
		var err error
		filter, err = getBsonQueryFromJSONString(mdbQuery)
		if err != nil {
			return fmt.Errorf("Failed to parse mdbQuery: %s", err)
		}
	}

	// Parse mdbPayload into bson.M update
	var update bson.M
	if mdbPayload != "" {
		err := json.Unmarshal([]byte(mdbPayload), &update)
		if err != nil {
			return fmt.Errorf("Failed to parse mdbPayload: %s", err)
		}
	}

	// Get collection from MongoDBClient and MongoDatabase
	collection := MongoDatabase.Collection(collectionName)

	// Update matching documents with filter and update
	result, err := collection.UpdateMany(context.Background(), filter, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("Failed to update documents: %s", err)
	}

	fmt.Printf("Updated %d documents\n", result.ModifiedCount)

	return nil
}

func GetCollectionObjects(entityName string, queryString string) ([]interface{}, error) {
	collection := MongoDatabase.Collection(entityName)

	filter, err := getBsonQueryFromJSONString(queryString)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse query string: %s", err)
	}

	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve documents: %s", err)
	}
	defer cur.Close(context.Background())

	// Получение значения структуры по имени с помощью отражения (reflection)
	entityType := getEntityStruct(entityName)
	if entityType == nil {
		return nil, fmt.Errorf("Unknown entity: %s", entityName)
	}

	// Создание слайса для хранения объектов
	objects := make([]interface{}, 0)

	for cur.Next(context.Background()) {
		obj := reflect.New(entityType).Interface() // Создание нового экземпляра структуры с помощью отражения (reflection)
		err := cur.Decode(obj)
		if err != nil {
			return nil, fmt.Errorf("Failed to decode document: %s", err)
		}
		objects = append(objects, obj)
	}

	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("Cursor error: %s", err)
	}

	return objects, nil
}

func CreateEntities(collectionName string, mdbPayload string) error {
	// Parse mdbPayload into bson.M document
	var document bson.M
	if mdbPayload != "" {
		err := json.Unmarshal([]byte(mdbPayload), &document)
		if err != nil {
			return fmt.Errorf("Failed to parse mdbPayload: %s", err)
		}
	}

	// Get collection from MongoDBClient and MongoDatabase
	collection := MongoDatabase.Collection(collectionName)

	// Create document
	_, err := collection.InsertOne(context.Background(), document)
	if err != nil {
		return fmt.Errorf("Failed to create document: %s", err)
	}

	return nil
}

func DeleteMongoEntities(collectionName string, mdbQuery string) error {
	// Parse mdbQuery into bson.M filter
	var filter bson.M
	if mdbQuery != "" {
		var err error
		filter, err = getBsonQueryFromJSONString(mdbQuery)
		if err != nil {
			return fmt.Errorf("Failed to parse mdbQuery: %s", err)
		}
	}

	// Get collection from MongoDBClient and MongoDatabase
	collection := MongoDatabase.Collection(collectionName)

	// Delete matching documents with filter
	result, err := collection.DeleteMany(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("Failed to delete documents: %s", err)
	}

	fmt.Printf("Deleted %d documents\n", result.DeletedCount)

	return nil
}

func getBsonQueryFromJSONString(jsonString string) (bson.M, error) {
	var filter bson.M
	err := json.Unmarshal([]byte(jsonString), &filter)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse JSON string: %s", err)
	}

	return filter, nil
}

func getBsonQueryFromString(queryString string) (bson.M, error) {
	queryMap := bson.M{}

	// Split the query string by the separator "&"
	pairs := strings.Split(queryString, "&")
	for _, pair := range pairs {
		// Split each key-value pair by the separator "="
		keyValue := strings.Split(pair, "=")
		if len(keyValue) != 2 {
			return nil, fmt.Errorf("Invalid query string format")
		}

		// Add the key-value pair to the query map
		queryMap[keyValue[0]] = keyValue[1]
	}

	return queryMap, nil
}

func InitMongoDB() error {
	dsn := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s",
		conf.Conf_Mongo.User,
		conf.Conf_Mongo.Password,
		conf.Conf_Mongo.Host,
		conf.Conf_Mongo.Port,
		conf.Conf_Mongo.DBName)

	clientOptions := options.Client().ApplyURI(dsn)

	h.Log("Создание клиента для подключения к серверу MongoDB... \n" + dsn)

	MongoDBClient, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		h.Err("Ошибка: \n" + err.Error())
		return err
	}

	h.Log("Проверка соединения с сервером MongoDB... \n" + dsn)

	MongoDBClient.Ping(context.Background(), nil)
	if err != nil {
		h.Err("Ошибка: \n" + err.Error())
		return err
	}

	tmpStr := h.YellowColor + conf.Conf_Mongo.DBName
	h.Log("Подключение к базе данных..." + tmpStr)
	MongoDatabase = MongoDBClient.Database(conf.Conf_Mongo.DBName)
	if MongoDatabase == nil {
		h.Err("Не удалось открыть базу данных: " + tmpStr)
	}
	h.Log("Connected to database " + tmpStr)

	return nil
}
