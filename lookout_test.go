package config

import (
	"testing"
	"os"
	"io/ioutil"
	"github.com/stretchr/testify/assert"

)

func Test_NewConfig_Yaml(t *testing.T) {

	configFile := createTempFile(yaml1())
	defer os.Remove(configFile)

	config, err := NewConfig(configFile)

	assert.NoError(t, err, "Should not have error when creating config")
	assert.Equal(t, config.GetString("foo", "default"), "bar")

	invalidConfigFile := createTempFile(invalidYaml())
	defer os.Remove(invalidConfigFile)

	_, err = NewConfig(invalidConfigFile)
	assert.Error(t, err, "Should error since invalid file")
}

func Test_NewConfig_Json(t *testing.T) {

	configFile := createTempFile(json1())
	defer os.Remove(configFile)

	config, err := NewConfig(configFile)

	assert.NoError(t, err, "Should not have error when creating config")
	assert.Equal(t, config.GetString("foo", "default"), "bar")
}

func Test_GetByType_Yaml(t *testing.T) {
	configFile := createTempFile(yaml1())
	defer os.Remove(configFile)

	config, err := NewConfig(configFile)

	assert.NoError(t, err, "Should not have error when creating config")

	//Test can get string
	assert.Equal(t, "bar", config.GetString("foo", "default"))
	assert.Equal(t, "world", config.GetString("hello", "default"))

	//Test can get int
	assert.Equal(t, 77,config.GetInt("age", 12))
	assert.NotEqual(t, 77, config.GetInt("hello", 12))
	assert.Equal(t, 12, config.GetInt("hello", 12))

	//Test can get bool
	assert.Equal(t, true, config.GetBool("bob", false))
	assert.Equal(t, false, config.GetBool("hello", false))


	//Invalid key
	assert.Equal(t, "default", config.GetString("golang", "default"))


}
func Test_UpdateOnWrite_Yaml(t *testing.T) {

	configFile := createTempFile(yaml1())
	defer os.Remove(configFile)

	config, err := NewConfig(configFile)

	assert.NoError(t, err, "Should not have error when creating config")
	assert.Equal(t, 77, config.GetInt("age", 12),)

	writeTempFile(yaml2(), configFile)

	assert.Equal(t, 32, config.GetInt("age", 12),)
}

func Test_UpdateOnWrite_Json(t *testing.T) {

	configFile := createTempFile(json1())
	defer os.Remove(configFile)

	config, err := NewConfig(configFile)

	assert.NoError(t, err, "Should not have error when creating config")
	assert.Equal(t, 77, config.GetInt("age", 12),)

	writeTempFile(json2(), configFile)

	assert.Equal(t, 32, config.GetInt("age", 12),)
}

func Test_UpdateOnWrite_Mixed(t *testing.T) {

	configFile := createTempFile(yaml1())
	defer os.Remove(configFile)

	config, err := NewConfig(configFile)

	assert.NoError(t, err, "Should not have error when creating config")
	assert.Equal(t, 77, config.GetInt("age", 12),)

	writeTempFile(json2(), configFile)

	assert.Equal(t, 32, config.GetInt("age", 12),)
}

func Test_UpdateWithBadConfig(t *testing.T) {
	configFile := createTempFile(yaml1())

	config, err := NewConfig(configFile)

	assert.NoError(t, err, "Should not have error when creating config")

	writeTempFile( invalidYaml(), configFile)

	assert.Equal(t, 77, config.GetInt("age", 12),)
}

func Test_DeletedConfig(t *testing.T) {
	configFile := createTempFile(yaml1())

	config, err := NewConfig(configFile)

	assert.NoError(t, err, "Should not have error when creating config")
	assert.Equal(t, 77, config.GetInt("age", 12),)
	os.Remove(configFile)

	assert.Equal(t, 77, config.GetInt("age", 12),)
}

func yaml1() []byte {
	return []byte(`hello: world
foo: bar
bob: true
age: 77
jersey: "77"`)
}
func yaml2() []byte {
	return []byte(`hello: golang
foo: bar
bob: false
age: 32
jersey: "77"`)
}

func invalidYaml() []byte {
	return []byte(`hello: world
foo: bar
bob: true
age: 77
jersey "77"`) //The last object is missing a :
}

func json1() []byte {
	return []byte(`{
  "hello": "world",
  "foo": "bar",
  "bob": "true",
  "age": 77,
  "jersey": "77"
}`)
}

func json2() []byte {
	return []byte(`{
  "hello": "golang",
  "foo": "bar",
  "bob": "false",
  "age": 32,
  "jersey": "77"
}`)
}
func createTempFile(data []byte) string {
	temp := tempFile()
	return writeTempFile(data, temp.Name())
}
func writeTempFile(data []byte, filenameString string) string {


	if err := ioutil.WriteFile(filenameString, data, 0644); err != nil {
		panic(err)
	}
	return filenameString
}

func tempFile() *os.File {
	tmpDir, err := ioutil.TempDir("", "test")

	if err != nil {
		panic(err)
	}
	tmpFile, err := ioutil.TempFile(tmpDir, "test")
	if err != nil {
		panic(err)
	}
	return tmpFile
}