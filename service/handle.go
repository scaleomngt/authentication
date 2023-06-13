package service

import (
	"errors"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	invalidator "github.com/guanguans/id-validator"
	"id-card-server/config"
	"id-card-server/gintool"
	"id-card-server/model"
	"id-card-server/utils"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

var once sync.Once
var Accounts map[string]model.Account
var Index uint32
var lock sync.RWMutex

var PrivateKey = config.Config.GetString("PrivateKey")
var ViewKey = config.Config.GetString("ViewKey")
var ApiUrl = config.Config.GetString("ApiUrl")
var Contract = config.Config.GetString("Contract")

const Query = "https://vm.aleo.org/api"
const Broadcast = "https://vm.aleo.org/api/testnet3/transaction/broadcast"
const Prefix = "at"
const RPrefix = "card-server"
const FieldPrefix = "11"
const FEE = "100000"

// CreateAccount 创建账户
func CreateAccount(c *gin.Context) {
	once.Do(func() {
		Accounts = make(map[string]model.Account, 0)
		Index = 0
	})

	cmd := "snarkos"
	args := []string{
		"account",
		"new"}
	result, err := utils.ExecCmdWithTimeout(60, cmd, args...)
	if err != nil {
		log.Println("err:", err, result)
		log.Println("result:", result)
		gintool.ResultFail(c, err.Error())
		return
	}
	log.Println("result:", result)
	keys := strings.Split(strings.TrimSpace(result), "\n")
	privateKey := strings.TrimSpace(strings.Split(keys[0], "Private Key")[1])
	viewKey := strings.TrimSpace(strings.Split(keys[1], "View Key")[1])
	address := strings.TrimSpace(strings.Split(keys[2], "Address")[1])
	Index = Index + 1
	account := model.Account{Index: Index, PrivateKey: privateKey, ViewKey: viewKey, Address: address}
	Accounts[account.Address] = account
	log.Println("Accounts:", Accounts)
	gintool.ResultOk(c, account)
}

// GetAccounts 获取GetAccounts
func GetAccounts(c *gin.Context) {
	if Accounts == nil {
		gintool.ResultOk(c, make([]model.Account, 0))
		return
	}

	gintool.ResultOk(c, Accounts)
}

// InitRedisId setRedis的id值
func InitRedisId(c *gin.Context) {
	id := c.Param("id")
	log.Println("id:", id)
	err := utils.SetId(id)
	if err != nil {
		log.Println("Error:", err)
		gintool.ResultFail(c, err.Error())
		return
	}
	gintool.ResultOkMsg(c, id, "success")
}

func CalcData(c *gin.Context) {
	lock.Lock()
	defer lock.Unlock()
	info := new(model.Transaction)
	if err := c.ShouldBindJSON(&info); err != nil {
		log.Println("Error:", err)
		gintool.ResultFail(c, err.Error())
		return
	}

	log.Println("info: ", info)
	id := info.Id

	if info.Id == "" || info.Address == "" {
		gintool.ResultFail(c, "Id or Address is not null")
		return
	}

	type Data struct {
		Ti     string `json:"ti"`
		Name   string `json:"name"`
		Gender string `json:"gender"`
		Result string `json:"result"`
		Owner  string `json:"owner"`
		Id     string `json:"id"`
	}

	hash, err := utils.GetOneHash(RPrefix+"-"+info.Address, info.Id)
	if err != nil {
		gintool.ResultFail(c, err.Error())
		return
	}

	if hash != "" {
		data := &Data{}
		data.Id = info.Id
		data.Result = hash
		gintool.ResultOkMsg(c, data, "success")
		return
	}

	time.Sleep(time.Duration(1000*10) * time.Millisecond)

	// 2. 获取步骤1中的output->value数据
	ciphertext, err := GetExecOutputValue(id)
	if err != nil {
		gintool.ResultFail(c, err.Error())
		return
	}

	// 获取计算Card数据的入参
	value, err := DecryptCiphertext(ciphertext)
	if err != nil {
		gintool.ResultFail(c, err.Error())
		return
	}

	// 获取最新record数据
	record, err := GetLatestFeeRecord()
	if err != nil {
		return
	}

	// 3.执行合约计算健康数据
	id, err = CalcCardData(record, value)
	if err != nil {
		gintool.ResultFail(c, err.Error())
		return
	}

	time.Sleep(time.Duration(1000*18) * time.Millisecond)

	// 获取计算结果
	ciphertext, err = GetExecOutputValue(id)
	if err != nil {
		gintool.ResultFail(c, err.Error())
		return
	}

	// 解析结果
	value, err = DecryptCiphertext(ciphertext)
	if err != nil {
		gintool.ResultFail(c, err.Error())
		return
	}

	data := &Data{}

	records := strings.Split(value, "\n")

	for i := 0; i < len(records); i++ {
		if strings.HasPrefix(strings.TrimSpace(records[i]), "owner:") {
			temp := strings.TrimSpace(strings.Replace(strings.Split(strings.TrimSpace(records[i]), ":")[1], ".private,", "", -1))
			data.Owner = temp
		} else if strings.HasPrefix(strings.TrimSpace(records[i]), "name:") {
			temp := strings.TrimSpace(strings.Replace(strings.Split(strings.TrimSpace(records[i]), ":")[1], "field.private,", "", -1))
			data.Name = string(StringToBytes(temp))
		} else if strings.HasPrefix(strings.TrimSpace(records[i]), "result:") {
			temp := strings.TrimSpace(strings.Replace(strings.Split(strings.TrimSpace(records[i]), ":")[1], "u32.private,", "", -1))
			data.Result = temp
		} else if strings.HasPrefix(strings.TrimSpace(records[i]), "id:") {
			temp := strings.TrimSpace(strings.Replace(strings.Split(strings.TrimSpace(records[i]), ":")[1], "field.private,", "", -1))
			log.Printf("Output id: %v", temp)
			data.Id = info.Id
		} else if strings.HasPrefix(strings.TrimSpace(records[i]), "gender:") {
			temp := strings.TrimSpace(strings.Replace(strings.Split(strings.TrimSpace(records[i]), ":")[1], "field.private,", "", -1))
			data.Gender = string(StringToBytes(temp))
		}
	}
	err = utils.SetHash(RPrefix+"-"+info.Address, info.Id, data.Result)
	if err != nil {
		gintool.ResultFail(c, err.Error())
		return
	}

	data.Ti = id
	gintool.ResultOkMsg(c, data, "success")
}

func StringToBytes(str string) []byte {
	var result []byte
	length := 3

	if len(str)%3 == len(FieldPrefix) {
		str = str[len(FieldPrefix):]
	}

	// 循环切割字符串并添加到结果切片中
	for i := 0; i < len(str); i += length {
		endIndex := i + length
		if endIndex > len(str) {
			endIndex = len(str)
		}
		substr := str[i:endIndex]

		decimal, err := strconv.Atoi(substr)
		if err != nil {
			log.Println("error parsing: ", err)
			result = append(result, 0)
		}
		result = append(result, byte(decimal))
	}

	return result
}

func calculateAge(birthDate time.Time) int {
	now := time.Now()
	age := now.Year() - birthDate.Year()

	// 检查生日是否已过
	if now.Month() < birthDate.Month() || (now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
		age--
	}

	return age
}

// SubmitData 提交数据
func SubmitData(c *gin.Context) {
	lock.Lock()
	defer lock.Unlock()
	info := new(model.CardInfo)
	if err := c.ShouldBindJSON(&info); err != nil {
		log.Println("Error:", err)
		gintool.ResultFail(c, err.Error())
		return
	}
	log.Println("info: ", info)

	// 计算年龄
	birthDate, err := time.Parse("2006-01-02", info.Birthdate)
	if err != nil {
		gintool.ResultFail(c, err.Error())
		return
	}
	age := calculateAge(birthDate)
	info.Birthdate = strings.Replace(info.Birthdate, "-", "", -1)
	// 获取最新record数据
	record, err := GetLatestFeeRecord()
	if err != nil {
		gintool.ResultFail(c, err.Error())
		return
	}

	// 1.执行合约提交Card数据
	id, err := SubmitCardData(record, info, age)
	if err != nil {
		gintool.ResultFail(c, err.Error())
		return
	}

	err = utils.SetHash(RPrefix+"-"+info.Addr, id, "")
	if err != nil {
		gintool.ResultFail(c, err.Error())
		return
	}

	time.Sleep(time.Duration(1000*10) * time.Millisecond)

	r := &model.Transaction{
		Id:      id,
		Address: info.Addr,
		Result:  "",
	}
	gintool.ResultOk(c, r)
}

func CalcCardData(record, value string) (string, error) {
	id := ""
	method := "calculate_age"

	args := []string{
		"developer",
		"execute",
		Contract,
		method,
		value,
		"--private-key",
		PrivateKey,
		"--query",
		Query,
		"--broadcast",
		Broadcast,
		"--fee",
		FEE,
		"--record",
		record}

	log.Println("args: ", args)

	cmd := "snarkos"
	result, err := utils.ExecCmdWithTimeout(60*5, cmd, args...)
	if err != nil {
		log.Println("执行合约计算Card数据 err:", err)
		log.Println("执行合约计算Card数据 result:", result)
		return id, err
	}
	log.Println("执行合约计算Card数据 result:", result)

	split := strings.Split(strings.TrimSpace(result), "\n")
	id = split[len(split)-1]
	if !strings.HasPrefix(id, Prefix) {
		log.Println("执行报错, 获取的数据有误, id: ", id)
		return "", errors.New(id)
	}

	err = utils.SetId(id)
	if err != nil {
		log.Println("utils.SetId err:", err)
		return "", err
	}
	return id, nil
}

func DecimalToPaddedString(data string) string {
	result := ""
	bytes := []byte(data)
	for i := 0; i < len(bytes); i++ {
		paddedString := fmt.Sprintf("%03d", bytes[i])
		result = result + paddedString
	}
	// field前缀不能为0
	if strings.HasPrefix(result, "0") {
		result = FieldPrefix + result
	}
	return result
}

func SubmitCardData(record string, info *model.CardInfo, age int) (string, error) {
	id := ""

	params := `{gender: {{gender}}field, name: {{name}}field, age: {{age}}u32, nation: {{nation}}field, birthdate: {{birthdate}}field, addr: {{addr}}}`
	params = strings.Replace(params, "{{gender}}", DecimalToPaddedString(info.Gender), -1)
	params = strings.Replace(params, "{{name}}", DecimalToPaddedString(info.Name), -1)
	params = strings.Replace(params, "{{age}}", strconv.Itoa(age), -1)
	params = strings.Replace(params, "{{nation}}", DecimalToPaddedString(info.Nation), -1)
	params = strings.Replace(params, "{{birthdate}}", info.Birthdate, -1)
	params = strings.Replace(params, "{{addr}}", info.Addr, -1)
	args := []string{
		"developer",
		"execute",
		Contract,
		"submit",
		params,
		"--private-key",
		PrivateKey,
		"--query",
		Query,
		"--broadcast",
		Broadcast,
		"--fee",
		FEE,
		"--record",
		record}
	log.Println("执行合约提交Card数据 args:", args)

	cmd := "snarkos"
	result, err := utils.ExecCmdWithTimeout(60*5, cmd, args...)
	if err != nil {
		log.Println("执行合约提交Card数据 err:", err)
		log.Println("执行合约提交Card数据 result:", result)
		return id, err
	}
	log.Println("执行合约提交Card数据 result:", result)

	split := strings.Split(strings.TrimSpace(result), "\n")
	id = split[len(split)-1]
	if !strings.HasPrefix(id, Prefix) {
		log.Println("执行报错, 获取的数据有误, id: ", id)
		return "", errors.New(id)
	}

	log.Println("执行合约提交Card数据 id: ", id)

	err = utils.SetId(id)
	if err != nil {
		log.Println("utils.SetId err:", err)
		return id, err
	}

	return id, nil
}

// GetLatestFeeRecord 获取fee transition outputs 0 value
func GetLatestFeeRecord() (string, error) {
	id, err := utils.GetId()
	if err != nil {
		return "", err
	}

	client := resty.New()
	resp, err := client.R().
		EnableTrace().
		Get(ApiUrl + id)
	if err != nil {
		log.Println("获取最新record数据 error: ", err)
		return "", err
	}

	ciphertext, err := jsonparser.GetString(resp.Body(), "fee", "transition", "outputs", "[0]", "value")
	if err != nil {
		log.Println("获取最新record数据 value error: ", err)
		return "", err
	}

	record, err := DecryptCiphertext(ciphertext)
	if err != nil {
		log.Println("获取最新record数据 进行解密 error: ", err)
		return "", err
	}
	log.Println("获取最新record数据: ", record)
	return record, nil
}

func DecryptCiphertext(ciphertext string) (string, error) {
	cmd := "snarkos"

	args := []string{
		"developer",
		"decrypt",
		"--ciphertext",
		ciphertext,
		"--view-key",
		ViewKey}
	log.Println("args: ", args)
	record, err := utils.ExecCmdWithTimeout(60, cmd, args...)
	if err != nil {
		log.Println("DecryptCiphertext err:", err, record)
		log.Println("result:", record)
		return "", err
	}
	log.Println("DecryptCiphertext record:", strings.TrimSpace(record))

	return strings.TrimSpace(record), nil
}

// GetExecOutputValue 获取 execution transitions value
func GetExecOutputValue(id string) (string, error) {
	client := resty.New()
	resp, err := client.R().
		EnableTrace().
		Get(ApiUrl + id)
	if err != nil {
		log.Println("发送http请求, 获取output->value数据 error: ", err)
		return "", err
	}
	cipherText, err := jsonparser.GetString(resp.Body(), "execution", "transitions", "[0]", "outputs", "[0]", "value")
	if err != nil {
		log.Println("获取json中 output->value数据 error: ", err)
		return "", err
	}
	log.Println("获取output->value数据 cipherText: ", cipherText)
	return cipherText, nil
}

func Test(c *gin.Context) {
	id := c.Param("id")
	log.Println(id)
	gintool.ResultOkMsg(c, invalidator.FakeId(), "success")
}
