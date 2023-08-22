// DPBPR project main.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

type individual struct {
	Name string `json:"name"`
	Age string  `json:"age"`
	IDNum string `json:"IDNum"`
	Address string `json:"address"`
	Job string `json:"job"`
	Extra string `json:"extra"`
}

type physiology struct {
	HeartRate string `json:"heartRate"`
	Pulse string `json:"pulse"`
	BloodPres string `json:"bloodPres"`
	Extra string `json:"extra"`
}

type diagnosis struct{ //对应JSON字段的首字母要大写
	Content string `json:"content"`  
	Author string `json:"author"`
	Time int64 `json:"time"`
}

type UserData struct {
	Individual individual `json:"individual"`
	Physiology physiology `json:"physiology"`
	Diagnosis  diagnosis  `json:"diagnosis"`
}


type Permission struct{
	Doctors      string
	Pharmacists  string
	Academics    string
}
///////////////////////////////////////////////////////////////数据记录的定义

type Right struct{
	Role string    `json:"role"`
	Individual int `json:"individual"`
	Physiology int `json:"physiology"`
	Diagnosis  int `json:"diagnosis"`
}

type FilePermission struct{
	FileId string        `json:"fileId"`
	Rights []Right       `json:"rights"`
}
///////////////////////////////////////////////////////////////数据权限的定义

func HashSha256(src string) string {
	m := sha256.New()
    m.Write([]byte(src))
    res := hex.EncodeToString(m.Sum(nil))
    return res	
}

func getTargetFile(w http.ResponseWriter, r *http.Request){  //获取服务器中的文件

	roles := make(map[string]string) 
	roles["0"] = "Doctors"  
	roles["1"] = "Pharmacists" 
	roles["2"] = "Academics"
	fileContent, err := os.Open("D:\\Truffle_test\\user_data.js") //读取健康检测数据
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("The file is download successfully...")
	defer fileContent.Close()
	byteResult, _ := ioutil.ReadAll(fileContent)
	
	var user_data UserData
	json.Unmarshal([]byte(byteResult), &user_data)
	
	params := r.URL.Query()
	fmt.Println(params["fileId"][0]); //获取请求的关联参数
	fmt.Println(params["role"][0]);
	
	req_url := "http://127.0.0.1:3000?op=get&fileId="+params["fileId"][0]
	resp, err := http.Get(req_url)    //访问文件的权限获取服务 fileId: 01234546789
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("Access-List; "+string(body))
	ret := string(body)
	if ret == "error" {
		fmt.Println("we get nothing now")
		fmt.Fprintf(w,"server error!") 
		return
	}
	
	str_perm := strings.Split(ret,":")  //权限列表拆分
	var perm Permission
	perm.Doctors=str_perm[0]
	perm.Pharmacists = str_perm[1]
	perm.Academics = str_perm[2]
	//fmt.Println(perm.Doctors)
	
	role_tag := roles[params["role"][0]]        //"Pharmacists"
	var specific_right string         //访问者具体角色
	switch role_tag {
      case "Doctors": specific_right = perm.Doctors
      case "Pharmacists": specific_right = perm.Pharmacists
      case "Academics": specific_right = perm.Academics
    }
	
	if specific_right[0]-'0' == 0{   //字段数据掩码
		user_data.Individual.Name = HashSha256(user_data.Individual.Name)
		user_data.Individual.Age = HashSha256(user_data.Individual.Age)
		user_data.Individual.Address = HashSha256(user_data.Individual.Address)
		user_data.Individual.IDNum = HashSha256(user_data.Individual.IDNum)
		user_data.Individual.Extra = HashSha256(user_data.Individual.Extra)
		user_data.Individual.Job = HashSha256(user_data.Individual.Job)
	}
	if specific_right[1]-'0' == 0{
		user_data.Physiology.HeartRate = HashSha256(user_data.Physiology.HeartRate)
		user_data.Physiology.Pulse = HashSha256(user_data.Physiology.Pulse)
		user_data.Physiology.BloodPres = HashSha256(user_data.Physiology.BloodPres)
		user_data.Physiology.Extra = HashSha256(user_data.Physiology.Extra)
	}
	
	if specific_right[2]-'0' == 0{
		user_data.Diagnosis.Content = HashSha256(user_data.Diagnosis.Content)
		user_data.Diagnosis.Author = HashSha256(user_data.Diagnosis.Author)
	}
	
	fmt.Println(user_data)
	
    jsonBytes, err := json.Marshal(user_data)
    if err != nil {
        // handle the error
    }
    jsonString := string(jsonBytes)
    fmt.Println(jsonString)
	//fmt.Println(resp.StatusCode)
	fmt.Fprintf(w, jsonString) 
}


func sendTargetFileAcl(w http.ResponseWriter, r *http.Request){ //w http.ResponseWriter, r *http.Request
	
	var fileRights FilePermission
	params := r.URL.Query()
	fileId := params["fileId"][0]  //从请求中获取文件标识符
	
	fileName := "D:\\Truffle_test\\user_files\\file_permissions_"+fileId+".js"  
	fileContent, err := os.Open(fileName) //读取健康检测数据的权限定义
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("The file is download successfully...")
	defer fileContent.Close()
	byteResult, _ := ioutil.ReadAll(fileContent)
	
	
    json.Unmarshal([]byte(byteResult), &fileRights)
	
	for j:=0; j<len(fileRights.Rights);j++{
		fmt.Println("Role:",fileRights.Rights[j].Role)
		fmt.Println("ACL:",fileRights.Rights[j].Individual,fileRights.Rights[j].Physiology,fileRights.Rights[j].Diagnosis)
	}
	
	var url string
	var d_rights string         //医生权限控制字符串
	var p_rights string         //Pharmacists
	var a_rights string 
	url = "http://127.0.0.1:3000?op=set&fileId="+fileRights.FileId
	for i:=0; i<len(fileRights.Rights);i++{
		
		if fileRights.Rights[i].Role == "Doctors"{
			d_rights = strconv.Itoa(fileRights.Rights[i].Individual)+ "_"+strconv.Itoa(fileRights.Rights[i].Physiology)+"_"+strconv.Itoa(fileRights.Rights[i].Diagnosis)
		}
		
		if fileRights.Rights[i].Role == "Pharmacists" {
			p_rights = strconv.Itoa(fileRights.Rights[i].Individual)+"_"+strconv.Itoa(fileRights.Rights[i].Physiology)+"_"+strconv.Itoa(fileRights.Rights[i].Diagnosis)
			
		}
		if fileRights.Rights[i].Role == "Academics" {
			a_rights = strconv.Itoa(fileRights.Rights[i].Individual)+"_"+strconv.Itoa(fileRights.Rights[i].Physiology)+"_"+strconv.Itoa(fileRights.Rights[i].Diagnosis)
		}
	}
	
	url = url + "&d="+d_rights+"&p="+p_rights+"&a="+ a_rights
	fmt.Println("Building request url as: "+url)
	resp, err := http.Get(url)    //访问文件的权限获取服务 fileId: 01234546789
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	ret := string(body)
	if ret == "error" {
		fmt.Println("we get nothing now")
		fmt.Fprintf(w,"server error!") 
		return
	}else{
		fmt.Fprintf(w,ret)
	}
	
}



func main() {
	//sendTargeFileAcl()
	http.HandleFunc("/obtainFile", getTargetFile)  
	log.Println("Server is working on port 8080")
	//http.ListenAndServe(":8080", nil)
	http.HandleFunc("/sendAcl",sendTargetFileAcl) //设定诊断数据的访问控制权限
	log.Fatal(http.ListenAndServe(":8080", nil)) 
}
