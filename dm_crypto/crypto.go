package dm_crypto

import (
	_ "dm_server/dm_configuration"
	conf "dm_server/dm_configuration"
	_ "dm_server/dm_helper"
	h "dm_server/dm_helper"

	"encoding/hex"
	"strings"

	"crypto/md5"

	"github.com/ftomza/gogost/gost34112012256"
	"github.com/ftomza/gogost/gost34112012512"
)

const HashFunctionStreebog = "streebog" //512
const HashFunctionStreebog256 = "streebog256"
const HashFunctionMD5 = "md5"

var currentHashFunction = ""
var _md5 = md5.New()
var _streebog = gost34112012512.New()
var _streebog256 = gost34112012256.New()

func SetHashFunction(value string) {
	currentHashFunction = strings.ToLower(value)

	switch currentHashFunction {
	case HashFunctionMD5:
		{
			currentHashFunction = HashFunctionMD5
		}
	case HashFunctionStreebog:
		{
			currentHashFunction = HashFunctionStreebog
		}
	case HashFunctionStreebog256:
		{
			currentHashFunction = HashFunctionStreebog256
		}
	default:
		{
			currentHashFunction = HashFunctionMD5
		}
	}

}

func InitCrypto() {
	SetHashFunction(conf.Crypto.HashFunction)
	data := []byte("DM Corporation")

	hashString := ""
	hashString = Sum(data)
	h.Log("Detault hash function is " + currentHashFunction)
	h.Log("testing " + currentHashFunction + ":\n\t" + hashString)
	h.Log("Switching hash function to " + HashFunctionMD5)
	SetHashFunction(HashFunctionMD5)
	hashString = Sum(data)
	h.Log("testing " + HashFunctionMD5 + ":\n\t" + hashString)
	h.Log("Switching hash function to " + HashFunctionStreebog)
	SetHashFunction(HashFunctionStreebog256)
	hashString = Sum(data)
	h.Log("testing " + HashFunctionStreebog256 + ":\n\t" + hashString)
	SetHashFunction(HashFunctionStreebog)
	hashString = Sum(data)
	h.Log("testing " + HashFunctionStreebog + ":\n\t" + hashString)

	h.Log("DM Cryptographic module initialized with hash function (" + currentHashFunction + ") from configuration file (" + conf.Crypto.HashFunction + ")")
	// fmt.Println("\t" + +"(" + conf.Crypto.HashFunction + ")")
}

func Sum(data []byte) string {
	result := ""

	switch currentHashFunction {
	case HashFunctionMD5:
		{
			_md5.Reset()
			_md5.Write(data)
			result = hex.EncodeToString(_md5.Sum(nil))
		}
	case HashFunctionStreebog:
		{
			_streebog.Reset()
			_streebog.Write(data)
			result = hex.EncodeToString(_streebog.Sum(nil))
		}
	case HashFunctionStreebog256:
		{
			_streebog256.Reset()
			_streebog256.Write(data)
			result = hex.EncodeToString(_streebog256.Sum(nil))
		}
	default:
		{
			h.Err("Не задана currentHashFunction (невозможно, надо хотя бы вызвать SetHashFunction)")
		}
	}

	return result
}
