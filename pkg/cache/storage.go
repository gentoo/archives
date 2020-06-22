package cache

var data map[string]interface{}

func Init(){
	data = make(map[string]interface{})
}

func Put(key string, value interface{}) {
	data[key] = value
}

func Get(key string) interface{} {
	return data[key]
}
